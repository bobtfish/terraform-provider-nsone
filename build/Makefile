.PHONY: itest_% clean shell
PROJECT = terraform-provider-nsone

BUILD_NUMBER?=0
ifdef upstream_build_number
	REAL_BUILD_NUMBER=$(upstream_build_number)
else
	REAL_BUILD_NUMBER?=$(BUILD_NUMBER)
endif
VERSION = 0.5
TF_VERSION = 0.8.2
ITERATION = default$(REAL_BUILD_NUMBER)

DOCKER_TAG = $(PROJECT)_$(shell date +%s)

CLEAN_CONTAINER := [ -e .docker_container_id ] && docker rm --force $$(cat .docker_container_id) || true; rm -f .docker_container_id

output: .docker_container_id
	docker cp $$(cat .docker_container_id):/dist ..

.docker_container_id: .docker_image_id
	docker run --rm=false \
		-v $(CURDIR)/..:/go/src/github.com/bobtfish/terraform-provider-nsone:ro \
		-v $(CURDIR)/build.sh:/build.sh:ro \
		--cidfile=$(CURDIR)/.docker_container_id \
		$$(cat .docker_image_id) \
		bash -xve /build.sh $(PROJECT) $(VERSION) $(ITERATION) $(TF_VERSION) || \
	(retval=$$?; $(CLEAN_CONTAINER); exit $$retval; )

.docker_image_id: Dockerfile
	docker build -t $(DOCKER_TAG) .
	docker inspect -f '{{ .Id }}' $(DOCKER_TAG) > .docker_image_id

clean:
	$(CLEAN_CONTAINER)
	rm -f .docker_image_id
	rm -rf ../dist

shell: .docker_image_id
	docker run --rm -t -i $$(cat .docker_image_id) /bin/bash

