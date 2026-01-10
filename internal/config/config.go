package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type config struct {
	countLines  uint64 //@todo вычислить количество строк доступных для вывода
	lengthLines int    //@todo
	sizeBuffer  uint64 //@todo максимальное количество строк в буфере
	processName string
	processIcon string //@todo указание иконки
	outputMode  StringValue
	template    StringValue
	indicator   StringValue
	help        bool
	version     bool
	command     string
	args        []string
}

func ReadConf() (Cfg, error) {
	cfg := config{
		outputMode: &defaultString{},
		template:   &defaultString{},
		indicator:  &defaultString{},
	}

	cmd := &cobra.Command{
		Use:           "tail",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cmd.SetArgs(os.Args[1:])

	cmd.Flags().Uint64VarP(&cfg.countLines, "lines", "n", 5, "maximum number of lines for output")
	cmd.Flags().IntVarP(&cfg.lengthLines, "length", "l", 0, "maximum length of the output line")
	cmd.Flags().Uint64VarP(&cfg.sizeBuffer, "size", "s", 500, "maximum size of the stored buffer in number of lines")
	cmd.Flags().StringVarP(&cfg.processName, "title", "a", "", "Name in the program title")
	cmd.Flags().StringVarP(&cfg.processIcon, "icon", "i", "", "Icon in the program title")
	cmd.Flags().VarP(cfg.outputMode, "output", "o", "text output process option")
	cmd.Flags().VarP(cfg.template, "template", "t", "data output template")
	cmd.Flags().VarP(cfg.indicator, "indicator", "r", "specifying an element to display the program's running process")
	cmd.Flags().StringVarP(&cfg.command, "command", "c", "", "the address of the executable file to be tracked")

	if err := cmd.Execute(); err != nil {
		return nil, err
	}
	cfg.args = cmd.Flags().Args()

	if cfg.countLines > 250 {
		_, _ = fmt.Fprintf(os.Stderr, "Error: n must be positive and < 250\n")
		os.Exit(1)
	}

	return &cfg, nil
}

func (c *config) GetCountLines() uint64 {
	return c.countLines
}
func (c *config) GetLengthLines() int {
	return c.lengthLines
}
func (c *config) GetSizeBuffer() uint64 {
	return c.sizeBuffer
}

func (c *config) GetProcessName() string {
	return c.processName
}
func (c *config) SetProcessName(str string) {
	c.processName = str
}
func (c *config) GetProcessIcon() string {
	return c.processIcon
}
func (c *config) SetProcessIcon(str string) {
	c.processIcon = str
}

func (c *config) IsHelp() bool {
	return c.help
} //@todo
func (c *config) IsVersion() bool {
	return c.version
} //@todo
func (c *config) GetCommand() string {
	return c.command
}
func (c *config) GetArgs() []string {
	return c.args
}

func (c *config) ReplaceOutputMode(container StringValue) error {
	err := container.Set(c.outputMode.String())
	if err == nil {
		c.outputMode = container
	}

	return err
}
func (c *config) GetOutputMode() string {
	return c.outputMode.String()
}

func (c *config) ReplaceTemplate(container StringValue) error {
	err := container.Set(c.template.String())
	if err == nil {
		c.template = container
	}

	return err
}
func (c *config) GetTemplate() string {
	return c.template.String()
}

func (c *config) ReplaceIndicator(container StringValue) error {
	err := container.Set(c.indicator.String())
	if err == nil {
		c.indicator = container
	}

	return err
}
func (c *config) GetIndicator() string {
	return c.indicator.String()
}
