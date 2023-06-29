# Golang Translate Datetime Formatting
This library allows you to use regular "formatting" rules, sometimes they are read from files etc and can be cumbersome to translate. This library does the job for you.

As an example, it will translate something like "YYYY-MM=dd HH:mm:ss" to "2006-01-02 15:04:05" that can be added right in to your time.Parse arguments.

## Sample code:
```golang
package main

import (
	"strings"
	"time"

	format "github.com/hektorinho/goDatetimeFormatTranslator"
)

const (
	myFormat = "YYYY-MM-dd HH:mm:ss"
)

func main() {
	dateTime := "2023-06-28 14:35:26"
	dec := format.NewDecoder(strings.NewReader(myFormat))
	dec.SetFormat(myFormat)
	golangFormat := format.TranslateTokens(dec.ReadTokens())

	myTime := time.Parse(golangFormat, dateTime)
}
```