package line

import "strings"

type line []string

func NewLine() Line {
	return &line{}
}

func (l *line) String(clean string, length int) string {
	str := ""
	style := ""
	leng := 0
	for i := 0; i+1 < len(*l); i += 2 {
		leng += len((*l)[i+1])
		if style == (*l)[i] {
			str += (*l)[i+1]
		} else {
			style = (*l)[i]
			if (*l)[i] == "" {
				str += clean + (*l)[i+1]
			} else {
				str += style + (*l)[i+1]
			}
		}
		if leng >= length && length > 0 {
			return str[:len(str)+length-leng] + " ."
		}
	}

	return str
}

func (l *line) skipToX(x int) (int, int) {
	return l.skipToXFromY(x, 1)
}

func (l *line) skipToXFromY(x, y int) (int, int) {
	var (
		length = len(*l)
	)
	if x <= 0 {
		return y, 0
	}
	for ; y < length; y += 2 {
		x -= len((*l)[y])
		if x <= 0 {
			return y, len((*l)[y]) + x
		}
	}

	return y, -x
}

func (l *line) addToEnd(style string, add string) {
	if len(*l) != 0 {
		if (*l)[len(*l)-2] == style {
			(*l)[len(*l)-1] += add

			return
		}
	}

	*l = append(*l, style, add)
}

func (l *line) inject(style string, add string, y, x int) {
	if x+len(add) <= len((*l)[y]) {
		str := (*l)[y]
		if (*l)[y-1] == style {
			(*l)[y] = str[:x] + add + str[x+len(add):]

			return
		}
		(*l)[y] = str[:x]
		if y+1 < len(*l) {
			*l = append((*l)[:y+1], append([]string{style, add, (*l)[y-1], str[x+len(add):]}, (*l)[y+1:]...)...)
		} else {
			*l = append((*l)[:y+1], []string{style, add, (*l)[y-1], str[x+len(add):]}...)
		}

		return
	}
	offset := len((*l)[y]) - x
	if offset > 0 {
		(*l)[y] = (*l)[y][:x]
	}
	ty, tx := l.skipToXFromY(len(add)-offset, y+2)

	if tx < 0 || ty > len(*l) {
		(*l) = (*l)[:y+1]
		l.addToEnd(style, add)

		return
	}

	if tx > 0 {
		(*l)[ty] = (*l)[ty][tx:]
	}

	if (*l)[y-1] == style {
		(*l)[y] += add
		*l = append((*l)[:y+1], (*l)[ty-1:]...)

		return
	}
	if (*l)[ty-1] == style {
		(*l)[ty] = add + (*l)[ty]

		*l = append((*l)[:y+1], (*l)[ty-1:]...)

		return
	}

	*l = append((*l)[:y+1], append([]string{style, add}, (*l)[ty-1:]...)...)
}

func (l *line) Set(style string, add string, x int) int {
	var (
		result = x + len(add)
		y      int
		length = len(*l)
	)

	y, x = l.skipToX(x)
	if x < 0 {
		l.addToEnd(style, strings.Repeat(" ", -x)+add)

		return result
	}
	if y < length {
		l.inject(style, add, y, x)
	} else {
		l.addToEnd(style, add)
	}

	return result
}

func (l *line) CleanPrefix(x int) int {
	ll := *l
	if x > 0 {
		for i := 1; i < len(ll); i += 2 {
			if len(ll[i]) < x {
				x -= len(ll[i])
			} else {
				if len(ll[i]) == x {
					if len(ll) > i+1 {
						(*l) = ll[i+1:]
					} else {
						(*l) = line{}
					}
					return 0
				}
				ll[i] = ll[i][x:]
				(*l) = ll[i-1:]

				return 0
			}
		}
	}
	(*l) = line{}

	return 0

}

func (l *line) CleanPostfix(x int) int {
	if x <= 0 {
		*l = line{}

		return 0
	}
	inX := x
	ll := *l
	for i := 1; i < len(ll); i += 2 {
		if len(ll[i]) < x {
			x -= len(ll[i])
		} else {
			if len(ll[i]) > x {
				(*l)[i] = ll[i][:x]
			}
			(*l) = ll[:i+1]

			return inX
		}
	}

	return inX
}
