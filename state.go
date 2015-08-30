package main

import (
	"fmt"
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

	// ObjIdent returns a string that uniquely identifies the Nagios
	// object that this state describes. Useful for hash-keying.
	ObjIdent() string
}

// HostState contains information about the state of a host object
// in Nagios. It implements the State interface.
type HostState struct {
	timestamp    time.Time
	hostname     string
	status       string
	hardness     string
	pluginOutput string
}

func (st HostState) Timestamp() time.Time { return st.timestamp }
func (st HostState) ObjectType() string   { return "HOST" }
func (st HostState) Hostname() string     { return st.hostname }
func (st HostState) Servicename() string  { return "" }
func (st HostState) Status() string       { return st.status }
func (st HostState) Hardness() string     { return st.hardness }
func (st HostState) PluginOutput() string { return st.pluginOutput }
func (st HostState) ObjIdent() string     { return st.hostname }

func NewHostState(timestamp time.Time, hostname, status, hardness, pluginOutput string) HostState {
	return HostState{timestamp, hostname, status, hardness, pluginOutput}
}

// ServiceState contains information about the state of a service object
// in Nagios. It implements the State interface.
type ServiceState struct {
	timestamp    time.Time
	hostname     string
	servicename  string
	status       string
	hardness     string
	pluginOutput string
}

func (st ServiceState) Timestamp() time.Time { return st.timestamp }
func (st ServiceState) ObjectType() string   { return "SERVICE" }
func (st ServiceState) Hostname() string     { return st.hostname }
func (st ServiceState) Servicename() string  { return st.servicename }
func (st ServiceState) Status() string       { return st.status }
func (st ServiceState) Hardness() string     { return st.hardness }
func (st ServiceState) PluginOutput() string { return st.pluginOutput }
func (st ServiceState) ObjIdent() string {
	return fmt.Sprintf("%s;%s", st.hostname, st.servicename)
}

func NewServiceState(timestamp time.Time, hostname, servicename, status,
	hardness, pluginOutput string) ServiceState {
	return ServiceState{timestamp, hostname, servicename, status, hardness, pluginOutput}
}
