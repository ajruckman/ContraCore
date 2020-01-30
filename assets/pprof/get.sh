#!/bin/bash

#wget "http://10.2.0.105:6053/debug/pprof/profile?seconds=$1" -O pprof_$(date +%s)
wget "http://10.2.0.105:6053/debug/pprof/trace?seconds=$1" -O pprof_$(date +%s)

