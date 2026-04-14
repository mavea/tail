package line

type Lines interface {
	Get(id uint64) Line
	Add(count uint64)
	LenHistory() uint64
	GetLastLines(count int) []Line
	GetFullLines() []Line
	CleanString(id uint64)
	CleanPostfix(id uint64)
	CleanPrefix(id uint64)
}

type Line interface {
	Set(style string, add string, x int) int
	String(clean string, length int) string
	CleanPrefix(int) int
	CleanPostfix(int) int
}
