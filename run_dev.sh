#!/bin/bash

export VEILLE_LISTEN=:8080
export AWS_ACCESS_KEY_ID=dummy
export AWS_SECRET_ACCESS_KEY=dummy
exec ./veille
