import json
import os

def lambda_handler(event, context):
    return {
        "statusCode": 200,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": json.dumps({
            "message": "Hello world!",
            "event": event,
            "environ": os.environ,
        })
    }
