package line

type lines []Line

func NewLines(count uint64) Lines {
	var (
		ll lines
	)
	if count == 0 {
		ll = make([]Line, count)
	} else {
		ll = newLines(count)
	}

	return &ll
}

func newLines(count uint64) []Line {
	ll := make([]Line, count)

	for i := uint64(0); i < count; i++ {
		ll[i] = NewLine()
	}

	return ll
}

func (l *lines) Add(count uint64) {
	if count == 0 {
		return
	}

	*l = append(*l, newLines(count)...)
}

func (l *lines) Get(id uint64) Line {
	return (*l)[id]
}

func (l *lines) CleanPostfix(id uint64) {
	if id < l.Len() && id > 0 {
		*l = (*l)[:id]
	}
	if id == 0 {
		*l = lines{}
	}
}

func (l *lines) CleanPrefix(id uint64) {
	if id < l.Len() {
		*l = (*l)[id:]
	} else {
		*l = lines{}
	}
}

func (l *lines) CleanString(id uint64) {
	(*l)[id] = NewLine()
}

func (l *lines) Len() uint64 {
	return uint64(len(*l))
}
