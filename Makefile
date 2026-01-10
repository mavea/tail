include $(PWD)/tools/make/main.mk

.PHONY: default
default: lint test build test-module

.PHONY: build
build:
	$(call TAIL__DOCKER_BUILD,/app/tail,./cmd/)

.PHONY: lint
lint:
	$(call TAIL__DOCKER_LINT)

.PHONY: test
test:
	$(call TAIL__DOCKER_TEST)


.PHONY: integration-test
integration-test:
	$(call TAIL__RUN_INTEGRATION_TESTS)

.PHONY: build-integration-expected
build-integration-expected:
	$(call TAIL__BUILD_INTEGRATION_EXPECTED)
