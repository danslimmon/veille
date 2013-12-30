package veille

import (
    "testing"
    "os"
    "io/ioutil"
)

func TestTest_Check(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    f, e := ioutil.TempFile("", "veille_test_")
    if e != nil {
        t.Log("Error creating temp file:", e)
        t.FailNow()
    }
    scriptPath := f.Name()
    defer os.Remove(scriptPath)
    f.Write([]byte(`#!/bin/bash
echo '{"status":"ok","message":"message","metrics":{}}'
`))
    f.Chmod(0755)
    f.Close()

    test := Test{
        Functionality: "Some arbitrary action works",
        Script: scriptPath,
    }
    rslt := test.Check()

    switch false {
    case rslt.Status == "ok":
        t.Log("Got wrong status from test script")
        t.FailNow()
    case rslt.Message == "message":
        t.Log("Got wrong message from test script")
        t.FailNow()
    case rslt.T.Functionality == test.Functionality:
        t.Log("Test result populated with wrong test struct")
        t.FailNow()
    }
}
