import json

from lambda_function import lambda_handler


def test_handler(startup_event_3):
    event = startup_event_3
    result = lambda_handler(event, None)
    print(result["body"] if "body" in result else result)


def test_handler_no_status(startup_event_3):
    event = startup_event_3
    body = json.loads(event["body"])
    body.pop("miner_status")
    event["body"] = json.dumps(body)
    result = lambda_handler(event, None)
    print(result["body"] if "body" in result else result)
