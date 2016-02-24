package rss

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/transform"

	"github.com/axgle/mahonia"
)

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	discarderReader := validUTF8Discarder{}
	switch {
	case isCharsetUTF8(charset):
		fmt.Println("FSAIJFAISO")
		return transform.NewReader(input, discarderReader), nil
	default:
		fmt.Println("FSAIJFAISO")
		if decoder := mahonia.NewDecoder(charset); decoder != nil {
			fmt.Println("FSAIJFAISOasdasdsa")
			return transform.NewReader(decoder.NewReader(input), discarderReader), nil
		}
	}

	return nil, errors.New("CharsetReader: unexpected charset: " + charset)
}

func isCharset(charset string, names []string) bool {
	charset = strings.ToLower(charset)
	for _, n := range names {
		if charset == strings.ToLower(n) {
			return true
		}
	}
	return false
}

func isCharsetUTF8(charset string) bool {
	names := []string{
		"UTF-8",
		// Default
		"",
	}
	return isCharset(charset, names)
}

type validUTF8Discarder struct {
}

func (r validUTF8Discarder) Transform(dst []byte, src []byte, atEOF bool) (nDst int, nSrc int, err error) {
	buf := src
	i := 0

	const undefinedCharacter = 0x3F // ?
	var err1 error

	for len(buf) > 0 {
		r, size := utf8.DecodeRune(buf)
		buf = buf[size:]

		if r == utf8.RuneError && size == 1 {
			dst[i] = undefinedCharacter
			i++
		} else if !isInCharacterRange(r) {
			for x := 0; x < size; x++ {
				dst[i] = undefinedCharacter
				i++
			}
		} else {
			for x := 0; x < size; x++ {
				dst[i] = src[i]
				i++
			}
		}
	}

	if i > 0 && dst[i-1] == 0 {
		err1 = transform.ErrShortDst
	}

	return i, i, err1
}

// Reset resets the state and allows a Transformer to be reused.
func (r validUTF8Discarder) Reset() {

}

// Decide whether the given rune is in the XML Character Range, per
// the Char production of http://www.xml.com/axml/testaxml.htm,
// Section 2.2 Characters.
func isInCharacterRange(r rune) (inrange bool) {
	return r == 0x09 ||
		r == 0x0A ||
		r == 0x0D ||
		r >= 0x20 && r <= 0xDF77 ||
		r >= 0xE000 && r <= 0xFFFD ||
		r >= 0x10000 && r <= 0x10FFFF
}

var hasUTF8 *regexp.Regexp

// DiscardInvalidUTF8IfUTF8 checks if input specifies itself as UTF8,
// and then runs a check to discard XML-invalid characters (because go xml parser throws up if present)
func DiscardInvalidUTF8IfUTF8(input []byte) []byte {
	if hasUTF8 == nil {
		hasUTF8 = regexp.MustCompile(`(?i)^.*<\?xml.*encoding=.*utf.?8`)
	}

	firstChunk := string(input[0:1024])

	if hasUTF8.MatchString(firstChunk) {
		fmt.Println("LOL")
		reader := bytes.NewReader(input)
		discarderReader := validUTF8Discarder{}

		transformer := transform.NewReader(reader, discarderReader)

		value, err := ioutil.ReadAll(transformer)
		if err != nil {
			return input
		}
		return value
	}
	return input
}
