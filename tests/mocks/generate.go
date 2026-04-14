package mocks

//go:generate go run github.com/matryer/moq@v0.5.3 -pkg mocks -out scanner_moq.go ../../internal/source/general Scanner
//go:generate go run github.com/matryer/moq@v0.5.3 -pkg mocks -out render_moq.go ../../internal/formatter Render
//go:generate go run github.com/matryer/moq@v0.5.3 -pkg mocks -out formatter_out_moq.go ../../internal/formatter Out
//go:generate go run github.com/matryer/moq@v0.5.3 -pkg mocks -out template_moq.go ../../internal/sink/console/template/general Template
//go:generate go run github.com/matryer/moq@v0.5.3 -pkg mocks -out template_indicator_moq.go ../../internal/sink/console/template/general Indicator
//go:generate go run github.com/matryer/moq@v0.5.3 -pkg mocks -out template_window_moq.go ../../internal/sink/console/template/general Window
//go:generate go run github.com/matryer/moq@v0.5.3 -pkg mocks -out parser_mainlist_moq.go ../../internal/formatter/buffer/parser MainList
//go:generate go run github.com/matryer/moq@v0.5.3 -pkg mocks -out parser_windowpalettes_moq.go ../../internal/formatter/buffer/parser WindowAndPalettes
//go:generate go run github.com/matryer/moq@v0.5.3 -pkg mocks -out parser_decprivatemodes_moq.go ../../internal/formatter/buffer/parser DECPrivateModes
//go:generate go run github.com/matryer/moq@v0.5.3 -pkg mocks -out parser_secondaryda_moq.go ../../internal/formatter/buffer/parser SecondaryDA
//go:generate go run github.com/matryer/moq@v0.5.3 -pkg mocks -out parser_otherlist_moq.go ../../internal/formatter/buffer/parser OtherList
