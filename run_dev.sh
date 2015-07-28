#!/bin/bash

export VEILLE_INPUT_MODE=tail
export VEILLE_LOG_LEVEL=debug

go get
go build

exec ./veille
