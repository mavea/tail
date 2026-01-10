package domain

type Scanner interface {
	Get() (string, bool)
	GetErr() (error, bool)
	GetStatus() (int, error)
}
