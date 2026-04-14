package console

import (
	"bufio"
	"context"
	"fmt"
	"strings"
	"sync"

	sinkConsoleIndicator "tail/internal/sink/console/indicator"
	sinkConsoleTemplate "tail/internal/sink/console/template"
	sinkConsoleTemplateGeneral "tail/internal/sink/console/template/general"
	sinkDomain "tail/internal/sink/domain"
	typo "tail/internal/typo"
)

type console struct {
	firstLine           bool
	countNextCleanLines int
	writer              *bufio.Writer
	cfg                 cfg
	template            sinkConsoleTemplateGeneral.Template
	out                 string
	clean               string

	mu sync.RWMutex
}

type CancelFunc func() error

func New(_ context.Context, cfg cfg, writer *bufio.Writer, window sinkConsoleTemplateGeneral.Window) (sinkDomain.Target, error) {
	indicate := sinkConsoleIndicator.New(cfg.GetIndicator())
	template, err := sinkConsoleTemplate.NewTemplate(cfg, indicate, window)
	if err != nil {
		return nil, err
	}

	cns := &console{
		cfg:       cfg,
		writer:    writer,
		template:  template,
		out:       "",
		clean:     "",
		firstLine: true,
	}

	return cns, nil
}

func (cns *console) flush() error {
	return cns.writer.Flush()
}

func (cns *console) GetDefaultStyle() (string, string, string) {
	return cns.template.StartLine(), cns.template.CleanLine(), cns.template.EndLine()
}

func (cns *console) print() error {
	_, err := fmt.Fprint(cns.writer, cns.template.GetCellarClean(cns.firstLine), cns.clean, cns.template.GetHeaderClean(cns.firstLine), cns.template.GetHeader(), cns.out, cns.template.GetCellar())
	cns.clean = strings.Repeat(typo.UpAndClean, cns.countNextCleanLines)
	return err
}

func (cns *console) Print() error {
	cns.mu.Lock()
	defer cns.mu.Unlock()
	if err := cns.print(); err != nil {
		return err
	}
	cns.firstLine = false
	return cns.flush()
}

func (cns *console) SetData(data []string) {
	str := ""
	for _, line := range data {
		str += cns.template.FormatLine(line) + typo.NewLine
	}
	cns.mu.Lock()
	defer cns.mu.Unlock()
	cns.countNextCleanLines = len(data)
	cns.out = str
}

func (cns *console) ClearScreen() error {
	cns.mu.Lock()
	defer cns.mu.Unlock()
	_, err := fmt.Fprint(cns.writer, cns.template.GetCellarClean(cns.firstLine), cns.clean, cns.template.GetHeaderClean(cns.firstLine))
	if err != nil {
		return err
	}
	cns.clean = strings.Repeat(typo.UpAndClean, 0)
	return cns.flush()
}

func (cns *console) PrintFull(buffer []string) error {
	for _, line := range buffer {
		if _, err := fmt.Fprint(cns.writer, cns.template.FormatLine(line)+typo.NewLine); err != nil {
			return err
		}
	}

	return cns.flush()
}
func (cns *console) Error(buffer []string, errs []string) error {
	for _, line := range buffer {
		if _, err := fmt.Fprint(cns.writer, cns.template.FormatLine(line)+typo.NewLine); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(cns.writer, "========= ERROR ========="+typo.NewLine); err != nil {
		return err
	}
	for _, line := range errs {
		if _, err := fmt.Fprint(cns.writer, cns.template.FormatLine(line)+typo.NewLine); err != nil {
			return err
		}
	}
	return cns.flush()
}

func (cns *console) SetStatus(status int) error {
	_, err := fmt.Fprint(cns.writer, "========= STATUS ========="+typo.NewLine, status, typo.NewLine)
	if err != nil {
		return err
	}
	return cns.flush()
}
