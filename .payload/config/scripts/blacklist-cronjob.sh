#!/bin/bash
# Cron script runs update-dnsmasq at random times within a 3 hour window

sleep $(( RANDOM \% 10800 )); /config/scripts/update-dnsmasq