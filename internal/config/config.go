package config

import (
	"fmt"
	"os"
	"tail/internal/formatter"

	"github.com/spf13/cobra"

	configGeneral "tail/internal/config/general"
	langGeneral "tail/internal/lang/general"
	sinkConsoleIndicator "tail/internal/sink/console/indicator"
	sinkConsoleTemplate "tail/internal/sink/console/template"
)

const (
	// MaxLineCount defines the maximum allowed value for line count
	maxLineCount       = 250
	maxCharsPerLine    = 1000
	maxBufferLines     = 1000000
	defaultBufferLines = 10 * 1024
)

type Cfg struct {
	maxLineCount    int
	maxCharsPerLine int
	maxBufferLines  uint64
	processName     string
	processIcon     string
	outputMode      configGeneral.StringValue
	outputTemplate  configGeneral.StringValue
	indicator       configGeneral.StringValue
	help            bool
	version         bool
	full            bool
	command         string
	args            []string
}

// ReadConf reads and parses the configuration from command-line arguments using the provided language.
func ReadConf(lang *langGeneral.Lang) (*Cfg, error) {
	cfg := Cfg{
		outputMode:     formatter.NewOutputMode(),
		outputTemplate: sinkConsoleTemplate.NewTemplateType(),
		indicator:      sinkConsoleIndicator.NewIndicatorType(),
	}

	return cfg.ReadArguments(&cobra.Command{
		Use:           "tail",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}, &cfg, lang)
}

// ReadArguments sets up command-line flags for the config using Cobra and executes the command.
func (c *Cfg) ReadArguments(cmd *cobra.Command, cfg *Cfg, lang *langGeneral.Lang) (*Cfg, error) {
	cmd.SetArgs(os.Args[1:])

	// Helper to add string flags
	addStringFlag := func(p *string, name, short, defaultVal string, desc string) {
		cmd.Flags().StringVarP(p, name, short, defaultVal, desc)
	}

	// Helper to add int flags
	addIntFlag := func(p *int, name, short string, defaultVal int, desc string) {
		cmd.Flags().IntVarP(p, name, short, defaultVal, desc)
	}

	// Helper to add int64 flags
	addUint64Flag := func(p *uint64, name, short string, defaultVal uint64, desc string) {
		cmd.Flags().Uint64VarP(p, name, short, defaultVal, desc)
	}

	// Helper to add var flags
	addVarFlag := func(p configGeneral.StringValue, name, short, desc string) {
		cmd.Flags().VarP(p, name, short, desc)
	}

	addIntFlag(&cfg.maxLineCount, "lines", "n", 5, lang.MaxLineCount.String())
	addIntFlag(&cfg.maxCharsPerLine, "length", "l", 0, lang.MaxCharsPerLine.String())
	addUint64Flag(&cfg.maxBufferLines, "size", "s", defaultBufferLines, lang.MaxBufferLines.String())
	addStringFlag(&cfg.processName, "title", "a", "", lang.ProcessName.String())
	addStringFlag(&cfg.processIcon, "icon", "i", "", lang.ProcessIcon.String())
	addVarFlag(cfg.outputMode, "output", "o", lang.OutputMode.String())
	addVarFlag(cfg.outputTemplate, "template", "t", lang.OutputTemplate.String())
	addVarFlag(cfg.indicator, "indicator", "r", lang.Indicator.String())
	addStringFlag(&cfg.command, "command", "c", "", lang.Command.String())
	cmd.PersistentFlags().BoolVarP(&cfg.help, "help", "h", false, lang.Help.String())
	cmd.PersistentFlags().BoolVarP(&cfg.version, "version", "v", false, lang.Version.String())
	cmd.PersistentFlags().BoolVarP(&cfg.full, "full", "f", false, lang.Full.String())

	if err := cmd.Execute(); err != nil {
		return nil, err
	}
	cfg.args = cmd.Flags().Args()

	return cfg, nil
}

func (c *Cfg) GetMaxLineCount() int {
	return c.maxLineCount
}

func (c *Cfg) SetMaxLineCount(n int) {
	c.maxLineCount = n
}

func (c *Cfg) GetMaxCharsPerLine() int {
	return c.maxCharsPerLine
}

func (c *Cfg) SetMaxCharsPerLine(n int) {
	c.maxCharsPerLine = n
}

func (c *Cfg) GetMaxBufferLines() uint64 {
	return c.maxBufferLines
}

func (c *Cfg) GetProcessName() string {
	return c.processName
}

func (c *Cfg) GetProcessIcon() string {
	return c.processIcon
}

func (c *Cfg) IsHelp() bool {
	return c.help
}

func (c *Cfg) IsVersion() bool {
	return c.version
}

func (c *Cfg) GetCommand() string {
	return c.command
}

func (c *Cfg) GetArgs() []string {
	return c.args
}

func (c *Cfg) GetOutputMode() string {
	return c.outputMode.String()
}

func (c *Cfg) GetOutputTemplate() string {
	return c.outputTemplate.String()
}

func (c *Cfg) GetIndicator() string {
	return c.indicator.String()
}

func (c *Cfg) IsCSIEnabled() bool {
	return false
}

func (c *Cfg) IsFullOutput() bool {
	return c.full
}

// Validate checks if the configuration values are within acceptable ranges.
func (c *Cfg) Validate() error {
	if c.maxLineCount < 1 {
		return fmt.Errorf("maxLineCount must be positive")
	}
	if c.maxLineCount > maxLineCount {
		return fmt.Errorf("maxLineCount must be < %d", maxLineCount)
	}

	if c.maxCharsPerLine < 0 {
		return fmt.Errorf("maxCharsPerLine must be non-negative")
	}
	if c.maxCharsPerLine > maxCharsPerLine {
		return fmt.Errorf("maxCharsPerLine must be <= %d", maxCharsPerLine)
	}

	if c.maxBufferLines < 1 {
		return fmt.Errorf("maxBufferLines must be positive")
	}
	if c.maxBufferLines > maxBufferLines {
		return fmt.Errorf("maxBufferLines must be <= %d", maxBufferLines)
	}

	if !c.outputMode.Validate() {
		return fmt.Errorf("output mode is invalid")
	}
	if !c.indicator.Validate() {
		return fmt.Errorf("indicator type is invalid")
	}
	if !c.outputTemplate.Validate() {
		return fmt.Errorf("output template is invalid")
	}

	if uint64(c.maxLineCount) >= c.maxBufferLines {
		return fmt.Errorf("lines must be < size")
	}

	if c.command != "" && len(c.args) > 0 {
		return fmt.Errorf("command mode and file arguments cannot be used together")
	}
	if len(c.args) > 1 {
		return fmt.Errorf("only one input file is supported")
	}

	return nil
}
