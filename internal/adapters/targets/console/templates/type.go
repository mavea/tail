package templates

import (
	"errors"
	"tail/internal/config"
)

type t struct {
	value string
}

func NewType() config.StringValue {
	return &t{
		value: "none",
	}
}

func (st *t) String() string {
	return st.value
}

func (st *t) Set(s string) error {
	switch s {
	case "":
		st.value = "none"
	case "none", "minimal", "full":
		st.value = s
	default:
		return errors.New("invalid screen type value: " + s)
	}

	return nil
}

func (st *t) Type() string {
	return "string"
}
