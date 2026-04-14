package buffer

import "testing"

type gfWindow struct{}

func (w *gfWindow) SetPosition(_ int, _ uint64)   {}
func (w *gfWindow) SetBufferSize(_ uint64, _ int) {}
func (w *gfWindow) SetIcon(_ string)              {}
func (w *gfWindow) SetTitle(_ string)             {}
func (w *gfWindow) Height() uint64                { return 24 }

type gfCfg struct{}

func (c *gfCfg) GetMaxLineCount() int      { return 5 }
func (c *gfCfg) GetMaxCharsPerLine() int   { return 20 }
func (c *gfCfg) GetMaxBufferLines() uint64 { return 10 }
func (c *gfCfg) GetProcessName() string    { return "p" }
func (c *gfCfg) GetProcessIcon() string    { return "i" }
func (c *gfCfg) IsCSIEnabled() bool        { return true }
func (c *gfCfg) IsFullOutput() bool        { return false }

func TestGetFull(t *testing.T) {
	b := New(&gfCfg{}, &gfWindow{})
	b.SetDefaultStyle("", "", "")
	b.Add("one")
	b.Add("two")
	full := b.GetFull()
	if len(full) == 0 {
		t.Fatal("expected non-empty full buffer")
	}
}
