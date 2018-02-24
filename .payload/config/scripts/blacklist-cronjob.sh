#!/bin/bash
# Cron script runs update-dnsmasq at random times with a 4 hour window

sleep $(( RANDOM % 14400 )) && /config/scripts/update-dnsmasq