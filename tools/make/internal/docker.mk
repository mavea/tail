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
            go mod download && \
            echo '$(call TAIL__MESSAGE_DOCKER_BUILD_GO_COMPILING)' && \
            sh ./tools/make/sh/build.sh $(1) $(2) && \
            echo '$(call TAIL__MESSAGE_DOCKER_BUILD_SUCCESS)' || \
            { echo '$(call TAIL__MESSAGE_DOCKER_BUILD_FAIL)'; exit 1; } \
        "
    @cp ${TAIL__APP_PATH}/tail ${TAIL__ROOT_PATH}/tail
endef

define TAIL__DOCKER_DEBUG
	$(call TAIL__DOCKER_ENSURE_VOLUME, tail-pkg-mod)
	$(call TAIL__DOCKER_ENSURE_VOLUME, tail-bin)
	@docker run --rm -it \
		--name tail-debug \
        -p 2345:2345 \
        -v "${TAIL__APP_PATH}":/app \
        -v "${TAIL__ROOT_PATH}":/src \
        -v tail-bin:/go/bin \
        -v tail-pkg-mod:/go/pkg/mod \
        -w /src \
        golang:1.25.3-alpine3.22 \
        sh -c ' \
            set -eu; \
            echo "$(call TAIL__MESSAGE_DOCKER_DEBUG_PREPARE)"; \
            go mod download; \
            echo "$(call TAIL__MESSAGE_DOCKER_BUILD_MOD_DOWNLOAD)"; \
            if [ ! -x /go/bin/dlv ]; then \
                GOBIN=/go/bin go install github.com/go-delve/delve/cmd/dlv@latest; \
            fi; \
            echo "$(call TAIL__MESSAGE_DOCKER_DEBUG_BUILD)"; \
			if ! sh ./tools/make/sh/debug_build.sh; then \
        		echo "$(call TAIL__MESSAGE_DOCKER_DEBUG_FAIL)"; \
				exit 1; \
			fi; \
            echo "$(call TAIL__MESSAGE_DOCKER_DEBUG_WAIT)"; \
            exec /go/bin/dlv exec /app/tail-debug --headless --listen=:2345 --api-version=2 --accept-multiclient --only-same-user=false -- "$$@" \
        ' sh $(1)
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
              set -o pipefail; \
              failed=0; \
              if ! command -v gosec >/dev/null 2>&1; then \
                if ! go install github.com/securego/gosec/v2/cmd/gosec@latest; then \
                  failed=1; \
                fi; \
              fi; \
                      golangci-lint migrate >/dev/null 2>&1 || true; \
              echo '$(call TAIL__MESSAGE_DOCKER_LINT_FORMAT)'; \
				  if ! gofmt -l -d . | sed 's|/src/| |g' > /tmp/gofmt.out; then \
					failed=1; \
				  fi; \
				  if [ -s /tmp/gofmt.out ]; then \
					cat /tmp/gofmt.out; \
				failed=1; \
              fi; \
              echo ''; \
              echo '$(call TAIL__MESSAGE_DOCKER_LINT_STATIC_ANALYZER)'; \
              if ! golangci-lint cache clean; then \
                failed=1; \
              fi; \
              if ! golangci-lint --config=.golangci.yml run ./... | sed 's|/src/| |g'; then \
                failed=1; \
              fi; \
              echo ''; \
              echo '$(call TAIL__MESSAGE_DOCKER_LINT_SECURITY_CHECK)'; \
              if ! gosec -quiet ./... | sed 's|/src/| |g' | sed 's|] - | ] - |g'; then \
                failed=1; \
              fi; \
              echo ''; \
              echo '$(call TAIL__MESSAGE_DOCKER_LINT_MOD_VERIFY)'; \
              if ! go mod verify | sed 's|/src/| |g'; then \
                failed=1; \
              fi; \
              echo ''; \
			  if [ "\$$failed" != "0" ]; then \
				echo '$(call TAIL__MESSAGE_DOCKER_LINT_FAIL)'; \
				exit 1; \
			  fi; \
              echo '$(call TAIL__MESSAGE_DOCKER_LINT_SUCCESS)'; \
        "
