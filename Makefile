.PHONY: cthulhu

SHELL=/bin/bash

default: cthulhu

cthulhu:

	mkdir -p bin
	go build -o bin/cthulhu cmd/cthulhu.go

tar: cthulhu

	mkdir -p cthulhu/{bin,conf,log}
	cp bin/* cthulhu/bin/
	cp configs/* cthulhu/conf/
	tar zcvf cthulhu/cthulhu.tar.gz cthulhu

clean:

	rm -fr bin cthulhu
