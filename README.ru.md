# tail

Продвинутая утилита `tail` на языке Go для мониторинга файлов, команд и конвейеров (pipe) в реальном времени с поддержкой ANSI-последовательностей, индикаторов процесса и гибкой настройки вывода.

## ✨ Ключевые особенности

*   **Мониторинг в реальном времени**: Отслеживание вывода команд, файлов или данных из стандартного ввода (stdin).
*   **Полная поддержка ANSI**: Корректная обработка управляющих последовательностей (CSI, SGR), включая цвета (256/RGB), перемещение курсора и очистку строк.
*   **Режимы вывода**:
    *   `direct`: Немедленный вывод при поступлении новых данных.
    *   `thread`: Периодическая перерисовка (уменьшает мерцание при интенсивном потоке).
    *   `roll`: Режим циклического буфера.
*   **Индикаторы процесса**: Наглядная визуализация работы (roller, pipe, process и др.).
*   **Гибкая шаблонизация**: Три уровня детализации вывода (`none`, `minimal`, `full`).
*   **Умная очистка**: Автоматическое удаление временных логов при успешном завершении команды.
*   **Локализация**: Полная поддержка русского и английского языков.

## 🚀 Установка и сборка

Проект использует Docker-систему сборки через `Makefile`.

### Для Windows (через WSL2):
```powershell
wsl.exe make build
```
Бинарный файл будет скомпилирован в папку `./app/tail`.

### Для Linux/macOS:
```bash
make build
```

## 🛠 Использование

### Основной синтаксис
```bash
tail [флаги] [-- команда аргументы]
```

### Примеры
1.  **Мониторинг выполнения команды с 10 строками вывода**:
    ```bash
    tail -n 10 -o thread -t minimal -- ping google.com
    ```
2.  **Использование в конвейере (pipe)**:
    ```bash
    cat log.txt | tail -o direct -t none
    ```
3.  **С индикатором процесса и заголовком**:
    ```bash
    tail -a "MyTask" -r roller -c "long_running_script.sh"
    ```

## ⚙️ Флаги конфигурации

| Флаг | Короткий | Описание | По умолчанию |
| :--- | :---: | :--- | :--- |
| `--lines` | `-n` | Макс. кол-во строк для вывода на экран (до 250) | `5` |
| `--output` | `-o` | Режим вывода: `direct`, `thread`, `roll` | `direct` |
| `--template` | `-t` | Шаблон вывода: `none`, `minimal`, `full` | `minimal` |
| `--indicator`| `-r` | Тип индикатора: `none`, `roller`, `roll5`, `bolded`, `pipe`, `process` и др. | `none` |
| `--command` | `-c` | Команда для выполнения и мониторинга | - |
| `--title` | `-a` | Название процесса в заголовке | - |
| `--icon` | `-i` | Иконка процесса в заголовке | - |
| `--length` | `-l` | Максимальная длина строки (0 - без ограничений) | `0` |
| `--size` | `-s` | Размер буфера хранения строк | `10240` |
| `--full` | `-f` | Режим полного вывода | `false` |
| `--version` | `-v` | Показать версию программы | `false` |
| `--help` | `-h` | Показать справку | `false` |

## 🌍 Локализация

Программа автоматически определяет язык системы. Вы можете принудительно изменить язык через переменную окружения `LANG`:

```bash
# Русский интерфейс
LANG=ru_RU.UTF-8 ./tail ...

# Английский интерфейс
LANG=en_US.UTF-8 ./tail ...
```

## ⚠️ Ограничения

*   **Лимит строк**: Максимальное количество отображаемых на экране строк ограничено 250 для предотвращения переполнения терминала.
*   **Производительность**: При экстремально больших объемах ANSI-кодов в режиме `direct` возможна повышенная нагрузка на CPU (рекомендуется использовать `-o thread`).

## 👨‍💻 Автор

*   **Makoveev Vitalii**
*   **Source**: [https://github.com/makoveev/tail](https://github.com/makoveev/tail)

## License

This project is licensed under the **MIT License**.

You are free to use, modify, and distribute this software, provided that the original copyright notice and permission notice are included in all copies or substantial portions of the software.

### Attribution

If you use this project or its code, please provide attribution by mentioning the author:
- **Author**: Makoveev Vitalii
- **Source**: [https://github.com/makoveev/tail](https://github.com/makoveev/tail)

---
### Разработка

Для запуска тестов и линтеров:
```bash
wsl.exe make test
wsl.exe make lint
wsl.exe make integration-test
```

