package indicator

type Indicator interface {
	Clean(bool) string
	Get() string
}
