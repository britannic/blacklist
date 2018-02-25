#!/bin/bash
# Cron script runs update-dnsmasq at random times within a 3 hour window
 
seconds=${1}

sleep $(( RANDOM \% ${seconds} ))
/config/scripts/update-dnsmasq
