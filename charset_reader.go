package rss

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/text/transform"

	"github.com/axgle/mahonia"
)

// CharsetReader is a lenient charset reader good for web inputs
func CharsetReader(charset string, input io.Reader) (io.Reader, error) {
	discarderReader := validUTF8Discarder{}

	switch {
	case isCharsetUTF8(charset):
		return transform.NewReader(input, discarderReader), nil
	// Go works really weird with UTF-16
	case isCharsetUTF16(charset):
		return transform.NewReader(input, discarderReader), nil
	default:
		if decoder := mahonia.NewDecoder(charset); decoder != nil {
			return transform.NewReader(decoder.NewReader(input), discarderReader), nil
		}
	}

	return nil, errors.New("CharsetReader: unexpected charset: " + charset)
}

func compareStringToStrings(charset string, names []string) bool {
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
		"UTF8",
		// Default
		"",
	}
	return compareStringToStrings(charset, names)
}

func isCharsetUTF16(charset string) bool {
	names := []string{
		"UTF-16",
		"UTF16",
	}
	return compareStringToStrings(charset, names)
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
var hasEncodingAtAll *regexp.Regexp

// DiscardInvalidUTF8IfUTF8 checks if input specifies itself as UTF8,
// and then runs a check to discard XML-invalid characters (because go xml parser throws up if present)
func DiscardInvalidUTF8IfUTF8(input []byte, responseHeaders http.Header) []byte {
	if hasUTF8 == nil {
		hasUTF8 = regexp.MustCompile(`(?i)^.*<\?xml.*encoding=.*utf.?8`)
	}

	if hasEncodingAtAll == nil {
		hasEncodingAtAll = regexp.MustCompile(`(?i)^.*<\?xml.*encoding=.*`)
	}

	var firstChunk string
	if len(input) > 1024 {
		firstChunk = string(input[0:1024])
	} else {
		firstChunk = string(input)
	}

	charsetFromHeaders := getCharsetFromHeaders(responseHeaders)

	// UTF-16 handling in go seems a bit weird
	// the Go XML parser wouldn't even attempt to decode unless we decode 16 first
	// Hopefully a smarter person than me will see this, rage so hard at the stupidity of this fix
	// and then open an issue letting me know how it's done
	if isCharsetUTF16(charsetFromHeaders) {
		utf16DecodedInput, err := decodeUTF16(input)
		if err == nil {
			return utf16DecodedInput
		}
	}

	if hasUTF8.MatchString(firstChunk) || !hasEncodingAtAll.MatchString(firstChunk) {
		// Some feeds respond with a <?xml encoding=utf8 even though their server
		// indicates another charset. Here we act to fix that, by detecting if a
		// header indicates something else. An example found in the wild:
		//     Content-Type: application/rss+xml; Charset=ISO-8859-9
		// this block would then convert ISO-8859-9 to UTF8 and then run the discarder on the input afterwards
		if charsetFromHeaders != "" && charsetFromHeaders != "utf-8" && charsetFromHeaders != "utf8" {
			dec := mahonia.NewDecoder(charsetFromHeaders)
			if dec != nil {
				convertedToUtf8 := dec.ConvertString(string(input))
				input = []byte(convertedToUtf8)
			}
		}

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

func getCharsetFromHeaders(responseHeaders http.Header) string {
	if responseHeaders == nil {
		return ""
	}

	contentType := responseHeaders.Get("Content-Type")
	pieces := strings.Split(contentType, ";")

	for _, piece := range pieces {
		charsetPieces := strings.Split(strings.ToLower(piece), "charset=")
		if len(charsetPieces) == 2 {
			return strings.ToLower(charsetPieces[1])
		}
	}

	return ""
}

func decodeUTF16(b []byte) ([]byte, error) {
	if len(b)%2 != 0 {
		return nil, errors.New("Must have even number of chars")
	}

	u16s := make([]uint16, 1)
	ret := &bytes.Buffer{}
	b8buf := make([]byte, 4)
	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.Bytes(), nil
}
