package processor

import "context"

type in interface {
	Get() (string, bool)
	GetErr() (error, bool)
	GetStatus() (int, error)
}

type out interface {
	Set(data string) error
	SetErr(err error) error
	SetStatus(status int) error
}

type Service interface {
	Run(ctx context.Context) error
}
