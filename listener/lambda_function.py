import json
import logging
import os
import re
import socket

from abc import ABC, abstractmethod
from base64 import b64encode
from functools import cached_property
from ipaddress import ip_address, ip_network

import requests

from aws_lambda_powertools.utilities import parameters

logger = logging.getLogger()

NO_CONTENT = {"statusCode": 204}


class EventHandler(ABC):
    def __init__(self, config):
        self.config = config
        self._used_up = False
        self.receive_exception = None

    def _handle(self, event):
        assert not self._used_up
        self._used_up = True
        try:
            response = self.receive(event)
            if response is not None:
                return response

        except Exception as exc:
            self.receive_exception = exc
            try:
                response = self.handle_event()
            finally:
                self.receive_exception = None
            logger.error(json_dumps(event.event), exc_info=exc)
            return response
        return self.handle_event()

    @abstractmethod
    def receive(self, event):  # pragma: no cover
        raise NotImplementedError

    @abstractmethod
    def handle_event(self):  # pragma: no cover
        raise NotImplementedError


class StatusRecorder(EventHandler):
    """Update status in Upstash Redis."""

    EVENT_ATTRS = ("unixtime_ms", "source_ip", "worker_id")
    REPORT_TTL = 28 * 24 * 60 * 60

    def receive(self, event):
        self.report = event.body.copy()

        error = self.report.pop("error", None)
        if error is not None:
            self.report.attach_error(error)

        for attr in self.EVENT_ATTRS:
            value = getattr(event, attr, None)
            if value is not None:
                self.report[attr] = value

    @property
    def api_url(self):
        return self.config["upstash_redis_rest_url"]

    @property
    def access_token(self):
        return self.config["upstash_redis_rest_token"]

    def handle_event(self):
        if self.receive_exception is not None:
            self.attach_error({"message": str(self.receive_exception)})

        try:
            miner_status = self.report["miner_status"]
            key = f"{miner_status['id']}-{miner_status['worker_id']}"
            typ = "worker"
        except Exception:
            key, typ = self.report["source_ip"], "host"
        key = f"{typ}:{key}"

        pxat = self.unixtime_ms + self.REPORT_TTL * 1000
        redis_commands = [
            ["set", key, json_dumps(self.report), "PXAT", pxat],
            ["zadd", "reports", self.unixtime, key],
        ]

        response = requests.post(
            f"{self.api_url}/pipeline",
            headers = {  # noqa: E251
                "Content-Type": "application/json",
                "Authorization": f"Bearer {self.access_token}",
            },
            data = json_dumps(redis_commands),  # noqa: E251
        )

        try:
            results = response.json()
            if len(results) != len(redis_commands):
                raise ValueError
            for result in results:
                if result.keys() != {"result"}:
                    raise ValueError
            return NO_CONTENT
        except Exception:
            pass

        log_unhandled_response(response)
        return NO_CONTENT

    @property
    def unixtime_ms(self):
        return self.report["unixtime_ms"]

    @property
    def unixtime(self):
        return self.unixtime_ms / 1000

    def attach_error(self, error):
        if "errors" not in self.report:
            self.report["errors"] = []
        self.report["errors"].append(error)


class MetricsRecorder(EventHandler):
    """Upload metrics to Grafana."""

    def receive(self, event):
        if not event.worker_id:
            return NO_CONTENT

        template = {
            "time": event.unixtime_ms // 1000,
            "interval": 10,
            "tags": (
                f"miner={event.worker_id}",
            ),
        }
        self.metrics = list(self._receive(
            template,
            event.miner_status,
        ))

    @classmethod
    def _receive(cls, template, metrics):
        for key, value in metrics.items():
            if isinstance(value, bool):
                value = int(value)
            if not isinstance(value, (int, float)):
                continue
            metric = template.copy()
            metric["name"] = f"miner.{key}"
            metric["value"] = value
            yield metric

    @property
    def api_url(self):
        return self.config["graphite_api_url"]

    @property
    def access_token(self):
        return self.config["graphite_access_token"]

    def handle_event(self):
        response = requests.post(
            f"{self.api_url}/metrics",
            headers = {  # noqa: E251
                "Content-Type": "application/json",
                "Authorization": f"Bearer {self.access_token}",
            },
            data = json_dumps(self.metrics),  # noqa: E251
        )

        try:
            if response.json()["published"] == len(self.metrics):
                return NO_CONTENT
        except Exception:
            pass

        log_unhandled_response(response)
        return NO_CONTENT


