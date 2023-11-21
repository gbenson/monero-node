from lambda_function import lambda_handler


def test_handler(startup_event_1):
    event = startup_event_1
    result = lambda_handler(event, None)
    print(result["body"] if "body" in result else result)
    assert False
