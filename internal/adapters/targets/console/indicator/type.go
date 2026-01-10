package indicator

import (
	"errors"
	"tail/internal/config"
)

type conf struct {
	value string
}

func NewType() config.StringValue {
	return &conf{
		value: "none",
	}
}

func (c *conf) String() string {
	return c.value
}

func (c *conf) Set(s string) error {
	switch s {
	case "":
		c.value = "none"
	case "none", "roll":
		c.value = s
	default:
		return errors.New("invalid screen type value: " + s)
	}

	return nil
}
func (c *conf) Type() string {
	return "string"
}
