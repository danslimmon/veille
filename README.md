veille
====

Veille is an event-based monitoring and alerting tool with a focus on
signal-to-noise ratio. Here are its goals:

Signal-to-noise ratio aware
----

* Requests and processes user feedback about whether an alert was real
* Allows you to sort pobes by their empirical [Positive predictive value][ppv]
* Can email when a test’s PPV gets too low
* Can detect and notify when two tests seem coupled/correlated

Minimal but sufficient set of alerts
----

* Separates detection from diagnosis
* Presents diagnostic results to the user during an incident

Hysteresis
----

* Tests can do basic statistical analysis on historical data
* Tests can integrate with [Graphite][graphite]
* Tests can integrate with [Elasticsearch][es]

Other nice features
----

* Provides searchable log of incidents with their full diagnostic history
* Publishes incident response metrics
* Can dynamically update the list of hosts on which to run diagnostics for a particular service


Configuration
=====

Veille's focus is on whether _work is getting done_. That's what you really care
about, right? Why test whether Apache is running, or whether you can ping the web
server, if you already know that your service is responding to requests?

In a Veille configuration you define __services__. A service is something that
people -- or other services -- need to be working. And for each service, you
define one or more __tests__. Each test makes sure that a given piece of
functionality is working.

A Simple Test
-----

Let's walk through the process of configuring Veille. We'll populate it with a
simple test that makes sure
[Github's API] [http://developer.github.com/v3/repos/#list-user-repositories]
is correctly returning info about the Veille project. This is a good example of
an end-to-end test, which is the bread and butter of Veille.

We'll start with a script that executes the actual test: making an HTTP GET
request for `http://www.example.com` and making sure that we get back a 200
response with the appropriate content. We'll implement it in Python 2 in the file
`tests/github_user_repos.py` (sure, you can write Veille tests in any language you
want, but right now there's only a library for Python).

```python
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
```

Make sure this file is executable by the user `veille`.

Now that we have a test, we will tell Veille to run it every 20 seconds. Edit
the file `veille.conf`, find the `services` section, and modify it thusly:

```yaml
---
services:
  - service_name: "Github API"
    tests:
      - functionality: "List a user's projects"
        script: github_user_repos.py
        run_every: 20
        alert_after: 3
        alert:
            mode: email
            target: developers@example.com
```

We've created a service called "Github API" with a single test. This test runs
the script `github_user_repos.py` every 20 seconds. If 3 tests in a row fail,
then Veille will email `developers@example.com` with an alert.

Let's reload Veille's config and watch it execute:

```
you@server:~/$ sudo /etc/init.d/veille reload
you@server:~/$ tail -F /var/log/veille.log
```


[ppv]: http://en.wikipedia.org/wiki/Positive_predictive_value
[graphite]: http://graphite.wikidot.com/
[es]: http://www.elasticsearch.org/
