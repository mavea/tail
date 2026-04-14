package general

type StringValue interface {
	String() string
	Set(string) error
	Type() string
	Validate() bool
}

type DefaultString struct {
	value string
}

func (st *DefaultString) String() string {
	return st.value
}
func (st *DefaultString) Set(s string) error {
	st.value = s
	return nil
}
func (st *DefaultString) Type() string {
	return "string"
}

func (st *DefaultString) Validate() bool {
	return true
}
