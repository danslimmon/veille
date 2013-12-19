#!/usr/bin/python

import json
import urllib2 as urllib

import veille

# Do the HTTP request
rsp = urllib.urlopen("https://api.github.com/users/danslimmon/repos")
if rsp.getcode() != 200:
    veille.problem("Received HTTP {0} response to Github API request".format(rsp.getcode()))

# Parse the JSON response
try:
    rsp_body = rsp.body()
    rsp_obj = json.load(rsp_body)
except e:
    veille.problem("Got error {0} when trying to parse Github API request".format(e),
                   extra_info={"github_response_body": rsp_body})

# Make sure Veille is listed in the response
if not [repo for repo in rsp_obj if rsp_obj["name"] == "veille"]:
    veille.problem("No project named 'veille' in the Github response",
                   extra_info={"github_response_body": rsp_body})

veille.ok()
