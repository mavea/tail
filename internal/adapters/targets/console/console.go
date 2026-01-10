package console

import (
	"bufio"
	"fmt"
	"strings"
	"sync"
	"tail/internal/adapters/targets/console/indicator"
	"tail/internal/adapters/targets/console/templates"
	template_domain "tail/internal/adapters/targets/console/templates/domain"
	"tail/internal/adapters/targets/domain"
	"tail/pkg/typo"
)

type console struct {
	countNextCleanLines int
	writer              *bufio.Writer
	cfg                 cfg
	template            template_domain.Template
	out                 string
	clean               string

	mu sync.RWMutex
}

func New(cfg cfg, writer *bufio.Writer) (domain.Target, error) {

	if err := cfg.ReplaceTemplate(templates.NewType()); err != nil {
		return nil, err
	}
	if err := cfg.ReplaceIndicator(indicator.NewType()); err != nil {
		return nil, err
	}

	indicate := indicator.New(cfg.GetIndicator())
	template, err := templates.NewTemplate(cfg, indicate)
	if err != nil {
		return nil, err
	}

	cns := &console{
		cfg:      cfg,
		writer:   writer,
		template: template,
		out:      "",
		clean:    "",
	}

	return cns, nil
}

func (cns *console) flush() {
	_ = cns.writer.Flush()
}

func (cns *console) GetDefaultStyle() (string, string, string) {
	return cns.template.StartLine(), cns.template.CleanLine(), cns.template.EndLine()
}

func (cns *console) print() {
	_, _ = fmt.Fprint(cns.writer, cns.template.GetCellarClean(), cns.clean, cns.template.GetHeaderClean(), cns.template.GetHeader(), cns.out, cns.template.GetCellar(0, 0, 0, 0))
	cns.clean = strings.Repeat(typo.UpAndClean, cns.countNextCleanLines)
}

func (cns *console) Print() {
	cns.mu.Lock()
	defer cns.mu.Unlock()
	cns.print()
	cns.flush()
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

func (cns *console) ClearScreen() {
	cns.mu.RLock()
	defer cns.mu.RUnlock()
	_, _ = fmt.Fprint(cns.writer, cns.template.GetCellarClean(), cns.clean, cns.template.GetHeaderClean())
	cns.flush()
	cns.clean = strings.Repeat(typo.UpAndClean, 0)
}

func (cns *console) Error(buffer []string, err []string) {
	for _, line := range buffer {
		_, _ = fmt.Fprint(cns.writer, cns.template.FormatLine(line)+typo.NewLine)
	}
	_, _ = fmt.Fprint(cns.writer, "========= ERROR ========="+typo.NewLine)
	for _, line := range err {
		_, _ = fmt.Fprint(cns.writer, cns.template.FormatLine(line)+typo.NewLine)
	}
	cns.flush()
}

func (cns *console) SetStatus(status int) {
	_, _ = fmt.Fprint(cns.writer, "========= STATUS ========="+typo.NewLine, status, typo.NewLine)
	cns.flush()
}
