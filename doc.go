package main

/*
UBNT edgeos-dnsmasq-blacklist dnsmasq DNS Blacklisting and Redirection

Follow the conversation @ https://community.ubnt.com/t5/EdgeMAX/DNS-Adblocking-amp-Blacklisting-dnsmasq-Configuration/td-p/2215008/jump-to/first-unread-message

Note: This is 3rd party software and isn't supported or endorsed by Ubiquiti Networks®
Overview

EdgeMax dnsmasq DNS blacklisting and redirection is inspired by the users at EdgeMAX Community

Copyright

Copyright © 2018 Helm Rock Consulting
Licenses

Redistribution and use in source and binary forms, with or without modification, are permitted
provided that the following conditions are met:

Redistributions of source code must retain the above copyright notice, this list of conditions
and the following disclaimer.

Redistributions in binary form must reproduce the above copyright notice, this list of conditions
and the following disclaimer in the documentation and/or other materials provided with the
distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR
IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND
FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR
CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN
ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

The views and conclusions contained in the software and documentation are those of the authors and
should not be interpreted as representing official policies, either expressed or implied, of the
FreeBSD Project.

Change Log

See changelog for details.

Features

° Adds DNS blacklisting integration to the EdgeRouter configuration

° Generates configuration files used directly by dnsmasq to redirect dns lookups

° Integrated with the EdgeMax OS CLI

° Any FQDN in the blacklist will force dnsmasq to return the configured dns redirect IP address
Compatibility

edgeos-dnsmasq-blacklist has been tested on the EdgeRouter ERLite-3, ERPoe-5, ER-X and UniFi
Security Gateway USG-3 routers versions EdgeMAX: v1.7.0-v1.9.7+hotfix.4, UniFi: v4.4.12-v4.4.18


EdgeRouter ERLite-3, ERPoe-5 & UniFi-Gateway-3

    curl https://community.ubnt.com/ubnt/attachments/ubnt/EdgeMAX/194030/13/edgeos-dnsmasq-blacklist_1.0.0.rc4_mips.deb.tgz | tar -xvz
    sudo dpkg -i edgeos-dnsmasq-blacklist_1.0.0.rc4_mips.deb

EdgeRouter ER-X & ER-X-SFP

    curl https://community.ubnt.com/ubnt/attachments/ubnt/EdgeMAX/194030/14/edgeos-dnsmasq-blacklist_1.0.0.rc4_mipsel.deb.tgz | tar -xvz
    sudo dpkg -i edgeos-dnsmasq-blacklist_1.0.0.rc4_mipsel.deb

Upgrade

Since dpkg cannot upgrade packages, follow the instructions under Installation and the previous package version will be automatically removed before the new package version is installed

Removal

• EdgeMAX - All Platforms

    sudo apt-get remove edgeos-dnsmasq-blacklist

Frequently Asked Questions

How do I disable/enable dnsmasq blacklisting?

Use these CLI configure commands:

° Disable:

    configure
    set service dns forwarding blacklist disabled true
    commit;save;exit

° Enable:
    configure
    set service dns forwarding blacklist disabled false
    commit;save;exit

Does the install backup my blacklist configuration before deleting it?

• If a blacklist configuration already exists, the install routine will automatically back it up to /config/user-data/blacklist.$(date +'%FT%H%M%S').cmds
How do I back up my blacklist configuration and restore it later?

• use the following commands (make a note of the file name):
    export DATE=$(date +'%FT%H%M%S'); echo "Backing up blacklist configuration to: /config/user-data/blacklist.${DATE}.cmds"; show configuration commands | grep blacklist > /config/user-data/blacklist.$(date +'%FT%H%M%S').cmds
    After installing the latest version, you can merge your backed up configuration:
    configure
    .  /config/user-data/blacklist.[date string].cmds
    commit;save;exit

° If you prefer to delete the default configuration and restore your previous configuration, run these commands:
    configure
    delete service dns forwarding blacklist
    .  /config/user-data/blacklist.[date string].cmds
    commit;save;exit

Which blacklist sources are installed by default?

° You can use this command in the CLI shell to view the current sources after installation or view the log and see previous downloads:
    show configuration commands | match blacklist | match source
    more /var/log/update-dnsmasq.log

How do I add or delete sources?

° To add a source, first check it can serve a text list and also note the prefix (if any) before the hosts or domains, e.g. http://www.malwaredomainlist.com/ has this format:
    #               MalwareDomainList.com Hosts List           #
    #   http://www.malwaredomainlist.com/hostslist/hosts.txt   #
    #         Last updated: Mon, 04 Dec 17 19:18:42 +0000      #


    127.0.0.1  localhost
    127.0.0.1  0koryu0.easter.ne.jp
    127.0.0.1  109-204-26-16.netconnexion.managedbroadband.co.uk
    127.0.0.1  1866809.securefastserver.com
    So the prefix is "127.0.0.1 "
    Here's how to creating the source in the CLI:
    configure
    set service dns forwarding blacklist hosts source malwaredomainlist description '127.0.0.1 based host and domain list'
    set service dns forwarding blacklist hosts source malwaredomainlist prefix '127.0.0.1  '
    set service dns forwarding blacklist hosts source malwaredomainlist url 'http://www.malwaredomainlist.com/hostslist/hosts.txt'
    commit;save;exit

° Using the CLI configure command, to delete domains and hosts sources:
    configure
    delete service dns forwarding blacklist domains source malc0de
    delete service dns forwarding blacklist hosts source yoyo.org
    commit;save;exit

How do I globally exclude or include hosts or a domains?

° Use these commands to globally include or exclude blacklisted entries:
    configure
    set service dns forwarding blacklist exclude cdn.visiblemeasures.com
    set service dns forwarding blacklist include www.nastywebsites.com
    commit;save;exit

    How do I exclude or include a host or a domain?

° Use these commands to include or exclude blacklisted entries:
    configure
    set service dns forwarding blacklist domains exclude visiblemeasures.com
    set service dns forwarding blacklist domains include domainsnastywebsites.com
    set service dns forwarding blacklist hosts exclude cdn.visiblemeasures.com
    set service dns forwarding blacklist hosts include www.nastywebsites.com
    commit;save;exit

    Does update-dnsmasq run automatically?

Yes, a scheduled task is created and run daily at midnight
° It can be reconfigured using these CLI configuration commands:
    set system task-scheduler task update_blacklists executable path /config/scripts/update-dnsmasq
    set system task-scheduler task update_blacklists interval 1d
    For example, to change the execution interval to every 6 hours, use this command:
    set system task-scheduler task update_blacklists interval 6h

° In daily use, no additional interaction with update-dnsmasq is required. By default, cron will run update-dnsmasq at midnight each day to download the blacklist sources and update the dnsmasq configuration files in /etc/dnsmasq.d. dnsmasq will automatically be reloaded after the configuration file update is completed.

How do I use the command line switches?

° update-dnsmasq has the following commandline switches available:
    /config/scripts/update-dnsmasq -h
        -dir string
                Override dnsmasq directory (default "/etc/dnsmasq.d")
        -f [full file path]
                [full file path] # Load a config.boot file
        -h   Display help
        -v   Verbose display
        -version
                Show version

EdgeOS dnsmasq Configuration

° dnsmasq may need to be configured to ensure blacklisting works correctly
° Here is an example using the EdgeOS configuration shell
    configure
    set service dns forwarding cache-size 2048
    set service dns forwarding except-interface [Your WAN i/f]
    set service dns forwarding name-server [Your choice of IPv4 Internet Name-Server]
    set service dns forwarding name-server [Your choice of IPv4 Internet Name-Server]
    set service dns forwarding name-server [Your choice of IPv6 Internet Name-Server]
    set service dns forwarding name-server [Your choice of IPv6 Internet Name-Server]
    set service dns forwarding options bogus-priv
    set service dns forwarding options domain-needed
    set service dns forwarding options domain=mydomain.local
    set service dns forwarding options enable-ra
    set service dns forwarding options expand-hosts
    set service dns forwarding options localise-queries
    set service dns forwarding options strict-order
    set service dns forwarding system
    set system name-server 127.0.0.1
    set system name-server '::1'
    commit; save; exit
*/
