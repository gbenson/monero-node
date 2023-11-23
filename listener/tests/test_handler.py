import json


def test_handler(lambda_handler, startup_event_3):
    event = startup_event_3
    result = lambda_handler(event, None)
    print(result["body"] if "body" in result else result)
    assert result["statusCode"] == 204


def test_handler_no_status(lambda_handler, startup_event_3):
    event = startup_event_3
    body = json.loads(event["body"])
    body.pop("miner_status")
    event["body"] = json.dumps(body)
    result = lambda_handler(event, None)
    print(result["body"] if "body" in result else result)
    assert result["statusCode"] == 204


def test_handler_named_worker(lambda_handler, startup_event_3):
    event = startup_event_3
    body = json.loads(event["body"])
    body["miner_status"]["worker_id"] = "miner1"
    event["body"] = json.dumps(body)
    result = lambda_handler(event, None)
    print(result["body"] if "body" in result else result)
    assert result["statusCode"] == 204


def test_handler_home_worker(lambda_handler, startup_event_3):
    event = startup_event_3
    lambda_handler.func.config["home_network"] = (
        event["requestContext"]["identity"]["sourceIp"])
    result = lambda_handler(event, None)
    print(result["body"] if "body" in result else result)
    assert result["statusCode"] == 204


def test_handler_unknown_home_worker(lambda_handler, startup_event_3):
    event = startup_event_3
    lambda_handler.func.config["home_network"] = (
        event["requestContext"]["identity"]["sourceIp"])
    body = json.loads(event["body"])
    body["miner_status"]["cpu"]["brand"] = (
        "Intel(R) Xeon(R) CPU E5-2620 v4 @ 2.10GHz")
    event["body"] = json.dumps(body)
    result = lambda_handler(event, None)
    print(result["body"] if "body" in result else result)
    assert result["statusCode"] == 204
