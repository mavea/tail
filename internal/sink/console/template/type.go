package template

import (
	"errors"

	configGeneral "tail/internal/config/general"
)

type conf struct {
	value string
}

func NewTemplateType() configGeneral.StringValue {
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
	case "none", "minimal", "full":
		c.value = s
	default:
		return errors.New("invalid screen type value: " + s)
	}

	return nil
}

func (c *conf) Type() string {
	return "string"
}

func (c *conf) Validate() bool {
	switch c.value {
	case "none", "minimal", "full":
		return true
	default:
		return false
	}
}
