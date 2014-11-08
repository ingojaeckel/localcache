#!/bin/sh
goapp test -cover -test.v=true -test.coverprofile=c.out
sed -i -e "s#.*/\(.*\.go\)#\./\\1#" c.out
goapp tool cover -html c.out -o coverage.html
