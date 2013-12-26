import json
import sys

def problem(message, extra_info={}):
    print json.dumps({
        "status": "problem",
        "message": message,
        "extra_info": extra_info,
        "metrics": [],
    })
    sys.exit(0)

def ok(message, extra_info={}):
    print json.dumps({
        "status": "ok",
        "message": "message",
        "extra_info": extra_info,
        "metrics": [],
    })
    sys.exit(0)
