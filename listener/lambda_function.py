import json

def lambda_handler(event, context):
    cdict = context.__dict__.copy()
    cdict["identity"] = {}
    for attr in context.identity.__class__.__slots__:
        cdict["identity"][attr] = getattr(context.identity, attr)

    return {
        "statusCode": 200,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": json.dumps({
            "event": event,
            "context": cdict,
        })
    }
