#!/bin/bash

# https://github.com/wg/wrk
# 12 threads
# 400 connections
# 2m duration
wrk -t 12 -c 400 -d 2m --latency "http://api.icndb.com/jokes/random?firstName=__NAME__&lastName=__SURNAME__&limitTo=[nerdy]" > icndb-wrk-$(date +%s).txt
