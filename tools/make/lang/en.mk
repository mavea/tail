define TAIL__MESSAGE_DOCKER_BUILD_MOD_DOWNLOAD
📥 Downloading dependencies...
endef

define TAIL__MESSAGE_DOCKER_BUILD_GO_COMPILING
🔨 Compiling...
endef

define TAIL__MESSAGE_DOCKER_BUILD_SUCCESS
✅ Compilation completed!
endef

define TAIL__MESSAGE_DOCKER_BUILD_FAIL
❌ Compilation failed!
endef

define TAIL__MESSAGE_DOCKER_DEBUG_PREPARE
🐞 Preparing debug environment and Delve...
endef

define TAIL__MESSAGE_DOCKER_DEBUG_BUILD
🔨 Building debug binary without optimizations...
endef

define TAIL__MESSAGE_DOCKER_DEBUG_FAIL
❌ Building debug failed!
endef

define TAIL__MESSAGE_DOCKER_DEBUG_WAIT
🧭 Delve is waiting for a connection on port 2345...
endef

define TAIL__MESSAGE_DOCKER_LINT_FORMAT
🔨 Checking formatting (gofmt)...
endef

define TAIL__MESSAGE_DOCKER_LINT_STATIC_ANALYZER
🔨 Static analysis (golangci-lint)...
endef

define TAIL__MESSAGE_DOCKER_LINT_SECURITY_CHECK
🔨 Security check (gosec)...
endef

define TAIL__MESSAGE_DOCKER_LINT_MOD_VERIFY
🔨 Verifying dependencies...
endef

define TAIL__MESSAGE_DOCKER_LINT_SUCCESS
✅ All checks completed successfully!
endef

define TAIL__MESSAGE_DOCKER_LINT_FAIL
❌ Checks failed!
endef

define TAIL__MESSAGE_DOCKER_TEST_RUN
🔨 Running unit tests...
endef

define TAIL__MESSAGE_DOCKER_TEST_SUCCESS
✅ Tests completed successfully!
endef

define TAIL__MESSAGE_DOCKER_TEST_FAIL
❌ Tests failed!
endef

define TAIL__MESSAGE_DOCKER_INTEGRATION_TEST_RUN
🔨 Building test scripts and running integration tests...
endef

define TAIL__MESSAGE_DOCKER_INTEGRATION_TEST_SUCCESS
✅ Tests completed successfully!
endef

define TAIL__MESSAGE_DOCKER_INTEGRATION_TEST_FAIL
❌ Tests failed!
endef

define TAIL__MESSAGE_DOCKER_BUILD_INTEGRATION_EXPECTED
🔨 Building test scripts and generating expected log files...
endef

define TAIL__MESSAGE_DOCKER_BUILD_INTEGRATION_EXPECTED_SUCCESS
✅ Generation completed successfully!
endef

define TAIL__MESSAGE_DOCKER_GENERATE_RUN
🔨 Generating mocks...
endef

define TAIL__MESSAGE_DOCKER_GENERATE_SUCCESS
✅ Generation completed successfully!
endef

define TAIL__MESSAGE_DOCKER_GENERATE_FAIL
❌ Generation failed!
endef

define TAIL__MESSAGE_DOCKER_GENERATE_CHECK_RUN
🔨 Checking generated mocks freshness...
endef

define TAIL__MESSAGE_DOCKER_GENERATE_CHECK_SUCCESS
✅ Generated mocks are up to date!
endef

define TAIL__MESSAGE_DOCKER_GENERATE_CHECK_FAIL
❌ Generated mocks are outdated! Run make generate.
endef

