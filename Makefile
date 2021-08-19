.PHONY: cthulhu

SHELL=/bin/bash

default: cthulhu

cthulhu:

	mkdir -p cthulhu/{bin,conf,log}
	go build -o cthulhu/bin/cthulhu cmd/cthulhu.go
	cp configs/* cthulhu/conf/

tar: cthulhu

	tar zcvf cthulhu/cthulhu.tar.gz cthulhu

clean:

	rm -fr bin cthulhu
