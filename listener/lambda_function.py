import json
import random

from base64 import b64decode, b64encode

from aws_lambda_powertools.utilities import parameters


def lambda_handler(event, context):
    graphite_api = json.loads(parameters.get_secret("graphite_api"))

    token = graphite_api.pop("access_token")
    prefix, data = token.split("_")
    creds = list(b64decode(data))
    random.shuffle(creds)
    creds = bytearray(creds)
    data = b64encode(creds).decode("ascii")
    token = "_".join((prefix, data))
    graphite_api["access_token"] = token

    return {
        "statusCode": 200,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": json.dumps({
            "event": event,
            "graphite_api": graphite_api,
        })
    }
