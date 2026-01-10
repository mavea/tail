define TAIL__GET_INTEGRATION_TESTS
$(shell find $(1) -maxdepth 1 -type d ! -name "$(notdir $(1))" | sort)
endef

define TAIL__RUN_INTEGRATION_TESTS
$(call TAIL__DOCKER_RUN_INTEGRATION_TESTS,./cmd/)
endef

define TAIL__BUILD_INTEGRATION_EXPECTED
$(call TAIL__DOCKER_BUILD_INTEGRATION_EXPECTED,./cmd/)
endef