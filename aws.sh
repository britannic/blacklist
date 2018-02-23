#!/bin/bash
# Updates the AWS repository and pushes it to github

VERSION="1.0"
host=${1}
pkg=${2}
tag=${3}

ssh -tt ${host} <<-EOF
	reprepro includedeb wheezy /tmp/${pkg}_*.deb
	git add --all
	git commit -am"Package repository release $(TAG)"
	git tag $(TAG)
	git push origin master
	git push --tags
	exit
EOF

