#!/bin/bash

# https://github.com/wg/wrk
# 12 threads
# 400 connections
# 2m duration
wrk -t 4 -c 40 -d 2m --latency "http://uinames.com/api/" > uinames-wrk-$(date +%s).txt
