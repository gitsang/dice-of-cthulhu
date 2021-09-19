.PHONY: cthulhu

SHELL=/bin/bash

default: cthulhu

help:
	# make help      : help
	# make ctuhlhu   : build project
	# make tar       : build and tar
	# make install   : install to systemd
	# make run       : install and run by systemd
	# make clean     : clean make file
	# make uninstall : uninstall

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
	cp configs/* /usr/local/etc/cthulhu/
	cp configs/cthulhu.service /usr/lib/systemd/system/cthulhu.service

run: install

	systemctl start cthulhu
	tail -f /var/log/cthulhu/cthulhu.log

clean:

	rm -fr bin cthulhu

uninstall:

	systemctl stop cthulhu
	rm -fr /usr/local/bin/cthuluh /usr/local/etc/cthulhu /usr/lib/systemd/system/cthulhu.service
