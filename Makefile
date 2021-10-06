GOFILES=$(shell ls **/*.go | grep -v '.*_test\.go')
dist/fac: $(GOFILES)
	go build -o dist/fac

dist:
	mkdir dist