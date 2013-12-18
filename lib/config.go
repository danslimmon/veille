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
var confMap map[string]interface{}

var confWatcherExists = false
var confWatcher *ConfigWatcher


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

// Returns the config value identified by the given strings.
//
// This may be a scalar, an array, or a map. It's up to the caller
// to expect the right one.
func GetConfValue(params ...string) (interface{}, error) {
    if confMap == nil {
        err := setConfFromFile(confFilePath)
        if err != nil {
            return nil, ConfigError{"Failed to read config file: " + err.Error()}
        }
    }

    m := confMap
    for _, p := range params[:len(params)-1] {
        nextLevel, ok := m[p].(map[string]interface{})
        if ! ok {
            return nil, ConfigError{"Config tree ends too early at [" + p + "]"}
        }
        if nextLevel == nil {
            return nil, ConfigError{"No config section '" + p + "'"}
        }
        m = nextLevel
    }

    return m[params[len(params)-1]], nil
}

// Reads and parses the given config file, barfing if it's invalid.
func LoadConfig(path string) error {
    if ! confWatcherExists {
        confWatcher = new(ConfigWatcher)
        confWatcher.RegisterSignals()
    }
    confWatcherExists = true

    err := setConfFromFile(path)
    if err != nil { return err }

    return nil
}

// Sets the global config map to the value read from the given YAML file.
func setConfFromFile(path string) error {
    yBlob, err := ioutil.ReadFile(path)
    if err != nil { return err }
    err = goyaml.Unmarshal(yBlob, &confMap)
    if err != nil { return err }
    return nil
}
