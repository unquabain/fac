GOFILES=$(shell ls **/*.go | grep -v '.*_test\.go')
dist/thing-doer: $(GOFILES)
	go build -o dist/thing-doer

dist:
	mkdir dist