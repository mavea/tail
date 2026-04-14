package style

import (
	"strconv"
)

type style struct {
	backColor uint64
	backRGB   string
	back256   uint64
	isBack256 bool
	textColor uint64
	textRGB   string
	text256   uint64
	isText256 bool
	style     []bool
}

func New() Style {
	return newStyle()
}

func newStyle() *style {
	return &style{
		style: make([]bool, 9),
	}
}

func (s *style) copy() *style {
	st := &style{
		backColor: s.backColor,
		backRGB:   s.backRGB,
		back256:   s.back256,
		isBack256: s.isBack256,
		textColor: s.textColor,
		textRGB:   s.textRGB,
		text256:   s.text256,
		isText256: s.isText256,
		style:     make([]bool, 9),
	}
	copy(st.style, s.style)

	return st
}

func (s *style) Set(data []uint64) Style {
	st := s.copy()
	for i := 0; i < len(data); i++ {
		n := data[i]
		switch n {
		case 0:
			st = newStyle()
		case 38, 48:
			if i+2 < len(data) && data[i+1] == 5 {
				if n == 38 {
					st.text256 = data[i+2]
					st.isText256 = true
					st.textRGB = ""
					st.textColor = 0
				} else {
					st.back256 = data[i+2]
					st.isBack256 = true
					st.backRGB = ""
					st.backColor = 0
				}
				i += 2
				continue
			}

			if i+4 < len(data) && data[i+1] == 2 {
				r := data[i+2]
				g := data[i+3]
				b := data[i+4]
				if n == 38 {
					st.textRGB = strconv.FormatUint(r, 10) + ";" + strconv.FormatUint(g, 10) + ";" + strconv.FormatUint(b, 10)
					st.textColor = 0
					st.isText256 = false
				} else {
					st.backRGB = strconv.FormatUint(r, 10) + ";" + strconv.FormatUint(g, 10) + ";" + strconv.FormatUint(b, 10)
					st.backColor = 0
					st.isBack256 = false
				}
				i += 4
			}
		default:
			switch true {
			case n > 0 && n <= 9:
				st.style[n-1] = true
			case n == 39:
				st.textColor = 0
				st.textRGB = ""
				st.isText256 = false
			case n == 49:
				st.backColor = 0
				st.backRGB = ""
				st.isBack256 = false
			case n >= 30 && n <= 37 || n >= 90 && n <= 97:
				st.textColor = n
				st.textRGB = ""
				st.isText256 = false
			case n >= 40 && n <= 47 || n >= 100 && n <= 107:
				st.backColor = n
				st.backRGB = ""
				st.isBack256 = false
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
	color := ""
	if s.backRGB != "" {
		color = "\033[48;2;" + s.backRGB + "m"
	} else if s.isBack256 {
		color = "\033[48;5;" + strconv.FormatUint(s.back256, 10) + "m"
	} else if s.backColor != 0 {
		str = strconv.FormatUint(s.backColor, 10)
	}

	if s.textRGB != "" {
		color += "\033[38;2;" + s.textRGB + "m"
	} else if s.isText256 {
		color += "\033[38;5;" + strconv.FormatUint(s.text256, 10) + "m"
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

	if color != "" {
		if str != "" {
			return color + "\033[" + str + "m"
		}

		return color
	}

	if str != "" {
		return "\033[" + str + "m"
	}

	return ""
}
