include $(PWD)/tools/make/main.mk

.PHONY: default
default: generate-check lint test integration-test build

.PHONY: build
build:
	$(call TAIL__DOCKER_BUILD,/app/tail,./cmd/)

.PHONY: debug
debug:
	$(call TAIL__DOCKER_DEBUG,$(ARGS))

.PHONY: lint
lint:
	$(call TAIL__DOCKER_LINT)

.PHONY: test
test:
	$(call TAIL__DOCKER_TEST)

.PHONY: generate
generate:
	$(call TAIL__DOCKER_GENERATE)

.PHONY: generate-check
generate-check:
	$(call TAIL__DOCKER_GENERATE_CHECK)

.PHONY: integration-test
integration-test:
	$(call TAIL__RUN_INTEGRATION_TESTS)

.PHONY: build-integration-expected
build-integration-expected:
	$(call TAIL__BUILD_INTEGRATION_EXPECTED)

# Interactive test runner: runs a test program, waits for Enter, then pipes through tail
# Usage: make run-test TEST=<number>
# Example: make run-test TEST=31
# This will:
#   1. Find and compile the test program for test #31 (e.g., 31-mixed-csi-sgr)
#   2. Run the program and display its output
#   3. Wait for user to press Enter
#   4. Pipe the program output through tail with 15 lines, no log clearing
.PHONY: run-test
run-test:
	@if [ -z "$(TEST)" ]; then \
		echo "Error: TEST parameter is required"; \
		echo "Usage: make run-test TEST=<number>"; \
		echo "Example: make run-test TEST=31"; \
		exit 1; \
	fi
	$(call TAIL__DOCKER_RUN_INTERACTIVE_TEST,$(TEST))
