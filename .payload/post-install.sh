#!/usr/bin/env bash

# Set up the Vyatta environment
declare -i DEC
source /opt/vyatta/etc/functions/script-template
OPRUN=/opt/vyatta/bin/vyatta-op-cmd-wrapper
CFGRUN=/opt/vyatta/sbin/vyatta-cfg-cmd-wrapper
API=/bin/cli-shell-api
shopt -s expand_aliases

alias AddImage='${OPRUN} add system image'
alias begin='${CFGRUN} begin'
alias cleanup='${CFGRUN} cleanup'
alias comment='${CFGRUN} comment'
alias commit='${CFGRUN} commit'
alias copy='${CFGRUN} copy'
alias delete='${CFGRUN} delete'
alias DeleteImage='/usr/bin/ubnt-upgrade --delete-noprompt'
alias discard='${CFGRUN} discard'
alias end='${CFGRUN} end'
alias load='${CFGRUN} load'
alias rename='${CFGRUN} rename'
alias save='sudo ${CFGRUN} save'
alias set='${CFGRUN} set'
alias show='${API} showConfig'
alias showddns=/opt/vyatta/bin/sudo-users/vyatta-op-dynamic-dns.pl
alias version='${OPRUN} show version'

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
	CTR=$( printf "%03x" ${DEC} )
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
	echo "post-install: ${MSG}" | fold -sw ${COLUMNS}
}

