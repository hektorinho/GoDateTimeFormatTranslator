# Golang Translate Datetime Formatting
This library allows you to use regular "formatting" rules, sometimes they are read from files etc and can be cumbersome to translate. This library does the job for you.

As an example, it will translate something like "YYYY-MM=dd HH:mm:ss" to "2006-01-02 15:04:05" that can be added right in to your time.Parse arguments.

## You can create you own Validation Function as well as your own Token Dictionary
The default validation function CheckNextPartOfToken(current, next byte) bool checks if next token is a different kind than previous. It makes a few exceptions for like M in to m to enable Mmm for Jan, Feb etc.

### Code to change CheckValidFunc:
```golang
...
dec.SetCheckValidFunc(CheckNextPartOfToken)
...
```

You can also update the Standard Token dictionary that is being used as input to Translate. StandardTokens is a map: map[string]string small subset can be seen below:

```golang
...
StandardTokens = map[string]string{
        ...
		"M":         "1",
		"MM":        "01",
		"MMM":       "Jan",
		"Mmm":       "Jan",
		"mmm":       "Jan",
		"MMMM":      "January",
		"Mmmm":      "January",
		"mmmm":      "January",
        ...
}
```

## Sample code:
```golang
package main

import (
    "fmt"
	"log"
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
	golangFormat, err := dec.Translate(format.StandardTokens)
	if err != nil {
		log.Fatalf("format: failed to translate tokens >> %s", err)
	}

	myTime := time.Parse(golangFormat, dateTime)
    fmt.Println(myTime)
	// 2023-06-28 14:35:26 +0000 UTC
}
```