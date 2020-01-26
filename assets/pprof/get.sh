#!/bin/bash

wget "http://localhost:6060/debug/pprof/profile?seconds=$1" -O pprof_$(date +%s)

