package veille

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/signal"
    "syscall"
    "log"

    "launchpad.net/goyaml"
)

type ConfigLoader interface {
    GetConfig() (*Config, error)
    ReloadConfig() (*Config, error)
}

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
    Loader ConfigLoader
    outputChans []chan *Config
}
func (cw *ConfigWatcher) Subscribe() chan *Config {
    ch := make(chan *Config)
    cw.outputChans = append(cw.outputChans, ch)
    return ch
}
func (cw *ConfigWatcher) Publish(c *Config) {
    log.Println("Notifying", len(cw.outputChans), "goroutines of config reload")
    for _, ch := range cw.outputChans {
        go func() {ch <- c}()
    }
}
func (cw *ConfigWatcher) RegisterSignals() {
    ch := make(chan os.Signal)
    go func() {
        for {
            <-ch
            conf, err := cw.Loader.GetConfig()
            if err != nil {
                fmt.Println("Received a SIGHUP, but failed to parse config: " + err.Error())
            }
            cw.Publish(conf)
        }
    }()
    signal.Notify(ch, syscall.SIGHUP)
}

// A ConfigLoader that loads from a YAML file.
//
// Must be initialized with SetPath() before you can load anything.
type YamlFileConfigLoader struct {
    Path string
    cachedConfig *Config
}

// Returns the active configuration.
func (loader *YamlFileConfigLoader) GetConfig() (*Config, error) {
    if loader.cachedConfig == nil {
        cc, err := loader.parseFile(loader.Path)
        if err != nil { return nil, err }
        loader.cachedConfig = cc
    }
    return loader.cachedConfig, nil
}

// Returns the active configuration.
func (loader *YamlFileConfigLoader) ReloadConfig() (*Config, error) {
    return nil, nil
}

// Reads the YAML file and returns the parsed Config.
func (loader *YamlFileConfigLoader) parseFile(path string) (*Config, error) {
    conf := new(Config)
    yBlob, err := ioutil.ReadFile(path)
    if err != nil { return nil, err }
    err = goyaml.Unmarshal(yBlob, conf)
    if err != nil { return nil, err }
    return conf, nil
}
