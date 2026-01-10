define TAIL__DOCKER_ENSURE_VOLUME
	@if ! docker volume inspect $(1) >/dev/null 2>&1; then \
		docker volume create $(1); \
	fi
endef

define TAIL__DOCKER_BUILD
	$(call TAIL__DOCKER_ENSURE_VOLUME, tail-pkg-mod)
	@docker run --rm \
		--name tail-build \
        -v "${TAIL__APP_PATH}":/app \
        -v "${TAIL__ROOT_PATH}":/src \
        -v tail-pkg-mod:/go/pkg/mod \
        -w /src \
        golang:1.25.3-alpine3.22 \
        sh -c " \
            echo '$(call TAIL__MESSAGE_DOCKER_BUILD_MOD_DOWNLOAD)' && \
            go mod download 2>/dev/null || go mod init temp-module && \
            echo '$(call TAIL__MESSAGE_DOCKER_BUILD_GO_COMPILING)' && \
            go build -ldflags='-w -s' -o $(1) $(2) | sed 's|/src/||g' && \
            echo '$(call TAIL__MESSAGE_DOCKER_BUILD_SUCCESS)' \
        "
    @chmod +x "${TAIL__APP_PATH}/tail"
endef

define TAIL__DOCKER_LINT
	$(call TAIL__DOCKER_ENSURE_VOLUME, tail-pkg-mod)
	$(call TAIL__DOCKER_ENSURE_VOLUME, tail-bin)
	@docker run --rm \
		--name tail-lint \
        -v "${TAIL__ROOT_PATH}":/src \
        -v tail-bin:/go/bin \
        -v tail-pkg-mod:/go/pkg/mod \
        -w /src \
		golangci/golangci-lint:v2.6.2 \
		sh -c " \
			command -v gosec >/dev/null 2>&1 || go install github.com/securego/gosec/v2/cmd/gosec@latest; \
			golangci-lint migrate >/dev/null 2>&1; \
			echo '$(call TAIL__MESSAGE_DOCKER_LINT_FORMAT)' && \
			gofmt -l -d . | sed 's|/src/| |g' && \
			echo '' && \
			echo '$(call TAIL__MESSAGE_DOCKER_LINT_STATIC_ANALYZER)' && \
			golangci-lint cache clean && \
			golangci-lint --config=.golangci.yml run ./... | sed 's|/src/| |g' && \
			echo '' && \
			echo '$(call TAIL__MESSAGE_DOCKER_LINT_SECURITY_CHECK)' && \
			gosec -quiet ./... | sed 's|/src/| |g' | sed 's|] - | ] - |g' && \
			echo '' && \
			echo '$(call TAIL__MESSAGE_DOCKER_LINT_MOD_VERIFY)' && \
			go mod verify | sed 's|/src/| |g' && \
			echo '' && \
			echo '$(call TAIL__MESSAGE_DOCKER_LINT_SUCCESS)' \
        "
endef

define TAIL__DOCKER_TEST
	$(call TAIL_ENSURE_VOLUME, tail-pkg-mod)
	$(call TAIL_ENSURE_VOLUME, tail-bin)
	@docker run --rm \
		--name tail-test \
        -v "${TAIL__ROOT_PATH}":/src \
        -v tail-bin:/go/bin \
        -v tail-pkg-mod:/go/pkg/mod \
        -w /src \
        golang:1.25.3-alpine3.22 \
        sh -c " \
        	ginkgo version >/dev/null 2>&1 || go install github.com/onsi/ginkgo/v2/ginkgo@v2.27.3; \
            echo '$(call TAIL__MESSAGE_DOCKER_TEST_RUN)' && \
            ginkgo -r ./... | sed 's|/src/||g' && \
            echo '$(call TAIL__MESSAGE_DOCKER_TEST_SUCCESS)' \
        "
endef

define TAIL__DOCKER_RUN_INTEGRATION_TESTS
	$(call TAIL__DOCKER_ENSURE_VOLUME, tail-pkg-mod)
	@docker run --rm \
		--name tail-build \
        -v "${TAIL__APP_PATH}":/app \
        -v "${TAIL__ROOT_PATH}":/src \
        -v tail-pkg-mod:/go/pkg/mod \
        -w /src \
        golang:1.25.3-alpine3.22 \
        sh -c " \
            echo '$(call TAIL__MESSAGE_DOCKER_BUILD_MOD_DOWNLOAD)' && \
            go mod download 2>/dev/null || go mod init temp-module && \
            echo '$(call TAIL__MESSAGE_DOCKER_INTEGRATION_TEST_RUN)' && \
            go build -ldflags='-w -s' -o /app/tail $(1) | sed 's|/src/||g' && \
            ./tools/tests/integration_tests.sh ./tests/integration /app | sed 's|/src/||g' && \
            echo '$(call TAIL__MESSAGE_DOCKER_INTEGRATION_TEST_SUCCESS)'\
        "
endef

define TAIL__DOCKER_BUILD_INTEGRATION_EXPECTED
	$(call TAIL__DOCKER_ENSURE_VOLUME, tail-pkg-mod)
	@docker run --rm \
		--name tail-build \
        -v "${TAIL__APP_PATH}":/app \
        -v "${TAIL__ROOT_PATH}":/src \
        -v tail-pkg-mod:/go/pkg/mod \
        -w /src \
        golang:1.25.3-alpine3.22 \
        sh -c " \
            echo '$(call TAIL__MESSAGE_DOCKER_BUILD_MOD_DOWNLOAD)' && \
            go mod download 2>/dev/null || go mod init temp-module && \
            echo '$(call TAIL__MESSAGE_DOCKER_BUILD_INTEGRATION_EXPECTED)' && \
            go build -ldflags='-w -s' -o /app/tail $(1) | sed 's|/src/||g' && \
            ./tools/tests/make_integration_expected.sh ./tests/integration /app | sed 's|/src/||g' && \
            echo '$(call TAIL__MESSAGE_DOCKER_BUILD_INTEGRATION_EXPECTED_SUCCESS)'\
        "
endef