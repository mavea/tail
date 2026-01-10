define TAIL__MESSAGE_DOCKER_BUILD_MOD_DOWNLOAD
📥 Загружаем зависимости...
endef

define TAIL__MESSAGE_DOCKER_BUILD_GO_COMPILING
🔨 Компилируем...
endef

define TAIL__MESSAGE_DOCKER_BUILD_SUCCESS
✅ Компиляция завершена!
endef

define TAIL__MESSAGE_DOCKER_LINT_FORMAT
🔨 Проверка форматирования (gofmt)...
endef

define TAIL__MESSAGE_DOCKER_LINT_STATIC_ANALYZER
🔨 Статический анализ (golangci-lint)...
endef

define TAIL__MESSAGE_DOCKER_LINT_SECURITY_CHECK
🔨 Проверка безопасности (gosec)...
endef

define TAIL__MESSAGE_DOCKER_LINT_MOD_VERIFY
🔨 Проверка зависимостей...
endef

define TAIL__MESSAGE_DOCKER_LINT_SUCCESS
✅ Все проверки завершены!
endef

define TAIL__MESSAGE_DOCKER_TEST_RUN
🔨 запуск тестов
endef

define TAIL__MESSAGE_DOCKER_TEST_SUCCESS
✅ Тестирование успешно завершено!
endef

define TAIL__MESSAGE_DOCKER_INTEGRATION_TEST_RUN
🔨 компиляция тестовых скриптов и запуск интеграционных тестов
endef

define TAIL__MESSAGE_DOCKER_INTEGRATION_TEST_SUCCESS
✅ Тестирование успешно завершено!
endef

define TAIL__MESSAGE_DOCKER_BUILD_INTEGRATION_EXPECTED
🔨 компиляция тестовых скриптов и создание файлов логов для сравнения
endef

define TAIL__MESSAGE_DOCKER_BUILD_INTEGRATION_EXPECTED_SUCCESS
✅ Создание успешно завершено!
endef