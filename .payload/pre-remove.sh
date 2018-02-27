#!/usr/bin/env bash

# Set up the Vyatta environment
declare -i DEC
source /opt/vyatta/etc/functions/script-template
DATE=$(date +'%FT%H%M%S')
shopt -s expand_aliases

alias bold='tput bold'
alias normal='tput sgr0'
alias reverse='tput smso'
alias underline='tput smul'

alias black='tput setaf 0'
alias blink='tput blink'
alias blue='tput setaf 4'
alias cyan='tput setaf 6'
alias green='tput setaf 2'
alias lime='tput setaf 190'
alias magenta='tput setaf 5'
alias powder='tput setaf 153'
alias purple='tput setaf 171'
alias red='tput setaf 1'
alias tan='tput setaf 3'
alias white='tput setaf 7'
alias yellow='tput setaf 3'

# Setup the echo_logger function
echo_logger() {
	local MSG
	shopt -s checkwinsize
	COLUMNS=$(tput cols)
	DEC+=1
	CTR=$(printf "%03x" ${DEC})
	TIME=$(date +%H:%M:%S.%3N)

	case "${1}" in
	E)
		shift
		MSG="$(red)$(bold)ERRO$(normal)[${CTR}]${TIME}: ${@} failed!"
		;;
	F)
		shift
		MSG="$(red)$(bold)FAIL$(normal)[${CTR}]${TIME}: ${@}"
		;;
	FE)
		shift
		MSG="$(red)$(bold)CRIT$(normal)[${CTR}]${TIME}: ${@}"
		;;
	I)
		shift
		MSG="$(green)INFO$(normal)[${CTR}]${TIME}: ${@}"
		;;
	S)
		shift
		MSG="$(green)$(bold)NOTI$(normal)[${CTR}]${TIME}: ${@}"
		;;
	T)
		shift
		MSG="$(tan)$(bold)TRYI$(normal)[${CTR}]${TIME}: ${@}"
		;;
	W)
		shift
		MSG="$(yellow)$(bold)WARN$(normal)[${CTR}]${TIME}: ${@}"
		;;
	*)
		echo "ERROR: usage: echo_logger MSG TYPE(E, F, FE, I, S, T, W) MSG."
		exit 1
		;;
	esac

	# MSG=$(echo "${MSG}" | ansi)
	let COLUMNS=${#MSG}-${#@}+${COLUMNS}
	echo "pre-upgrade: ${MSG}" | fold -sw ${COLUMNS}
}

# Fix the group so that the admin user will be able to commit configs
set_vyattacfg_grp() {
	try chgrp -R vyattacfg /opt/vyatta/config/active
}

# Function to output command status of success or failure to screen and log
try() {
	if eval "${@}"; then
		echo_logger I "${@}"
		return 0
	else
		echo_logger E "${@}"
		return 1
	fi
}

# Back up [service dns forwarding blacklist]
backup_dns_config() {
	/bin/cli-shell-api existsActive service dns forwarding blacklist
	if [[ $? == 0 ]]; then
		echo_logger I "Backing up blacklist configuration to: /config/user-data/blacklist.${DATE}.cmds"
		run show configuration commands | grep blacklist >/config/user-data/blacklist.${DATE}.cmds ||
			echo_logger E 'Blacklist configuration backup failed!'
	fi
}

backup_dns_config