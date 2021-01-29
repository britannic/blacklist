#!/usr/bin/env bash

# Set up the Vyatta environment
declare -i DEC
API=/bin/cli-shell-api
CFGRUN=/opt/vyatta/sbin/vyatta-cfg-cmd-wrapper
shopt -s expand_aliases

alias begin='${CFGRUN} begin'
alias cleanup='${CFGRUN} cleanup'
alias commit='${CFGRUN} commit'
alias delete='${CFGRUN} delete'
alias end='${CFGRUN} end'
alias save='sudo ${CFGRUN} save'
alias set='${CFGRUN} set'

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

# Set the group so that the admin user will be able to commit configs
set_vyattacfg_grp() {
if [[ 'vyattacfg' != $(id -ng) ]]; then
  exec sg vyattacfg -c "$0 $@"
fi
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

noblacklist() {
	${API} existsActive service dns forwarding blacklist && return 1
	return 0
}

# Load the [service dns forwarding blacklist] configuration
update_dns_config() {
	try begin
	# try set service dns forwarding blacklist disabled false
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
	try set service dns forwarding blacklist domains source malwaredomains.com description '"Just Domains"'
	try set service dns forwarding blacklist domains source malwaredomains.com url 'http://mirror1.malwaredomains.com/files/justdomains'
	try set service dns forwarding blacklist domains source NoBitCoin description '"Blocking Web Browser Bitcoin Mining"'
	try set service dns forwarding blacklist domains source NoBitCoin prefix '0.0.0.0'
	try set service dns forwarding blacklist domains source NoBitCoin url 'https://raw.githubusercontent.com/hoshsadiq/adblock-nocoin-list/master/hosts.txt'
	try set service dns forwarding blacklist domains source OISD description '"OISD Domains Basic"'
	try set service dns forwarding blacklist domains source OISD url 'https://dbl.oisd.nl/basic/'
	try set service dns forwarding blacklist domains source simple_tracking description '"Basic tracking list by Disconnect"'
	try set service dns forwarding blacklist domains source simple_tracking url 'https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt'
	try set service dns forwarding blacklist exclude 1e100.net
	try set service dns forwarding blacklist exclude 2o7.net
	try set service dns forwarding blacklist exclude adjust.com
	try set service dns forwarding blacklist exclude adobedtm.com
	try set service dns forwarding blacklist exclude akamai.net
	try set service dns forwarding blacklist exclude akamaihd.net
	try set service dns forwarding blacklist exclude amazon.com
	try set service dns forwarding blacklist exclude amazonaws.com
	try set service dns forwarding blacklist exclude ampproject.org
	try set service dns forwarding blacklist exclude android.clients.google.com
	try set service dns forwarding blacklist exclude apple.com
	try set service dns forwarding blacklist exclude apresolve.spotify.com
	try set service dns forwarding blacklist exclude ask.com
	try set service dns forwarding blacklist exclude avast.com
	try set service dns forwarding blacklist exclude avira-update.com
	try set service dns forwarding blacklist exclude bannerbank.com
	try set service dns forwarding blacklist exclude bazaarvoice.com
	try set service dns forwarding blacklist exclude bing.com
	try set service dns forwarding blacklist exclude bit.ly
	try set service dns forwarding blacklist exclude bitdefender.com
	try set service dns forwarding blacklist exclude bonsaimirai.us9.list-manage.com
	try set service dns forwarding blacklist exclude c.s-microsoft.com
	try set service dns forwarding blacklist exclude cdn.ravenjs.com
	try set service dns forwarding blacklist exclude cdn.visiblemeasures.com
	try set service dns forwarding blacklist exclude clientconfig.passport.net
	try set service dns forwarding blacklist exclude clients2.google.com
	try set service dns forwarding blacklist exclude clients4.google.com
	try set service dns forwarding blacklist exclude cloudfront.net
	try set service dns forwarding blacklist exclude coremetrics.com
	try set service dns forwarding blacklist exclude dickssportinggoods.com
	try set service dns forwarding blacklist exclude dl.dropboxusercontent.com
	try set service dns forwarding blacklist exclude dropbox.com
	try set service dns forwarding blacklist exclude ebay.com
	try set service dns forwarding blacklist exclude edgesuite.net
	try set service dns forwarding blacklist exclude evernote.com
	try set service dns forwarding blacklist exclude express.co.uk
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
	try set service dns forwarding blacklist exclude i.s-microsoft.com
	try set service dns forwarding blacklist exclude images-amazon.com
	try set service dns forwarding blacklist exclude live.com
	try set service dns forwarding blacklist exclude logmein.com
	try set service dns forwarding blacklist exclude m.weeklyad.target.com
	try set service dns forwarding blacklist exclude magnetmail1.net
	try set service dns forwarding blacklist exclude microsoft.com
	try set service dns forwarding blacklist exclude microsoftonline.com
	try set service dns forwarding blacklist exclude msdn.com
	try set service dns forwarding blacklist exclude msecnd.net
	try set service dns forwarding blacklist exclude msftncsi.com
	try set service dns forwarding blacklist exclude mywot.com
	try set service dns forwarding blacklist exclude nsatc.net
	try set service dns forwarding blacklist exclude outlook.office365.com
	try set service dns forwarding blacklist exclude paypal.com
	try set service dns forwarding blacklist exclude pop.h-cdn.co
	try set service dns forwarding blacklist exclude products.office.com
	try set service dns forwarding blacklist exclude quora.com
	try set service dns forwarding blacklist exclude rackcdn.com
	try set service dns forwarding blacklist exclude rarlab.com
	try set service dns forwarding blacklist exclude s.youtube.com
	try set service dns forwarding blacklist exclude schema.org
	try set service dns forwarding blacklist exclude shopify.com
	try set service dns forwarding blacklist exclude skype.com
	try set service dns forwarding blacklist exclude smacargo.com
	try set service dns forwarding blacklist exclude sourceforge.net
	try set service dns forwarding blacklist exclude spclient.wg.spotify.com
	try set service dns forwarding blacklist exclude spotify.com
	try set service dns forwarding blacklist exclude spotify.edgekey.net
	try set service dns forwarding blacklist exclude spotilocal.com
	try set service dns forwarding blacklist exclude ssl-on9.com
	try set service dns forwarding blacklist exclude ssl-on9.net
	try set service dns forwarding blacklist exclude sstatic.net
	try set service dns forwarding blacklist exclude static.chartbeat.com
	try set service dns forwarding blacklist exclude storage.googleapis.com
	try set service dns forwarding blacklist exclude twimg.com
	try set service dns forwarding blacklist exclude video-stats.l.google.com
	try set service dns forwarding blacklist exclude viewpoint.com
	try set service dns forwarding blacklist exclude weeklyad.target.com
	try set service dns forwarding blacklist exclude weeklyad.target.com.edgesuite.net
	try set service dns forwarding blacklist exclude windows.net
	try set service dns forwarding blacklist exclude www.msftncsi.com
	try set service dns forwarding blacklist exclude xboxlive.com
	try set service dns forwarding blacklist exclude yimg.com
	try set service dns forwarding blacklist exclude ytimg.com
	try set service dns forwarding blacklist hosts exclude cfvod.kaltura.com
	try set service dns forwarding blacklist hosts include ads.feedly.com
	try set service dns forwarding blacklist hosts include beap.gemini.yahoo.com
	try set service dns forwarding blacklist hosts source githubSteveBlack description '"Blacklists adware and malware websites"'
	try set service dns forwarding blacklist hosts source githubSteveBlack prefix '0.0.0.0 '
	try set service dns forwarding blacklist hosts source githubSteveBlack url 'https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts'
	try set service dns forwarding blacklist hosts source openphish description '"OpenPhish automatic phishing detection"'
	try set service dns forwarding blacklist hosts source openphish prefix 'http'
	try set service dns forwarding blacklist hosts source openphish url 'https://openphish.com/feed.txt'
	try set system task-scheduler task update_blacklists executable arguments 10800
	try set system task-scheduler task update_blacklists executable path /config/scripts/update-dnsmasq-cronjob.sh
	try set system task-scheduler task update_blacklists interval 1d
	try commit
	try save
	try end
}

# echo "$@"
# Set group to vyattacfg
set_vyattacfg_grp

# Set UPGRADE flag
UPGRADE=0
[[ "${1}" == "configure" ]] && [[ -z "${2}" ]] && UPGRADE=1

noblacklist && UPGRADE=1

# Only run the post installation script if this is a first time installation
if [[ ${UPGRADE} == 1 ]] ; then
	echo "Installing blacklist configuration settings..."
	update_dns_config
fi