endef

define TAIL__DOCKER_TEST
	$(call TAIL__DOCKER_ENSURE_VOLUME, tail-pkg-mod)
	$(call TAIL__DOCKER_ENSURE_VOLUME, tail-bin)
	@docker run --rm \
		--name tail-test \
        -v "${TAIL__ROOT_PATH}":/src \
        -v tail-bin:/go/bin \
        -v tail-pkg-mod:/go/pkg/mod \
        -w /src \
        golang:1.25.3-alpine3.22 \
        sh -c " \
            set -o pipefail; \
        	ginkgo version >/dev/null 2>&1 || go install github.com/onsi/ginkgo/v2/ginkgo@v2.27.3; \
            echo '$(call TAIL__MESSAGE_DOCKER_TEST_RUN)' && \
            ./tools/make/sh/tests/unit.sh && \
            echo '$(call TAIL__MESSAGE_DOCKER_TEST_SUCCESS)' || \
            { echo '$(call TAIL__MESSAGE_DOCKER_TEST_FAIL)'; exit 1; } \
        "
endef

define TAIL__DOCKER_GENERATE
	$(call TAIL__DOCKER_ENSURE_VOLUME, tail-pkg-mod)
	@docker run --rm \
		--name tail-generate \
		-v "${TAIL__ROOT_PATH}":/src \
		-v tail-pkg-mod:/go/pkg/mod \
		-w /src \
		golang:1.25.3-alpine3.22 \
		sh -c " \
      set -e; \
      echo '$(call TAIL__MESSAGE_DOCKER_GENERATE_RUN)' && \
			go mod download && \
      go generate ./tests/mocks | sed 's|/src/||g' && \
      echo '$(call TAIL__MESSAGE_DOCKER_GENERATE_SUCCESS)' || \
      { echo '$(call TAIL__MESSAGE_DOCKER_GENERATE_FAIL)'; exit 1; } \
		"
endef

define TAIL__DOCKER_GENERATE_CHECK
  $(call TAIL__DOCKER_ENSURE_VOLUME, tail-pkg-mod)
  @docker run --rm \
    --name tail-generate-check \
    -v "${TAIL__ROOT_PATH}":/src \
    -v tail-pkg-mod:/go/pkg/mod \
    -w /src \
    golang:1.25.3-alpine3.22 \
    sh -c " \
      set -e; \
      echo '$(call TAIL__MESSAGE_DOCKER_GENERATE_CHECK_RUN)' && \
      go mod download && \
      tmp_dir=/tmp/tail-generate-check-$$ && \
      rm -rf $$tmp_dir && \
      mkdir -p $$tmp_dir/repo && \
      cp -R /src/. $$tmp_dir/repo && \
      (cd $$tmp_dir/repo && go generate ./tests/mocks >/dev/null) && \
      if diff -ruN ./tests/mocks $$tmp_dir/repo/tests/mocks > $$tmp_dir/generate.diff; then \
        echo '$(call TAIL__MESSAGE_DOCKER_GENERATE_CHECK_SUCCESS)'; \
      else \
        cat $$tmp_dir/generate.diff; \
        echo '$(call TAIL__MESSAGE_DOCKER_GENERATE_CHECK_FAIL)'; \
        exit 1; \
      fi \
    "
endef