class Event:
    def __init__(self, event, config):
        self.config = config
        self.event = event

    @cached_property
    def request_context(self):
        return self.event["requestContext"]

    @cached_property
    def source_ip(self):
        return self.request_context["identity"]["sourceIp"]

    @cached_property
    def unixtime_ms(self):
        return self.request_context["requestTimeEpoch"]

    @cached_property
    def body(self):
        return json.loads(self.event["body"])

    @cached_property
    def miner_status(self):
        return dict(self._flatten(self.body["miner_status"]))

    @cached_property
    def worker_id(self):
        if "miner_status" not in self.body:
            return None

        result = self.miner_status["worker_id"]
        if not is_in_container_hostname(result):
            if result in self.config.home_hostnames_by_cpu.values():
                return self._qualify_home_hostname(result)
            return result

        result = self.source_ip
        if ip_address(result) not in self.config.home_network:
            return gethostbyaddr(result)

        for word in self.miner_status["cpu.brand"].split():
            name = self.config.home_hostnames_by_cpu.get(word)
            if name is not None:
                return self._qualify_home_hostname(name)

        return gethostbyaddr(result)

    @classmethod
    def _qualify_home_hostname(cls, name, suffix=".nx"):
        if name.endswith(suffix):
            return name
        return f"{name}{suffix}"

    LIST_SUFFIXES = {
        "hashrate.total": ("10s", "1m", "15m"),
        "hugepages": ("got", "want"),
        "resources.load_average": ("1m", "5m", "15m"),
    }

    def _flatten(cls, d, prefix=""):
        if prefix:
            prefix += "."
        for key, value in d.items():
            key = f"{prefix}{key}"
            if isinstance(value, dict):
                for item in cls._flatten(value, key):
                    yield item
                continue

            if isinstance(value, list):
                suffixes = cls.LIST_SUFFIXES.get(key)
                if suffixes is not None:
                    for suffix, value in zip(suffixes, value):
                        yield f"{key}.{suffix}", value
                    continue

            yield key, value


class LambdaHandler:
    EVENT_HANDLERS = (
        StatusRecorder,
        MetricsRecorder,
    )

    def __call__(self, event, context):
        try:
            return self._handle(Event(event, self))
        except Exception as exc:
            logger.error(json_dumps(event), exc_info=exc)
        return NO_CONTENT

    def _handle(self, event):
        responses = []
        for handlerclass in self.EVENT_HANDLERS:
            handler = handlerclass(self.config)
            try:
                responses.append(handler._handle(event))
            except Exception as exc:
                logger.error(json_dumps(event.event), exc_info=exc)
        responses = [r for r in responses if r is not NO_CONTENT]
        if not responses:
            return NO_CONTENT
        if len(responses) == 1:
            return responses[0]
        return lambda_response(200, {"responses": responses})

    @cached_property
    def function_name(self):
        return os.environ.get("AWS_LAMBDA_FUNCTION_NAME")

    @cached_property
    def config(self):
        return json.loads(parameters.get_secret(self.function_name))

    @cached_property
    def home_network(self):
        return ip_network(self.config["home_network"])

    @cached_property
    def home_hostnames_by_cpu(self):
        return json.loads(self.config["home_hostnames_by_cpu"])


def is_in_container_hostname(s):
    """Return `True` if `s` could be the hostname in a container.
    """
    return re.match(r"^[0-9a-f]{12}$", s) is not None


def gethostbyaddr(ip_address):
    try:
        return socket.gethostbyaddr(ip_address)[0]
    except Exception:
        return ip_address


def as_dict(obj):
    """JSON serialization helper.
    """
    d = getattr(obj, "__dict__", None)
    if d is not None:
        return d
    slots = getattr(obj, "__slots__", None)
    if slots is not None:
        return dict((attr, getattr(obj, attr))
                    for attr in obj.__slots__
                    if hasattr(obj, attr))
    if isinstance(obj, bytes):
        try:
            return obj.decode("utf-8")
        except Exception:
            pass
    logger.warning("%s: serializing with repr()", obj)
    return repr(obj)


def json_dumps(obj, separators=(",", ":"), default=as_dict, **kwargs):
    return json.dumps(obj, separators=separators, default=default, **kwargs)


def lambda_response(http_status_code, body, content_type=None):
    if content_type is None:
        content_type = "application/json"
        body = json_dumps(body)

    result = {
        "statusCode": http_status_code,
        "headers": {
            "Content-Type": content_type,
        },
    }

    if isinstance(body, bytes):
        body = b64encode(body)
        result["isBase64Encoded"] = True

    result["body"] = body
    return result


def log_unhandled_response(response, **kwargs):
    if isinstance(response, requests.Response):
        try:
            response = lambda_response(
                response.status_code,
                response.json(),
            )
        except Exception:
            response = lambda_response(
                response.status_code,
                response.content,
                response.headers["content-type"],
            )
    logger.error("%s: unhandled response", json_dumps(response), **kwargs)


lambda_handler = LambdaHandler()
