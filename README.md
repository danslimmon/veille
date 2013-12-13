guet
====

Guet is an event-based monitoring and alerting tool with a focus on
signal-to-noise ratio. Here are its goals:

Signal-to-noise ratio aware
----

* Requests and processes user feedback about whether an alert was real
* Allows you to sort pobes by their empirical [Positive predictive value][ppv]
* Can email when a probe’s PPV gets too low
* Can detect and notify when two probes seem coupled/correlated

Minimal but sufficient set of alerts
----

* Separates detection from diagnosis
* Presents diagnostic results to the user during an incident

Hysteresis
----

* Probes can do basic statistical analysis on historical data
* Probes can integrate with [Graphite][graphite]
* Probes can integrate with [Elasticsearch][es]

Other nice features
----

* Provides searchable log of incidents with their full diagnostic history
* Publishes incident response metrics
* Can dynamically update the list of hosts on which to run diagnostics for a particular service


[ppv]: http://en.wikipedia.org/wiki/Positive_predictive_value
[graphite]: http://graphite.wikidot.com/
[es]: http://www.elasticsearch.org/
