package line

import "math"

const maxBufferSize = math.MaxUint64 - math.MaxInt - 1

type lines struct {
	// buffer - фиксированный буфер строк.
	buffer []Line

	// cap - максимальное число сохраняемых строк.
	cap uint64

	// start - физический индекс первой логической строки.
	start uint64

	// history - количество добавленных строк за всё время.
	// Это не может быть больше cap, так как мы перезаписываем старые строки, но может быть меньше,
	// если мы ещё не заполнили буфер.
	history uint64
}

func MakeLines(max uint64) Lines {
	if max == 0 {
		max = 1
	}
	if max > maxBufferSize {
		max = maxBufferSize
	}

	buf := make([]Line, max)
	buf[0] = NewLine()

	return &lines{
		buffer:  buf,
		cap:     max,
		start:   0,
		history: 0,
	}
}

func (l *lines) reset() {
	l.start = l.physicalIndex(l.history)
	l.buffer[l.start] = NewLine()
}

func (l *lines) physicalIndex(id uint64) uint64 {
	return id % l.cap
}

// Add добавляет count новых строк в конец логического списка, вытесняя старые при переполнении.
func (l *lines) Add(count uint64) {
	if count == 0 {
		return
	}

	if count >= l.cap {
		l.history += count
		l.start = l.physicalIndex(l.history)
		l.buffer[l.start] = NewLine()

		return
	}

	idx := l.physicalIndex(l.history + count + 1)
	for i := l.physicalIndex(l.history + 1); i != idx; i++ {
		if i == l.cap {
			i = 0
			if i == idx {
				break
			}
		}
		l.buffer[i] = NewLine()
		if i == l.start {
			l.start++
			if l.start == l.cap {
				l.start = 0
			}
		}
	}
	l.history += count
}

func (l *lines) Get(id uint64) Line {
	if id > l.history {
		return nil
	}
	idx := l.physicalIndex(id)
	if idx > l.history || idx < l.start && idx > l.physicalIndex(l.history) {
		return nil
	}

	return l.buffer[idx]
}

func (l *lines) CleanPrefix(id uint64) {
	if id >= l.history {
		l.reset()
		return
	}
	idx := l.physicalIndex(id)
	if idx < l.start && idx > l.physicalIndex(l.history) {
		return
	}
	l.start = idx
}

func (l *lines) CleanPostfix(id uint64) {
	if id >= l.history {
		return
	}

	idx := l.physicalIndex(id)
	if idx < l.start && idx > l.physicalIndex(l.history) {
		l.reset()
		return
	}

	l.history = id
}

func (l *lines) CleanString(id uint64) {
	l.buffer[l.physicalIndex(id)] = NewLine()
}

func (l *lines) LenHistory() uint64 {
	return l.history + 1
}

func (l *lines) GetLastLines(count int) []Line {
	if count <= 0 {
		return []Line{}
	}
	uCount := uint64(count)
	if uCount > l.cap {
		uCount = l.cap
	}
	if uCount > l.history+1 {
		uCount = l.history + 1
	}

	finish := l.physicalIndex(l.history)
	var countCache uint64
	if l.start <= finish {
		countCache = finish - l.start + 1
	} else {
		countCache = l.cap + finish - l.start + 1
	}

	if countCache < uCount {
		uCount = countCache
	}

	result := make([]Line, 0, uCount)

	x := l.physicalIndex(l.history - uCount + 1)
	for ; uCount > 0; uCount-- {
		result = append(result, l.buffer[x])
		x++
		if x == l.cap {
			x = 0
		}
	}

	return result

}

func (l *lines) GetFullLines() []Line {
	finish := l.physicalIndex(l.history)
	var count uint64
	if l.start <= finish {
		count = finish - l.start + 1
	} else {
		count = l.cap + finish - l.start + 1
	}
	result := make([]Line, count)

	x := l.start
	for i := uint64(0); i < count; i++ {
		result[i] = l.buffer[x]
		x++
		if x == l.cap {
			x = 0
		}
	}

	return result
}
