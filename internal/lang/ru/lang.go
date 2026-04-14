package ru

import langGeneral "tail/internal/lang/general"

const (
	Code langGeneral.Code = "ru"
)

func NewLang() *langGeneral.Lang {
	return &langGeneral.Lang{
		MaxLineCount:          "максимальное количество строк для вывода на экран",
		MaxCharsPerLine:       "максимальное количество символов в строке вывода",
		MaxBufferLines:        "максимальное количество строк, хранимых в обрабатываемом буфере",
		ProcessName:           "название в заголовке программы",
		ProcessIcon:           "иконка в заголовке программы",
		OutputMode:            "режим вывода информации (режимы: direct - вывод при новых данных, thread - периодическая перерисовка)",
		OutputTemplate:        "шаблон вывода информации (none, minimal, full)",
		Indicator:             "индикатор процесса (none, roller, bolded-roller, pipe)",
		Command:               "команда, выполнение которой мониторит процесс. Как альтернатива возможен запуск через pipe",
		Help:                  "отображение информации о помощи",
		Version:               "отображение информации о версии",
		HelpDescription:       "Tail\n Программа для мониторинга выполнения команд и отображения их вывода в реальном времени с различными режимами и шаблонами отображения. Поддерживает как запуск команд напрямую, так и получение данных через pipe. Программа удаляет лог в случае его успешного завершения, обеспечивая чистоту и актуальность отображаемой информации.\n",
		HelpBottomDescription: "Для изменения языка, при запуске укажите переменную LANG в формате ru_RU.UTF-8, или en_US.UTF-8 (важны первые 2 символа)\n",
		HelpExample: "Примеры использования:\n" +
			"  tail -n 10 -o thread -t minimal -- ping google.com\n" +
			"  tail -s 1000 -o direct -t full -- ping google.com\n" +
			"  echo 'Hello World' | tail -o direct -t none\n",
		VersionDescription: "Tail\n Исходный код:    https://github.com/mavea/tail \n  Версия:          %s\n  Автор:           %s\n  Дата компиляции: %s\n",
		Full:               "полный вывод",
	}
}
