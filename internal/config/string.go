package config

type StringValue interface {
	String() string
	Set(string) error
	Type() string
}

type defaultString struct {
	value string
}

func (st *defaultString) String() string {
	return st.value
}
func (st *defaultString) Set(s string) error {
	return nil
}
func (st *defaultString) Type() string {
	return "string"
}
