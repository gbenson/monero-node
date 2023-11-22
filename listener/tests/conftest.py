import json
import logging
import os

import pytest

from lambda_function import LambdaHandler

TESTDIR = os.path.dirname(__file__)


class FailOnLog:
    def __init__(self, caplog, func, level=logging.WARNING):
        self.func = func
        self.caplog = caplog
        self.level = level

    def __call__(self, *args, **kwargs):
        with self.caplog.at_level(self.level):
            result = self.func(*args, **kwargs)
        assert not self.caplog.records
        return result


@pytest.fixture
def lambda_handler(caplog):
    yield FailOnLog(caplog, LambdaHandler())


def load_json(filename):
    with open(os.path.join(TESTDIR, "resources", filename)) as fp:
        return json.load(fp)


@pytest.fixture
def startup_event_1():
    return load_json("startup-event-1.json")


@pytest.fixture
def startup_event_3():
    return load_json("startup-event-2.json")
