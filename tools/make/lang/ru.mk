define TAIL__MESSAGE_DOCKER_BUILD_MOD_DOWNLOAD
📥 Загружаем зависимости...
endef

define TAIL__MESSAGE_DOCKER_BUILD_GO_COMPILING
🔨 Компилируем...
endef

define TAIL__MESSAGE_DOCKER_BUILD_SUCCESS
✅ Компиляция завершена!
endef

define TAIL__MESSAGE_DOCKER_BUILD_FAIL
❌ Ошибка компеляции
endef

define TAIL__MESSAGE_DOCKER_DEBUG_PREPARE
🐞 Подготавливаем debug-окружение и Delve...
endef

define TAIL__MESSAGE_DOCKER_DEBUG_BUILD
🔨 Собираем debug-бинарник без оптимизаций...
endef

define TAIL__MESSAGE_DOCKER_DEBUG_FAIL
❌ Ошибка debug сборки
endef


define TAIL__MESSAGE_DOCKER_DEBUG_WAIT
🧭 Delve ожидает подключение на порту 2345...
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
✅ Все проверки завершены успешно!
endef

define TAIL__MESSAGE_DOCKER_LINT_FAIL
❌ Проверки провалены!
endef

define TAIL__MESSAGE_DOCKER_TEST_RUN
🔨 Запуск unit-тестов...
endef

define TAIL__MESSAGE_DOCKER_TEST_SUCCESS
✅ Тестирование успешно завершено!
endef

define TAIL__MESSAGE_DOCKER_TEST_FAIL
❌ Тестирование провалено!
endef

define TAIL__MESSAGE_DOCKER_INTEGRATION_TEST_RUN
🔨 Компиляция тестовых скриптов и запуск интеграционных тестов...
endef

define TAIL__MESSAGE_DOCKER_INTEGRATION_TEST_SUCCESS
✅ Тестирование успешно завершено!
endef

define TAIL__MESSAGE_DOCKER_INTEGRATION_TEST_FAIL
❌ Тестирование провалено!
endef

define TAIL__MESSAGE_DOCKER_BUILD_INTEGRATION_EXPECTED
🔨 Компиляция тестовых скриптов и создание файлов логов для сравнения...
endef

define TAIL__MESSAGE_DOCKER_BUILD_INTEGRATION_EXPECTED_SUCCESS
✅ Создание успешно завершено!
endef

define TAIL__MESSAGE_DOCKER_GENERATE_RUN
🔨 Генерация моков...
endef

define TAIL__MESSAGE_DOCKER_GENERATE_SUCCESS
✅ Генерация успешно завершена!
endef

define TAIL__MESSAGE_DOCKER_GENERATE_FAIL
❌ Генерация провалена!
endef

define TAIL__MESSAGE_DOCKER_GENERATE_CHECK_RUN
🔨 Проверка актуальности сгенерированных моков...
endef

define TAIL__MESSAGE_DOCKER_GENERATE_CHECK_SUCCESS
✅ Сгенерированные моки актуальны!
endef

define TAIL__MESSAGE_DOCKER_GENERATE_CHECK_FAIL
❌ Сгенерированные моки устарели! Выполните make generate.
endef
