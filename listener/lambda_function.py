import json
import logging
import os
import re

from functools import cached_property

import requests

from aws_lambda_powertools.utilities import parameters

logger = logging.getLogger()

NO_CONTENT = {"statusCode": 204}
DOCKER_WORKER_ID = re.compile(r"^[0-9a-f]{12}$")


class Handler:
    def __call__(self, event, context):
        try:
            result = self.handle(Event(event))
            if result is not None:
                return result

        except Exception as e:
            logger.error(json_dumps(event), exc_info=e)

        return NO_CONTENT

    def handle(self, event):
        response = requests.post(
            f"{self.config['graphite_api_url']}/metrics",
            headers = {
                "Content-Type": "application/json",
                "Authorization": f"Bearer {self.config['graphite_access_token']}",
            },
            data = json_dumps(event.metrics),
        )

        return {
            "statusCode": 200,
            "headers": {
                "Content-Type": "application/json"
            },
            "body": json_dumps(response.json()),
        }

    @cached_property
    def function_name(self):
        return os.environ.get("AWS_LAMBDA_FUNCTION_NAME")

    @cached_property
    def config(self):
        return json.loads(parameters.get_secret(self.function_name))


class Event:
    __slots__ = "context", "error", "miner_api", "miner_status"

    def __init__(self, src):
        for attr, value in json.loads(src["body"]).items():
            setattr(self, attr, value)
        self.context = EventContext(src["requestContext"])

    @property
    def worker_id(self):
        worker_id = self.miner_status["worker_id"]
        if DOCKER_WORKER_ID.match(worker_id) is None:
            return worker_id
        return self.context.source_ip

    @property
    def metrics(self):
        return self._descend_metrics(self.miner_status, {
                "name": "miner",
                "interval": 10,
                "time": self.context.unixtime_ms // 1000,
                "tags": [f"miner={self.worker_id}"],
        })

    SUFFIXES = {
        "miner.hashrate.total": ("10s", "1m", "15m"),
        "miner.hugepages": ("got", "want"),
        "miner.resources.load_average": ("1m", "5m", "15m"),
    }

    @classmethod
    def _descend_metrics(cls, src, base, dst=None):
        if dst is None:
            dst = []
        for name, value in src.items():
            if value is None:
                continue
            if isinstance(value, (str, bool)):
                continue

            name = f"{base['name']}.{name}"
            if hasattr(value, "items"):
                child = base.copy()
                child["name"] = name
                cls._descend_metrics(value, child, dst)
                continue

            if not isinstance(value, list):
                metric = base.copy()
                metric["name"] = name
                metric["value"] = value
                dst.append(metric)
                continue

            if not value:
                continue
            if isinstance(value[0], str):
                continue

            suffixes = cls.SUFFIXES.get(name)
            if suffixes is not None:
                for suffix, value in zip(suffixes, value):
                    if value is None:
                        continue
                    metric = base.copy()
                    metric["name"] = f"{name}.{suffix}"
                    metric["value"] = value
                    dst.append(metric)
                continue

            if name != "miner.results.best":
                logger.warning("%s?", name)
                continue

        return dst


class EventContext:
    def __init__(self, src):
        self.source_ip = src["identity"]["sourceIp"]
        self.unixtime_ms = src["requestTimeEpoch"]


def as_dict(obj):
    """JSON serialization helper"""
    d = getattr(obj, "__dict__", None)
    if d is not None:
        return d
    return dict((attr, getattr(obj, attr))
                for attr in obj.__slots__
                if hasattr(obj, attr))


def json_dumps(obj, separators=(",", ":"), default=as_dict, **kwargs):
    return json.dumps(obj, separators=separators, default=default, **kwargs)


lambda_handler = Handler()
