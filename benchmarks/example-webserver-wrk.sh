#!/bin/bash

# https://github.com/wg/wrk
# 12 threads
# 400 connections
# 2m duration
wrk -t 12 -c 400 -d 2m --latency "http://localhost:5000/v1/" > example-webserver-wrk-$(date +%s).txt
