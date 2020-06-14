#!/bin/sh

set -e

./prepare

exec keep-client -config ./config.toml start
