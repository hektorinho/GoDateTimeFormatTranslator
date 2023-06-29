package format

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type Decoder struct {
	r   *bufio.Reader
	tok Token

	CheckValidFunc func(byte, byte) bool
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: bufio.NewReader(r), CheckValidFunc: CheckNextPartOfToken}
}

func (d *Decoder) SetCheckValidFunc(fn func(byte, byte) bool) {
	d.CheckValidFunc = fn
}

type Token struct {
	Type   string
	Length int

	dateType    DateType
	currentByte byte
	nextByte    byte
	prevByte    byte
	idx         int

	lastByte       bool
	validDateToken bool
}

type DateType uint8

const (
	Year DateType = iota
	Month
	Day
	Hour
	Minute
	Second
	Fractional

	DayNightDivision
	Timezone
	TimezoneLocation
)

func (d *Decoder) Token() (Token, error) {
	var tok Token
	for {
		// advance reader and read current byte
		b, err := d.r.ReadByte()
		if err != io.EOF && err != nil {
			return Token{}, err
		}
		// set previous byte before setting current byte
		tok.prevByte = tok.currentByte
		tok.currentByte = b
		// Peek in to the next byte, first check if we're on the last byte
		if !tok.lastByte {
			p, err := d.r.Peek(1)
			if err == io.EOF {
				tok.lastByte = true
			}
			if !tok.lastByte {
				tok.nextByte = p[0]
			}
		}
		// Continue string if next byte is part of Token
		if d.CheckValidFunc(tok.currentByte, tok.nextByte) {
			tok.Type += string(tok.currentByte)
		}
		// Break token if next byte is from a different Token
		if !d.CheckValidFunc(tok.currentByte, tok.nextByte) {
			tok.Type += string(tok.currentByte)
			tok.Length = len(tok.Type)
			d.tok = tok
			tok.Type = strings.TrimSuffix(tok.Type, string([]byte{0}))
			return tok, nil
		}
		// If we reached end of format, then break last Token
		if err == io.EOF {
			tok.Type += string(tok.currentByte)
			tok.Length = len(tok.Type)
			d.tok = tok
			tok.Type = strings.TrimSuffix(tok.Type, string([]byte{0}))
			return tok, io.EOF
		}
		tok.idx++
	}
}

func (d *Decoder) ReadTokens() ([]Token, error) {
	var tokens []Token
	for {
		token, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

var (
	ErrBadFormat = errors.New("format: bad format, not matching dictionary")
)

func (d *Decoder) Translate(dict map[string]string) (string, error) {
	tokens, err := d.ReadTokens()
	if err != nil {
		return "", ErrBadFormat
	}
	var golangFormat string
	for _, token := range tokens {
		if val, ok := dict[token.Type]; ok {
			golangFormat += val
		} else {
			golangFormat += token.Type
		}
	}
	return strings.TrimSpace(golangFormat), nil
}

func CheckNextPartOfToken(current, next byte) bool {
	if current == next {
		return true
	}
	// add exceptions here
	// Y can be followed by lowercase y
	if (current == byte(89)) && (next == byte(121)) {
		return true
	}
	// M followed by mm will be Jan
	if (current == byte(77)) && (next == byte(109)) {
		return true
	}
	// D can be followed by lowercase d
	if (current == byte(68)) && (next == byte(100)) {
		return true
	}
	// H can be followed by lowercase h
	if (current == byte(72)) && (next == byte(104)) {
		return true
	}
	// S can be followed by lowercase s
	if (current == byte(83)) && (next == byte(115)) {
		return true
	}
	// S can be followed by lowercase s
	if (current == byte(70)) && (next == byte(102)) {
		return true
	}
	// Z can be followed by lowercase z
	if (current == byte(90)) && (next == byte(122)) {
		return true
	}
	// D can be followed by lowercase a which can be followed by lowercase y, within Day
	if (current == byte(68)) && (next == byte(97)) || (current == byte(97)) && (next == byte(121)) {
		return true
	}

	return false
}

var (
	ValidTokens = map[byte]DateType{
		89:  Year,
		121: Year,
		77:  Month,
		68:  Day,
		100: Day,
		72:  Hour,
		104: Hour,
		109: Minute,
		115: Second,
		70:  Fractional,
		102: Fractional,
		65:  DayNightDivision,
		97:  DayNightDivision,
		90:  Timezone,
		122: Timezone,
		79:  TimezoneLocation,
		111: TimezoneLocation,
	}
	StandardTokens = map[string]string{
		"yyyy":      "2006",
		"yy":        "06",
		"YYYY":      "2006",
		"YY":        "06",
		"Yyyy":      "2006",
		"Yy":        "06",
		"M":         "1",
		"MM":        "01",
		"MMM":       "Jan",
		"Mmm":       "Jan",
		"mmm":       "Jan",
		"MMMM":      "January",
		"Mmmm":      "January",
		"mmmm":      "January",
		"D":         "2",
		"DD":        "02",
		"Dd":        "02",
		"d":         "2",
		"dd":        "02",
		"h":         "3",
		"H":         "3",
		"hh":        "03",
		"HH":        "15",
		"Hh":        "15",
		"m":         "4",
		"mm":        "04",
		"s":         "5",
		"ss":        "05",
		"S":         "5",
		"SS":        "05",
		"Ss":        "05",
		"f":         "9",
		"ff":        "99",
		"fff":       "999",
		"ffff":      "9999",
		"fffff":     "99999",
		"ffffff":    "999999",
		"fffffff":   "9999999",
		"ffffffff":  "99999999",
		"fffffffff": "999999999",
		"F":         "9",
		"FF":        "99",
		"FFF":       "999",
		"FFFF":      "9999",
		"FFFFF":     "99999",
		"FFFFFF":    "999999",
		"FFFFFFF":   "9999999",
		"FFFFFFFF":  "99999999",
		"FFFFFFFFF": "999999999",
		"Ff":        "99",
		"Fff":       "999",
		"Ffff":      "9999",
		"Fffff":     "99999",
		"Ffffff":    "999999",
		"Fffffff":   "9999999",
		"Ffffffff":  "99999999",
		"Fffffffff": "999999999",
		"A":         "PM",
		"a":         "pm",
		"z":         "-07",
		"zz":        "-0700",
		"zzz":       "-7:00",
		"Z":         "-07",
		"ZZ":        "-0700",
		"ZZZ":       "-7:00",
		"Zz":        "-0700",
		"Zzz":       "-7:00",
		"O":         "MST",
		"o":         "mst",
		"Day":       "Monday",
	}
)
