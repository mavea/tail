package indicator

type ind struct {
	icons  []string
	cleans []string
	pos    int
}

func (f *ind) Get() string {
	f.pos++
	if f.pos >= len(f.icons) {
		f.pos = 0
	}

	return f.icons[f.pos]
}

func (f *ind) Clean() string {
	if f.pos < 0 {
		return ""
	}
	return f.cleans[f.pos]
}

func New(name string) Indicator {
	result := &ind{
		pos: -1,
	}
	switch name {
	case "roll":
		result.icons = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
		result.cleans = []string{"", "", "", "", "", "", "", ""}
	default:
		result.icons = []string{""}
		result.cleans = []string{""}
	}

	return result
}
