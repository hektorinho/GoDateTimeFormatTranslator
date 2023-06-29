package format

import (
	"fmt"
	"strings"
	"testing"
)

func TestDecodeFormat(t *testing.T) {
	testFormats := []struct {
		external         string
		expectedInternal string
		actualOutput     string
	}{
		{
			external:         "YYYY-MM-dd HH:mm:ss",
			expectedInternal: "2006-01-02 15:04:05",
		},
		{
			external:         "YYYY-MM-DD",
			expectedInternal: "2006-01-02",
		},
		{
			external:         "YYYY/MM/DD",
			expectedInternal: "2006/01/02",
		},
		{
			external:         "DD-MM-YYYY",
			expectedInternal: "02-01-2006",
		},
		{
			external:         "DD/MM/YYYY",
			expectedInternal: "02/01/2006",
		},
		{
			external:         "YYYY-MM-DD hh:mm:ss A",
			expectedInternal: "2006-01-02 03:04:05 PM",
		},
		{
			external:         "MMM DD,YYYY",
			expectedInternal: "Jan 02,2006",
		},
		{
			external:         "MMMM DD,YYYY",
			expectedInternal: "January 02,2006",
		},
		{
			external:         "DD/MM/YY",
			expectedInternal: "02/01/06",
		},
		{
			external:         "HH:mm:ss.ffffff",
			expectedInternal: "15:04:05.999999",
		},
		{
			external:         "YY-Mmm-dd",
			expectedInternal: "06-Jan-02",
		},
		{
			external:         "Day",
			expectedInternal: "Monday",
		},
		{
			external:         "Yyyy-Mmm-Dd",
			expectedInternal: "2006-Jan-02",
		},
	}

	for _, format := range testFormats {
		dec := NewDecoder(strings.NewReader(format.external))
		outputFormat, err := dec.Translate(StandardTokens)
		if err != nil {
			t.Errorf("format: failed to translate tokens > %s", err)
		}
		format.actualOutput = outputFormat

		if format.actualOutput != format.expectedInternal {
			fmt.Println([]byte(format.expectedInternal))
			fmt.Println([]byte(format.actualOutput))
			t.Errorf("format: failed | expected=%s | got=%s", format.expectedInternal, format.actualOutput)
		}
	}
}
