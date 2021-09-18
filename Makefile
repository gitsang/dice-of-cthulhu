.PHONY: cthulhu

SHELL=/bin/bash

default: cthulhu

cthulhu:

	mkdir -p cthulhu/{bin,conf,log}
	go build -o cthulhu/bin/cthulhu cmd/cthulhu.go
	cp configs/*.yml cthulhu/conf/

tar: cthulhu

	tar zcvf cthulhu/cthulhu.tar.gz cthulhu

install: cthulhu

	systemctl stop cthulhu
	mkdir -p /usr/local/etc/cthulhu
	mkdir -p /var/log/cthulhu
	cp cthulhu/bin/cthulhu /usr/local/bin/cthulhu
	cp configs/cthulhu.yml /usr/local/etc/cthulhu/cthulhu.yml
	cp configs/cthulhu.service /usr/lib/systemd/system/cthulhu.service

run: install

	systemctl start cthulhu
	systemctl status cthulhu

clean:

	rm -fr bin cthulhu
