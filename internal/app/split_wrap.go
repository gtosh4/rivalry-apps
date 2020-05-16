package app

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

func SplitWrap(s string, lim uint) (parts []string) {
	buf := bytes.NewBuffer(make([]byte, 0, int(lim)))

	var lastSep rune
	for len(s) > 0 {
		idx := strings.IndexFunc(s, unicode.IsSpace)
		var word string
		var sep rune
		if idx >= 0 {
			var sepW int
			sep, sepW = utf8.DecodeRuneInString(s[idx:])
			word, s = s[:idx], s[idx+sepW:]
		} else {
			word, s = s, ""
		}

		if buf.Len()+len(word) > int(lim) {
			parts = append(parts, buf.String())
			buf.Reset()
		} else if lastSep > 0 {
			buf.WriteRune(lastSep)
		}
		buf.WriteString(word)
		lastSep = sep
	}
	if buf.Len() > 0 {
		parts = append(parts, buf.String())
	}
	return
}
