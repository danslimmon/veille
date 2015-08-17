package main

import (
	"time"
)

// A State implementation represents host or service state as gleaned
// from a Nagios log entry.
type State interface {
	// Timestamp returns the timestamp at which the state was recorded.
	Timestamp() time.Time
	// ObjectType returns the type of Nagios object whose state is
	// recorded. This is either "HOST" or "SERVICE".
	ObjectType() string
	// Hostname returns the hostname of the host or service whose state
	// is recorded.
	Hostname() string
	// Servicename returns the name of the service whose state is
	// recorded. If it's a HOST status rather than a SERVICE state,
	// this returns the empty string.
	Servicename() string
	// Status returns the Nagios status that we've recorded. That's
	// usually "OK", "CRITICAL", or "WARNING" for services and "UP" or
	// "DOWN" for hosts.
	Status() string
	// Hardness returns the Nagios "hardness" of the status we've
	// recorded. That's either "SOFT" or "HARD".
	Hardness() string
	// PluginOutput returns the output of the most recent check for
	// the Nagios object.
	PluginOutput() string
}

// HostState contains information about the state of a host object
// in Nagios. It implements the Status interface.
type HostState struct {
	timestamp    time.Time
	hostname     string
	status       string
	hardness     string
	pluginOutput string
}

func (hs HostState) Timestamp() time.Time { return hs.timestamp }
func (hs HostState) ObjectType() string   { return "HOST" }
func (hs HostState) Hostname() string     { return hs.hostname }
func (hs HostState) Servicename() string  { return "" }
func (hs HostState) Status() string       { return hs.status }
func (hs HostState) Hardness() string     { return hs.hardness }
func (hs HostState) PluginOutput() string { return hs.pluginOutput }

func NewHostState(timestamp time.Time, hostname, status, hardness, pluginOutput string) HostState {
	return HostState{timestamp, hostname, status, hardness, pluginOutput}
}
