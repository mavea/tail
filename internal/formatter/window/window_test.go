package window

import "testing"

func TestWindowSettersAndGetters(t *testing.T) {
	w := NewWindow("icon", "title")

	w.SetPosition(-3, 7)
	x, y := w.GetPosition()
	if x != 0 || y != 7 {
		t.Fatalf("unexpected position: %d,%d", x, y)
	}

	w.SetPosition(11, 12)
	x, y = w.GetPosition()
	if x != 11 || y != 12 {
		t.Fatalf("unexpected position: %d,%d", x, y)
	}

	w.SetBufferSize(3, -1)
	lines, cols := w.GetBufferSize()
	if lines != 3 || cols != 0 {
		t.Fatalf("unexpected buffer size: %d,%d", lines, cols)
	}

	w.SetBufferSize(4, 5)
	lines, cols = w.GetBufferSize()
	if lines != 4 || cols != 5 {
		t.Fatalf("unexpected buffer size: %d,%d", lines, cols)
	}

	w.SetIcon("new-icon")
	w.SetTitle("new-title")
	if w.GetIcon() != "new-icon" {
		t.Fatalf("unexpected icon: %s", w.GetIcon())
	}
	if w.GetTitle() != "new-title" {
		t.Fatalf("unexpected title: %s", w.GetTitle())
	}

	w.SetMaxSize(80, 25)
	if w.Width() != 80 {
		t.Fatalf("unexpected width: %d", w.Width())
	}
	if w.Height() != 25 {
		t.Fatalf("unexpected height: %d", w.Height())
	}
}
