# We need access to Bash features, such as arrays
SHELL := /bin/bash

.PHONY: itest_% clean shell
PROJECT = nsone


BUILD_NUMBER?=0
ifdef upstream_build_number
	REAL_BUILD_NUMBER=$(upstream_build_number)
else
	REAL_BUILD_NUMBER?=$(BUILD_NUMBER)
endif
VERSION = 0.9.1
TF_VERSION = 0.12 # TODO: nix
TF_VERSIONS := 0.12 0.13
ITERATION = yelp$(REAL_BUILD_NUMBER)
ARCH := $(shell dpkg --print-architecture)

PACKAGE_NAMES := $(foreach TF_VERSION,$(TF_VERSIONS),terraform-provider-$(PROJECT)-$(TF_VERSION)_$(VERSION)-$(ITERATION)_amd64.deb)
PACKAGE_FILES := $(foreach PACKAGE_NAME,$(PACKAGE_NAMES),dist/$(PACKAGE_NAME))

DOCKER_TAG = terraform-provider-$(PROJECT)_$(shell date +%s)

CLEAN_CONTAINER := [ -e .docker_container_id ] && docker rm --force $$(cat .docker_container_id) || true; rm -f .docker_container_id

itest_%: $(PACKAGE_FILES)
	export TF_VERSIONS=( $(TF_VERSIONS) ) && \
	export PACKAGE_FILES=( $(PACKAGE_FILES) ) && \
	export VERSION=( $(VERSION) ) && \
	for ((i=0; i<$${#TF_VERSIONS[@]}; ++i)); do \
		docker run --rm -v $(CURDIR)/../dist:/dist:ro -v $(CURDIR)/itest.sh:/itest.sh:ro docker-dev.yelpcorp.com/$*_yelp:latest bash /itest.sh /$${PACKAGE_FILES[i]} $${TF_VERSIONS[i]} $${VERSION}; \
	done

$(PACKAGE_FILES): .docker_container_id
	mkdir -p ../dist && \
	for PACKAGE_FILE in $(PACKAGE_FILES); do \
		docker cp $$(cat .docker_container_id):/$$PACKAGE_FILE ../dist/; \
	done

.docker_container_id: .docker_image_id
	docker run --rm=false \
		-v $(CURDIR)/..:/go/src/terraform-provider-nsone:ro \
		-v $(CURDIR)/build.sh:/build.sh:ro \
		--cidfile=$(CURDIR)/.docker_container_id \
		$$(cat .docker_image_id) \
		bash /build.sh $(PROJECT) $(VERSION) $(ITERATION) $(TF_VERSIONS) || \
	(retval=$$?; $(CLEAN_CONTAINER); exit $$retval; )

.docker_image_id: Dockerfile go.mod go.sum
	docker build -t $(DOCKER_TAG) .
	docker inspect -f '{{ .Id }}' $(DOCKER_TAG) > .docker_image_id

clean:
	$(CLEAN_CONTAINER)
	rm -f .docker_image_id
	rm -rf ../dist

shell: .docker_image_id
	docker run --rm -t -i $$(cat .docker_image_id) /bin/bash
