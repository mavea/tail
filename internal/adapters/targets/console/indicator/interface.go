package indicator

type Indicator interface {
	Clean() string
	Get() string
}
