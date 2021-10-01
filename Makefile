.PHONY: all clean test

PROJECT := nsone
OUTPUT ?= bin/terraform-provider-${PROJECT}

all: terraform-provider-${PROJECT} .git/hooks/pre-commit

install: terraform-provider-${PROJECT}
	cp -f $(GOPATH)/bin/terraform-provider-${PROJECT} $$(dirname $$(which terraform))

terraform-provider-${PROJECT}: main.go nsone/*.go
	mkdir -p $(GOPATH)/bin
	go build -o $(OUTPUT)

fmt:
	go fmt ./...

test: .git/hooks/pre-commit
	cd nsone ; go test -v .

clean:
	rm -f bin/terraform-provider-nsone
	make -C yelppack clean

.git/hooks/pre-commit:
	    if [ ! -f .git/hooks/pre-commit ]; then ln -s ../../git-hooks/pre-commit .git/hooks/pre-commit; fi

itest_%:
	cp -vp go.mod go.sum yelppack/
	mkdir -p dist
	make -C yelppack $@

package: itest_bionic

