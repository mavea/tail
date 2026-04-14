package window

import "sync"

type Window struct {
	icon               string
	title              string
	countBufferLines   uint64
	countBufferColumns uint64
	x                  uint64
	y                  uint64
	maxHeight          uint64
	height             uint64
	maxWidth           uint64
	width              uint64

	mu sync.RWMutex
}

func NewWindow(icon, title string) *Window {
	return &Window{
		icon:               icon,
		title:              title,
		countBufferLines:   0,
		countBufferColumns: 0,
		x:                  0,
		y:                  0,
	}
}

func (w *Window) SetPosition(x int, y uint64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if x < 0 {
		w.x = 0
	} else {
		w.x = uint64(x)
	}

	w.y = y
}

func (w *Window) GetPosition() (uint64, uint64) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.x, w.y
}

func (w *Window) SetBufferSize(countBufferLines uint64, countBufferColumns int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.countBufferLines = countBufferLines

	if countBufferColumns <= 0 {
		w.countBufferColumns = 0
	} else {
		w.countBufferColumns = uint64(countBufferColumns)
	}
}

func (w *Window) GetBufferSize() (uint64, uint64) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.countBufferLines, w.countBufferColumns
}

func (w *Window) SetIcon(icon string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.icon = icon
}

func (w *Window) GetIcon() string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.icon
}

func (w *Window) SetTitle(title string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.title = title
}

func (w *Window) GetTitle() string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.title
}

func (w *Window) SetMaxSize(width, height int) {
	if width <= 0 {
		width = 1
	}
	if height <= 0 {
		height = 1
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.maxWidth = uint64(width)
	w.width = w.maxWidth
	w.maxHeight = uint64(height)
	w.height = w.maxHeight
}
func (w *Window) Height() uint64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.height
}

func (w *Window) Width() uint64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.width
}
