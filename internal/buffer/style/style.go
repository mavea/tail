package style

import (
	"strconv"
)

type style struct {
	backColor uint64
	backRGB   string
	textColor uint64
	textRGB   string
	style     []bool
}

func New() Style {
	return newStyle()
}

func newStyle() *style {
	return &style{
		style: make([]bool, 8),
	}
}

func (s *style) copy() *style {
	st := &style{
		backColor: s.backColor,
		backRGB:   s.backRGB,
		textColor: s.textColor,
		textRGB:   s.textRGB,
		style:     make([]bool, 8),
	}
	copy(st.style, s.style)

	return st
}

func (s *style) Set(data []uint64) Style {
	st := s.copy()
	if len(data) > 0 {
		switch data[0] {
		case 38, 48:
			if len(data) > 1 && data[1] == 2 {
				var r, g, b uint64
				if len(data) > 2 {
					r = data[2]
				}
				if len(data) > 3 {
					g = data[3]
				}
				if len(data) > 4 {
					b = data[4]
				}
				if data[0] == 38 {
					st.textRGB = strconv.FormatUint(r, 10) + ";" + strconv.FormatUint(g, 10) + ";" + strconv.FormatUint(b, 10)
					st.textColor = 0
				} else {
					st.backRGB = strconv.FormatUint(r, 10) + ";" + strconv.FormatUint(g, 10) + ";" + strconv.FormatUint(b, 10)
					st.backColor = 0
				}
			}
		default:
			for _, n := range data {
				switch true {
				case n == 0:
					st = newStyle()
				case n > 0 && n <= 9:
					st.style[n-1] = true
				case n >= 30 && n <= 37 || n >= 90 && n <= 97:
					st.textColor = n
					st.textRGB = ""
				case n >= 40 && n <= 47 || n >= 100 && n <= 107:
					st.backColor = n
					st.backRGB = ""
				}
			}
		}
	}

	return st
}

func (s *style) getStyle() string {
	str := ""
	for i := 0; i < len(s.style); i++ {
		if s.style[i] {
			if str == "" {
				str = strconv.Itoa(i + 1)
			} else {
				str += ";" + strconv.Itoa(i+1)
			}
		}
	}

	return str
}

func (s *style) String() string {
	str := ""
	rgb := ""
	if s.backRGB != "" {
		rgb = "\033[48;2;" + s.backRGB + "m"
	} else if s.backColor != 0 {
		str = strconv.FormatUint(s.backColor, 10)
	}

	if s.textRGB != "" {
		rgb += "\033[38;2;" + s.textRGB + "m"
	} else if s.textColor != 0 {
		if str == "" {
			str = strconv.FormatUint(s.textColor, 10)
		} else {
			str += ";" + strconv.FormatUint(s.textColor, 10)
		}
	}

	add := s.getStyle()
	if str != "" && add != "" {
		str += ";" + add
	} else if str != add {
		str += add
	}

	if rgb != "" {
		if str != "" {
			return rgb + "\033[" + str + "m"
		}

		return rgb
	} else {
		if str != "" {
			return "\033[" + str + "m"
		}

		return ""
	}
}
