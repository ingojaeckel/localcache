#!/bin/sh

FILE=go_appengine_sdk_linux_386-1.9.15.zip

if [ ! -e ${FILE} ]
then
	wget https://storage.googleapis.com/appengine-sdks/featured/${FILE}
	unzip ${FILE}
else
	echo "Skip downloading Go GAE SDK since it already exists"
fi

export GOROOT=go_appengine/goroot
export PATH=go_appengine:$GOROOT/bin:$PATH

./tests.sh
./coverage.sh
