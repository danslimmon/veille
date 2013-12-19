package veille

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/signal"
    "syscall"

    "launchpad.net/goyaml"
)

var confFilePath string
var conf Config

var confWatcherExists = false
var confWatcher *ConfigWatcher

type Config struct {
    Services []ServiceConfig
}
type ServiceConfig struct {
    Service_Name string
    Tests []TestConfig
}
type TestConfig struct {
    Functionality string
    Script string
    Run_Every int
    Alert_After int
    Alert AlertConfig
}
type AlertConfig struct {
    Mode string
    Target string
}


type ConfigError struct { msg string }
func (e ConfigError) Error() string { return "Config error: " + e.msg }

// Sends notifications on channels whenever the config changes.
type ConfigWatcher struct {
    OutputChans []chan byte
}
func (cw *ConfigWatcher) Subscribe() chan byte {
    ch := make(chan byte)
    cw.OutputChans = append(cw.OutputChans, ch)
    return ch
}
func (cw *ConfigWatcher) Publish() {
    fmt.Println("Notifying", len(cw.OutputChans), "goroutines of config reload")
    for _, ch := range cw.OutputChans {
        go cw.write(ch)
    }
}
func (cw *ConfigWatcher) RegisterSignals() {
    ch := make(chan os.Signal)
    go func() {
        for {
            <-ch
            cw.Publish()
        }
    }()
    signal.Notify(ch, syscall.SIGHUP)
}
func (cw *ConfigWatcher) write(outputChan chan byte) {
    outputChan <- 0
}

func ConfigSubscribe() chan byte {
    return confWatcher.Subscribe()
}

// Returns the Config struct that's been built by LoadConfig.
func GetConfig() *Config {
    return &conf
}

// Reads and parses the given config file, barfing if it's invalid.
func LoadConfig(path string) error {
    if ! confWatcherExists {
        confWatcher = new(ConfigWatcher)
        confWatcher.RegisterSignals()
    }
    confWatcherExists = true

    err := setConfigFromFile(path)
    if err != nil { return err }

    return nil
}

// Sets the global config map to the value read from the given YAML file.
func setConfigFromFile(path string) error {
    yBlob, err := ioutil.ReadFile(path)
    if err != nil { return err }
    err = goyaml.Unmarshal(yBlob, &conf)
    if err != nil { return err }
    return nil
}
