package general

type Code string

type Message string

func (m *Message) String() string {
	return string(*m)
}

type Lang struct {
	MaxLineCount          Message
	MaxCharsPerLine       Message
	MaxBufferLines        Message
	ProcessName           Message
	ProcessIcon           Message
	OutputMode            Message
	OutputTemplate        Message
	Indicator             Message
	Command               Message
	Help                  Message
	Version               Message
	HelpDescription       Message
	HelpBottomDescription Message
	HelpExample           Message
	VersionDescription    Message
	Full                  Message
}
