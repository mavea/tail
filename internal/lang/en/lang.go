package en

import langGeneral "tail/internal/lang/general"

const (
	Code langGeneral.Code = "en"
)

func NewLang() *langGeneral.Lang {
	return &langGeneral.Lang{
		MaxLineCount:          "maximum number of lines for output to the screen",
		MaxCharsPerLine:       "maximum length of the output line",
		MaxBufferLines:        "maximum size of the stored buffer in number of lines",
		ProcessName:           "name in the program title",
		ProcessIcon:           "icon in the program title",
		OutputMode:            "text output process option (modes: direct - output on new data, thread - periodic redraw)",
		OutputTemplate:        "data output template (none, minimal, full)",
		Indicator:             "specifying an element to display the program's running process (none, roller, bolded-roller, pipe)",
		Command:               "the address of the executable file to be tracked. As an alternative, the pipe can be used to start",
		Help:                  "display help information",
		Version:               "display version information",
		HelpDescription:       "Tail\n Program for monitoring command execution and displaying their output in real-time with various modes and display templates. Supports both direct command execution and data retrieval via pipe. Deletes log upon successful completion, ensuring cleanliness and relevance of displayed information.\n",
		HelpBottomDescription: "For changing the language, when starting, specify the LANG variable in the format ru_RU.UTF-8, or en_US.UTF-8 (the first two characters are important)\n",
		HelpExample: "Examples of use:\n" +
			"  tail -n 10 -o thread -t minimal -- ping google.com\n" +
			"  tail -s 1000 -o direct -t full -- ping google.com\n" +
			"  echo 'Hello World' | tail -o direct -t none\n",
		VersionDescription: "Tail\n Source code:     https://github.com/mavea/tail \n Version:         %s\n  Author:          %s\n  Build time:      %s\n",
		Full:               "full output",
	}
}
