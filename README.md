# tail

Advanced Go-based `tail` utility for real-time monitoring of files, commands, and pipelines (pipes) with ANSI sequence support, process indicators, and flexible output configuration.

## ✨ Key Features

*   **Real-time Monitoring**: Track output from commands, files, or stdin data.
*   **Full ANSI Support**: Correct handling of control sequences (CSI, SGR), including colors (256/RGB), cursor movement, and line clearing.
*   **Output Modes**:
    *   `direct`: Immediate output upon receiving new data.
    *   `thread`: Periodic redrawing (reduces flicker during high-intensity flows).
    *   `roll`: Circular buffer mode.
*   **Process Indicators**: Visual indicators of activity (roller, pipe, process, etc.).
*   **Flexible Templating**: Three levels of output detail (`none`, `minimal`, `full`).
*   **Smart Cleanup**: Automatically removes temporary logs upon successful command completion.
*   **Localization**: Full support for English and Russian languages.

## 🚀 Installation and Build

The project uses a Docker-based build system via `Makefile`.

### For Windows (via WSL2):
```powershell
wsl.exe make build
```
The binary will be compiled to the `./app/tail` folder.

### For Linux/macOS:
```bash
make build
```

## 🛠 Usage

### Basic Syntax
```bash
tail [flags] [-- command arguments]
```

### Examples
1.  **Monitor a command with 10 lines of output**:
    ```bash
    tail -n 10 -o thread -t minimal -- ping google.com
    ```
2.  **Use in a pipeline (pipe)**:
    ```bash
    cat log.txt | tail -o direct -t none
    ```
3.  **With a process indicator and title**:
    ```bash
    tail -a "MyTask" -r roller -c "long_running_script.sh"
    ```

## ⚙️ Configuration Flags

| Flag | Short | Description | Default |
| :--- | :---: | :--- | :--- |
| `--lines` | `-n` | Max number of lines to display on screen (up to 250) | `5` |
| `--output` | `-o` | Output mode: `direct`, `thread`, `roll` | `direct` |
| `--template` | `-t` | Output template: `none`, `minimal`, `full` | `minimal` |
| `--indicator`| `-r` | Indicator type: `none`, `roller`, `roll5`, `bolded`, `pipe`, `process`, etc. | `none` |
| `--command` | `-c` | Command to execute and monitor | - |
| `--title` | `-a` | Process name in the title | - |
| `--icon` | `-i` | Process icon in the title | - |
| `--length` | `-l` | Maximum line length (0 - no limit) | `0` |
| `--size` | `-s` | Storage buffer size in lines | `10240` |
| `--full` | `-f` | Full output mode | `false` |
| `--version` | `-v` | Show program version | `false` |
| `--help` | `-h` | Show help | `false` |

## 🌍 Localization

The program automatically detects the system language. You can force a language change via the `LANG` environment variable:

```bash
# Russian interface
LANG=ru_RU.UTF-8 ./tail ...

# English interface
LANG=en_US.UTF-8 ./tail ...
```

## ⚠️ Limitations

*   **Line Limit**: The maximum number of lines displayed on the screen is limited to 250 to prevent terminal overflow.
*   **Performance**: Extremely high volumes of ANSI codes in `direct` mode may cause increased CPU load (using `-o thread` is recommended).

## License

This project is licensed under the **MIT License**.

You are free to use, modify, and distribute this software, provided that the original copyright notice and permission notice are included in all copies or substantial portions of the software.

### Attribution

If you use this project or its code, please provide attribution by mentioning the author:
- **Author**: Makoveev Vitalii
- **Source**: [https://github.com/makoveev/tail](https://github.com/makoveev/tail)

---
### Development

To run tests and linters:
```bash
wsl.exe make test
wsl.exe make lint
wsl.exe make integration-test
```
