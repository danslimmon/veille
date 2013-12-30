package veille

import (
    "testing"

    "os"
    "syscall"
    "io/ioutil"
)

func TestYamlFileConfigLoader_GetConfig(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    confStr := `---
services:
  - service_name: "Github API"
    tests:
      - functionality: "List a user's projects"
        script: "github_user_repos.py"
        run_every: 20
        alert_after: 3
        alert:
            mode: "email"
            target: "developers@example.com"`

    // Write test configuration file
    confFile, err := ioutil.TempFile("", "veille_test_")
    if err != nil { t.FailNow() }
    confFile.Write([]byte(confStr))
    confPath := confFile.Name()
    defer os.Remove(confPath)
    confFile.Close()

    // Load the config from that temp file
    loader := &YamlFileConfigLoader{
        Path: confPath,
    }
    c, err := loader.GetConfig()
    if err != nil { t.Log("Error loading YAML configuration: " + err.Error()) }

    // Make sure that the config was correctly parsed.
    switch false {
    case c.Services[0].Service_Name == "Github API":
        t.Log("Failed to load service name from configuration")
        t.Fail()
    case c.Services[0].Tests[0].Functionality == "List a user's projects":
        t.Log("Failed to load test name from configuration")
        t.Fail()
    case c.Services[0].Tests[0].Script == "github_user_repos.py":
        t.Log("Failed to load test script from configuration")
        t.Fail()
    case c.Services[0].Tests[0].Run_Every == 20:
        t.Log("Failed to load test run interval from configuration")
        t.Fail()
    case c.Services[0].Tests[0].Alert_After == 3:
        t.Log("Failed to load test failure threshold from configuration")
        t.Fail()
    case c.Services[0].Tests[0].Alert.Mode == "email":
        t.Log("Failed to load test alert mode from configuration")
        t.Fail()
    case c.Services[0].Tests[0].Alert.Target == "developers@example.com":
        t.Log("Failed to load test alert target")
        t.Fail()
    }
}


func TestConfigWatcher(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    conf := &Config{[]ServiceConfig{}}
    confWatcher := new(ConfigWatcher)
    ch := confWatcher.Subscribe()
    confWatcher.Publish(conf)

    newConf := <- ch
    if conf != newConf {
        t.Log("ConfigWatcher didn't return the new Config object")
        t.Fail()
    }
}


func TestConfigWatcher_PublishOnSignals(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    conf := &Config{[]ServiceConfig{}}
    loader := &MockConfigLoader{
        Config: conf,
    }
    confWatcher := new(ConfigWatcher)
    confWatcher.Loader = loader

    subCh := confWatcher.Subscribe()
    sigCh := confWatcher.PublishOnSignals()

    sigCh <- syscall.SIGHUP
    newConf := <- subCh
    if conf != newConf {
        t.Log("ConfigWatcher didn't return the new Config object in response to signal")
        t.Fail()
    }
}
