package ansi

import "testing"

func TestDECPrivateModesAllBranches(t *testing.T) {
	a := GetActionsDECPrivateModes()

	if _, ok := a.GetFunction('x'); ok {
		t.Fatal("unexpected function for unknown key")
	}

	f, ok := a.GetFunction('h')
	if !ok {
		t.Fatal("expected function for h")
	}
	if _, ok := f([]uint64{7}); !ok {
		t.Fatal("expected valid mode for h")
	}
	if _, ok := f([]uint64{999}); ok {
		t.Fatal("expected invalid mode for h")
	}

	f, ok = a.GetFunction('l')
	if !ok {
		t.Fatal("expected function for l")
	}
	if _, ok := f([]uint64{25}); !ok {
		t.Fatal("expected valid mode for l")
	}

	f, ok = a.GetFunction('c')
	if !ok {
		t.Fatal("expected function for c")
	}
	if _, ok := f([]uint64{1}); !ok {
		t.Fatal("expected valid VT100 mode")
	}
}

func TestSecondaryDABranches(t *testing.T) {
	a := GetActionsSecondaryDA()
	if _, ok := a.GetFunction('x'); ok {
		t.Fatal("unexpected function for unknown key")
	}

	f, ok := a.GetFunction('c')
	if !ok {
		t.Fatal("expected function c")
	}
	if _, ok := f([]uint64{1}); !ok {
		t.Fatal("expected valid c action")
	}

	f, ok = a.GetFunction('m')
	if !ok {
		t.Fatal("expected function m")
	}
	if _, ok := f([]uint64{1}); !ok {
		t.Fatal("expected valid m action")
	}
}

func TestWindowAndPalettesBranches(t *testing.T) {
	a := GetActionsWindowAndPalettes()
	if _, ok := a.GetFunction('x'); ok {
		t.Fatal("unexpected function for unknown key")
	}

	for _, ch := range []rune{'?', 's', 't', 5} {
		f, ok := a.GetFunction(ch)
		if !ok {
			t.Fatalf("expected function for %q", ch)
		}
		if ch == 5 {
			if _, ok := f(nil, "text"); !ok {
				t.Fatal("expected text action")
			}
			continue
		}
		if _, ok := f([]uint64{1, 2}, "title"); !ok {
			t.Fatalf("expected valid action for %q", ch)
		}
	}
}

func TestMainListAllKeysAndValidation(t *testing.T) {
	a := GetActionsMainList()
	if _, ok := a.GetFunction('x'); ok {
		t.Fatal("unexpected function for unknown key")
	}

	keys := []rune{'A', 'B', 'C', 'D', 'E', 'F', 'H', 'f', 'G', 'd', 'n', 's', 'u', 'J', 'K', 'm'}
	for _, key := range keys {
		f, ok := a.GetFunction(key)
		if !ok {
			t.Fatalf("missing function for key %q", key)
		}
		vals := []uint64{1}
		switch key {
		case 'n':
			vals = []uint64{5}
		case 's', 'u', 'J', 'K':
			vals = []uint64{0}
		case 'm':
			vals = []uint64{31}
		}
		if _, ok := f(vals); !ok {
			t.Fatalf("expected valid command for key %q", key)
		}
	}

	if f, ok := a.GetFunction('n'); ok {
		if _, ok := f([]uint64{99}); ok {
			t.Fatal("expected invalid status query")
		}
	}
	if f, ok := a.GetFunction('s'); ok {
		if _, ok := f([]uint64{1}); ok {
			t.Fatal("expected invalid save cursor arg")
		}
	}
	if f, ok := a.GetFunction('u'); ok {
		if _, ok := f([]uint64{1}); ok {
			t.Fatal("expected invalid restore cursor arg")
		}
	}
	if f, ok := a.GetFunction('J'); ok {
		if _, ok := f([]uint64{9}); ok {
			t.Fatal("expected invalid clear screen arg")
		}
	}
	if f, ok := a.GetFunction('K'); ok {
		if _, ok := f([]uint64{9}); ok {
			t.Fatal("expected invalid clear line arg")
		}
	}
	if f, ok := a.GetFunction('m'); ok {
		if _, ok := f(nil); ok {
			t.Fatal("expected invalid empty color args")
		}
		if _, ok := f([]uint64{38, 5, 1}); !ok {
			t.Fatal("expected valid 256-color args")
		}
		if _, ok := f([]uint64{38, 5, 255}); !ok {
			t.Fatal("expected valid 256-color upper bound")
		}
		if _, ok := f([]uint64{48, 5, 0}); !ok {
			t.Fatal("expected valid 256-color lower bound")
		}
		if _, ok := f([]uint64{38, 5, 256}); ok {
			t.Fatal("expected invalid 256-color overflow")
		}
		if _, ok := f([]uint64{38, 5}); ok {
			t.Fatal("expected invalid incomplete 256-color sequence")
		}
		if _, ok := f([]uint64{48, 2, 1, 2, 3}); !ok {
			t.Fatal("expected valid truecolor args")
		}
		if _, ok := f([]uint64{48, 2, 255, 255, 255}); !ok {
			t.Fatal("expected valid truecolor upper bound")
		}
		if _, ok := f([]uint64{48, 2, 255, 255, 256}); ok {
			t.Fatal("expected invalid truecolor overflow")
		}
		if _, ok := f([]uint64{38, 5, 202, 1}); !ok {
			t.Fatal("expected valid combined color and style")
		}
		if _, ok := f([]uint64{1, 38, 5, 202}); !ok {
			t.Fatal("expected valid prefixed style and color")
		}
	}
}

func TestOtherListAndGetLists(t *testing.T) {
	o := GetActionsOtherList()
	if _, ok := o.GetFunction('S'); ok {
		t.Fatal("other list should not expose functions")
	}
	if len(o) != 0 {
		t.Fatalf("expected empty other list, got %d", len(o))
	}

	d, s, w, m, other := GetLists()
	if len(d) == 0 || len(s) == 0 || len(w) == 0 || len(m) == 0 {
		t.Fatal("expected all list groups to be initialized")
	}
	if len(other) != 0 {
		t.Fatal("expected empty other list from GetLists")
	}
}
