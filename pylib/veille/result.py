import json
import sys

def problem(message, extra_info={}):
    print json.dumps({
        "status": "problem",
        "message": message,
        "metrics": {},
    })
    sys.exit(0)

def ok(message, extra_info={}):
    print json.dumps({
        "status": "ok",
        "message": message,
        "metrics": {},
    })
    sys.exit(0)