# Fix the group so that the admin user will be able to commit configs
set_vyattacfg_grp() {
	try chgrp -R vyattacfg /opt/vyatta/config
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

# Load the [service dns forwarding blacklist] configuration
update_dns_config() {
	try begin
	try set service dns forwarding blacklist disabled false
	try set service dns forwarding blacklist dns-redirect-ip 0.0.0.0
	try set service dns forwarding blacklist domains include adk2x.com
	try set service dns forwarding blacklist domains include adsrvr.org
	try set service dns forwarding blacklist domains include adtechus.net
	try set service dns forwarding blacklist domains include advertising.com
	try set service dns forwarding blacklist domains include centade.com
	try set service dns forwarding blacklist domains include doubleclick.net
	try set service dns forwarding blacklist domains include fastplayz.com
	try set service dns forwarding blacklist domains include free-counter.co.uk
	try set service dns forwarding blacklist domains include hilltopads.net
	try set service dns forwarding blacklist domains include intellitxt.com
	try set service dns forwarding blacklist domains include kiosked.com
	try set service dns forwarding blacklist domains include patoghee.in
	try set service dns forwarding blacklist domains include themillionaireinpjs.com
	try set service dns forwarding blacklist domains include traktrafficflow.com
	try set service dns forwarding blacklist domains include wwwpromoter.com
	try set service dns forwarding blacklist domains source githubSteveBlack url 'https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews-gambling-porn-social/hosts'
	try set service dns forwarding blacklist domains source githubSteveBlack description '"Blacklists fakenews, gambling, porn and social websites"'
	try set service dns forwarding blacklist domains source githubSteveBlack prefix '0.0.0.0 '
	try set service dns forwarding blacklist domains source malc0de description '"List of zones serving malicious executables observed by malc0de.com/database/"'
	try set service dns forwarding blacklist domains source malc0de prefix 'zone '
	try set service dns forwarding blacklist domains source malc0de url 'http://malc0de.com/bl/ZONES'
	try set service dns forwarding blacklist domains source malwaredomains.com description '"Just Domains"'
	try set service dns forwarding blacklist domains source malwaredomains.com url 'http://mirror1.malwaredomains.com/files/justdomains'
	try set service dns forwarding blacklist domains source simple_tracking description '"Basic tracking list by Disconnect"'
	try set service dns forwarding blacklist domains source simple_tracking url 'https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt'
	try set service dns forwarding blacklist domains source zeus description '"abuse.ch ZeuS domain blocklist"'
	try set service dns forwarding blacklist domains source zeus url 'https://zeustracker.abuse.ch/blocklist.php?download=domainblocklist'
	try set service dns forwarding blacklist domains source NoBitCoin description '"Blocking Web Browser Bitcoin Mining"'
	try set service dns forwarding blacklist domains source NoBitCoin prefix '0.0.0.0'
	try set service dns forwarding blacklist domains source NoBitCoin url 'https://raw.githubusercontent.com/hoshsadiq/adblock-nocoin-list/master/hosts.txt'	
	try set service dns forwarding blacklist exclude 1e100.net
	try set service dns forwarding blacklist exclude 2o7.net
	try set service dns forwarding blacklist exclude adobedtm.com
	try set service dns forwarding blacklist exclude akamai.net
	try set service dns forwarding blacklist exclude akamaihd.net
	try set service dns forwarding blacklist exclude amazon.com
	try set service dns forwarding blacklist exclude amazonaws.com
	try set service dns forwarding blacklist exclude apple.com
	try set service dns forwarding blacklist exclude ask.com
	try set service dns forwarding blacklist exclude avast.com
	try set service dns forwarding blacklist exclude avira-update.com
	try set service dns forwarding blacklist exclude bannerbank.com
	try set service dns forwarding blacklist exclude bing.com
	try set service dns forwarding blacklist exclude bit.ly
	try set service dns forwarding blacklist exclude bitdefender.com
	try set service dns forwarding blacklist exclude cdn.ravenjs.com
	try set service dns forwarding blacklist exclude cdn.visiblemeasures.com
	try set service dns forwarding blacklist exclude cloudfront.net
	try set service dns forwarding blacklist exclude coremetrics.com
	try set service dns forwarding blacklist exclude dropbox.com
	try set service dns forwarding blacklist exclude ebay.com
	try set service dns forwarding blacklist exclude edgesuite.net
	try set service dns forwarding blacklist exclude evernote.com
	try set service dns forwarding blacklist exclude feedly.com
	try set service dns forwarding blacklist exclude freedns.afraid.org
	try set service dns forwarding blacklist exclude github.com
	try set service dns forwarding blacklist exclude githubusercontent.com
	try set service dns forwarding blacklist exclude global.ssl.fastly.net
	try set service dns forwarding blacklist exclude google.com
	try set service dns forwarding blacklist exclude googleads.g.doubleclick.net
	try set service dns forwarding blacklist exclude googleadservices.com
	try set service dns forwarding blacklist exclude googleapis.com
	try set service dns forwarding blacklist exclude googletagmanager.com
	try set service dns forwarding blacklist exclude googleusercontent.com
	try set service dns forwarding blacklist exclude gstatic.com
	try set service dns forwarding blacklist exclude gvt1.com
	try set service dns forwarding blacklist exclude gvt1.net
	try set service dns forwarding blacklist exclude hb.disney.go.com
	try set service dns forwarding blacklist exclude herokuapp.com
	try set service dns forwarding blacklist exclude hp.com
	try set service dns forwarding blacklist exclude hulu.com
	try set service dns forwarding blacklist exclude images-amazon.com
	try set service dns forwarding blacklist exclude live.com
	try set service dns forwarding blacklist exclude microsoft.com
	try set service dns forwarding blacklist exclude microsoftonline.com
	try set service dns forwarding blacklist exclude msdn.com
	try set service dns forwarding blacklist exclude msecnd.net
	try set service dns forwarding blacklist exclude msftncsi.com
	try set service dns forwarding blacklist exclude mywot.com
	try set service dns forwarding blacklist exclude nsatc.net
	try set service dns forwarding blacklist exclude paypal.com
	try set service dns forwarding blacklist exclude pop.h-cdn.co
	try set service dns forwarding blacklist exclude rackcdn.com
	try set service dns forwarding blacklist exclude rarlab.com
	try set service dns forwarding blacklist exclude schema.org
	try set service dns forwarding blacklist exclude shopify.com
	try set service dns forwarding blacklist exclude skype.com
	try set service dns forwarding blacklist exclude smacargo.com
	try set service dns forwarding blacklist exclude sourceforge.net
	try set service dns forwarding blacklist exclude spotify.com
	try set service dns forwarding blacklist exclude spotify.edgekey.net
	try set service dns forwarding blacklist exclude spotilocal.com
	try set service dns forwarding blacklist exclude ssl-on9.com
	try set service dns forwarding blacklist exclude ssl-on9.net
	try set service dns forwarding blacklist exclude sstatic.net
	try set service dns forwarding blacklist exclude static.chartbeat.com
	try set service dns forwarding blacklist exclude storage.googleapis.com
	try set service dns forwarding blacklist exclude twimg.com
	try set service dns forwarding blacklist exclude viewpoint.com
	try set service dns forwarding blacklist exclude windows.net
	try set service dns forwarding blacklist exclude xboxlive.com
	try set service dns forwarding blacklist exclude yimg.com
	try set service dns forwarding blacklist exclude ytimg.com
	try set service dns forwarding blacklist hosts include beap.gemini.yahoo.com
	try set service dns forwarding blacklist hosts include ads.feedly.com
	try set service dns forwarding blacklist hosts source hostsfile.org description '"hostsfile.org bad hosts blacklist"'
	try set service dns forwarding blacklist hosts source hostsfile.org prefix '127.0.0.1'
	try set service dns forwarding blacklist hosts source hostsfile.org url 'http://www.hostsfile.org/Downloads/hosts.txt'
	try set service dns forwarding blacklist hosts source openphish description '"OpenPhish automatic phishing detection"'
	try set service dns forwarding blacklist hosts source openphish prefix 'http'
	try set service dns forwarding blacklist hosts source openphish url 'https://openphish.com/feed.txt'
	# try set service dns forwarding blacklist hosts source raw.github.com description '"This hosts file is a merged collection of hosts from reputable sources"'
	# try set service dns forwarding blacklist hosts source raw.github.com prefix '0.0.0.0 '
	# try set service dns forwarding blacklist hosts source raw.github.com url 'https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts'
	try set service dns forwarding blacklist hosts source sysctl.org description '"This hosts file is a merged collection of hosts from Cameleon"'
	try set service dns forwarding blacklist hosts source sysctl.org prefix '127.0.0.1 '
	try set service dns forwarding blacklist hosts source sysctl.org url 'http://sysctl.org/cameleon/hosts'
	# try set service dns forwarding blacklist hosts source yoyo description '"Fully Qualified Domain Names only - no prefix to strip"'
	# try set service dns forwarding blacklist hosts source yoyo prefix '127.0.0.1 '
	# try set service dns forwarding blacklist hosts source yoyo url "http://pgl.yoyo.org/as/serverlist.php"
	try set system task-scheduler task update_blacklists executable path /config/scripts/update-dnsmasq
	try set system task-scheduler task update_blacklists interval 1d
	try commit
	try save
	try end
}

update_dns_config
set_vyattacfg_grp