define TAIL__DOCKER_RUN_INTEGRATION_TESTS
	$(call TAIL__DOCKER_ENSURE_VOLUME, tail-pkg-mod)
	@docker run --rm \
		--name tail-integration_test \
        -v "${TAIL__APP_PATH}":/app \
        -v "${TAIL__ROOT_PATH}":/home/runner/work/tail/tail/ \
        -v tail-pkg-mod:/go/pkg/mod \
        -w /home/runner/work/tail/tail/ \
        golang:1.25.3-alpine3.22 \
        sh -c " \
            set -o pipefail; \
            echo '$(call TAIL__MESSAGE_DOCKER_BUILD_MOD_DOWNLOAD)' && \
            go mod download && \
            echo '$(call TAIL__MESSAGE_DOCKER_INTEGRATION_TEST_RUN)' && \
            go build -ldflags='-w -s -X tail/internal/bootstrap.Version=integration.tests -X tail/internal/bootstrap.BuildTime=0000.00.00' -o /app/tail $(1) | sed 's|/home/runner/work/tail/tail/||g' && \
            sh ./tools/make/sh/tests/integration_tests.sh ./tests/integration /app --color=always && \
            echo '$(call TAIL__MESSAGE_DOCKER_INTEGRATION_TEST_SUCCESS)' || \
            { echo '$(call TAIL__MESSAGE_DOCKER_INTEGRATION_TEST_FAIL)'; exit 1; } \
        "
endef

define TAIL__DOCKER_BUILD_INTEGRATION_EXPECTED
	$(call TAIL__DOCKER_ENSURE_VOLUME, tail-pkg-mod)
	@docker run --rm \
		--name tail-build \
        -v "${TAIL__APP_PATH}":/app \
        -v "${TAIL__ROOT_PATH}":/home/runner/work/tail/tail \
        -v tail-pkg-mod:/go/pkg/mod \
        -w /home/runner/work/tail/tail \
        golang:1.25.3-alpine3.22 \
        sh -c " \
            echo '$(call TAIL__MESSAGE_DOCKER_BUILD_MOD_DOWNLOAD)' && \
            go mod download && \
            echo '$(call TAIL__MESSAGE_DOCKER_BUILD_INTEGRATION_EXPECTED)' && \
            go build -ldflags='-w -s -X tail/internal/bootstrap.Version=integration.tests -X tail/internal/bootstrap.BuildTime=0000.00.00' -o /app/tail $(1) | sed 's|/home/runner/work/tail/tail/||g' && \
            sh ./tools/make/sh/tests/make_integration_expected.sh ./tests/integration /app | sed 's|/home/runner/work/tail/tail/||g' && \
            echo '$(call TAIL__MESSAGE_DOCKER_BUILD_INTEGRATION_EXPECTED_SUCCESS)'\
        "
endef

define TAIL__DOCKER_RUN_INTERACTIVE_TEST
	$(call TAIL__DOCKER_ENSURE_VOLUME, tail-pkg-mod)
	@docker run --rm -it \
		--name tail-interactive-test \
		-v "${TAIL__APP_PATH}":/app \
		-v "${TAIL__ROOT_PATH}":/src \
		-v tail-pkg-mod:/go/pkg/mod \
		-w /src \
		golang:1.25.3-alpine3.22 \
		sh -c ' \
			set -o pipefail; \
			TEST_NUM="$(1)"; \
			ACTUAL_DIR=$$(ls -d ./tests/integration/tail/$$TEST_NUM-* 2>/dev/null | head -1); \
			if [ -z "$$ACTUAL_DIR" ]; then \
				echo "Test directory not found for test number: $(1)"; \
				exit 1; \
			fi; \
			TEST_NAME=$$(basename "$$ACTUAL_DIR"); \
			echo "📥 Загружаем зависимости..."; \
      go mod download; \
			echo "🔨 Собираем тестовую программу: $$TEST_NAME"; \
			go build -ldflags="-w -s" -o /app/test "$$ACTUAL_DIR/main.go"; \
			echo "▶️  Запуск тестовой программы (нажмите Enter для продолжения):"; \
			/app/test; \
			echo ""; \
			read -p "Нажмите Enter для запуска через tail с 20 строками..." dummy; \
			echo ""; \
			echo "▶️  Запуск через tail..."; \
			/app/tail -a test -n 20 -o direct -t none -r roller -f -c /app/test \
		'
endef

