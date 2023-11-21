import json
import os

import pytest

TESTDIR = os.path.dirname(__file__)


def load_json(filename):
    with open(os.path.join(TESTDIR, "resources", filename)) as fp:
        return json.load(fp)


@pytest.fixture
def startup_event_1():
    return load_json("startup-event-1.json")
