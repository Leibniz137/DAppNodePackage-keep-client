#!/bin/sh

./prepare

exec keep-client -config ./config.toml start
