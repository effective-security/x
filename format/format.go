package format

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func YesNo(val bool) string {
	if val {
		return "yes"
	}
	return "no"
}

func Enabled(val bool) string {
	if val {
		return "enabled"
	}
	return "disabled"
}

func Number[T ~int | ~int32 | ~uint | ~uint32 | ~int64 | ~uint64](val T) string {
	return fmt.Sprintf("%d", val)
}

func Float[T ~float32 | ~float64](val T) string {
	return fmt.Sprintf("%0.2f", val)
}

// StringMax returns the string value with a maximum length of max.
func StringMax(limit int, val string) string {
	if len(val) > limit {
		return val[:limit] + "..."
	}
	return val
}

func Strings(val []string) string {
	if len(val) == 0 {
		return ""
	}
	var buf strings.Builder
	for i, s := range val {
		if s == "" {
			continue
		}
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(s)
	}
	return buf.String()
}

// StringsMax returns the string value with a maximum length of max.
func StringsMax(limit int, val []string) string {
	if len(val) == 0 {
		return ""
	}
	var buf strings.Builder
	for i, s := range val {
		if s == "" {
			continue
		}
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(s)
		if buf.Len() >= limit {
			// add ellipsis
			if i < len(val)-1 {
				buf.WriteString("...")
			}
			return buf.String()
		}
	}
	return buf.String()
}

func StringsAndMore(first int, val []string) string {
	count := len(val)
	if count <= first {
		return Strings(val)
	}
	var buf strings.Builder
	for i, s := range val {
		if s == "" {
			continue
		}
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(s)
		if i >= first-1 {
			break
		}
	}
	// add ellipsis
	if count > first {
		fmt.Fprintf(&buf, ", %d more...", count-first)
	}
	return buf.String()
}

// DisplayName returns display name of a field to preserve common acronyms
func DisplayName(name string) string {
	s := Split(name)
	if len(s) == 0 {
		return name
	}
	return strings.Join(s, " ")
}

// Split splits the camelcase word and returns a list of words. It also
// supports digits. Both lower camel case and upper camel case are supported.
// For more info please check: http://en.wikipedia.org/wiki/CamelCase
//
// Examples
//
//	"" =>                     [""]
//	"lowercase" =>            ["lowercase"]
//	"Class" =>                ["Class"]
//	"MyClass" =>              ["My", "Class"]
//	"MyC" =>                  ["My", "C"]
//	"HTML" =>                 ["HTML"]
//	"PDFLoader" =>            ["PDF", "Loader"]
//	"AString" =>              ["A", "String"]
//	"SimpleXMLParser" =>      ["Simple", "XML", "Parser"]
//	"vimRPCPlugin" =>         ["vim", "RPC", "Plugin"]
//	"GL11Version" =>          ["GL", "11", "Version"]
//	"99Bottles" =>            ["99", "Bottles"]
//	"May5" =>                 ["May", "5"]
//	"BFG9000" =>              ["BFG", "9000"]
//	"BöseÜberraschung" =>     ["Böse", "Überraschung"]
//	"Two  spaces" =>          ["Two", "  ", "spaces"]
//	"BadUTF8\xe2\xe2\xa1" =>  ["BadUTF8\xe2\xe2\xa1"]
//
// Splitting rules
//
//  1. If string is not valid UTF-8, return it without splitting as
//     single item array.
//  2. Assign all unicode characters into one of 4 sets: lower case
//     letters, upper case letters, numbers, and all other characters.
//  3. Iterate through characters of string, introducing splits
//     between adjacent characters that belong to different sets.
//  4. Iterate through array of split strings, and if a given string
//     is upper case:
//     if subsequent string is lower case:
//     move last character of upper case string to beginning of
//     lower case string
func Split(src string) (entries []string) {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return []string{src}
	}
	entries = []string{}
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass || (lastClass == 2 && class == 3) {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}
	return
}

func TextWithIndent(doc string, indent string, noFirstIndent bool) string {
	if doc == "" {
		return ""
	}
	var buf strings.Builder
	parts := strings.Split(doc, "\n")
	for idx, part := range parts {
		if idx > 0 || !noFirstIndent {
			buf.WriteString(indent)
		}
		buf.WriteString(part)
		buf.WriteString("\n")
	}
	return buf.String()
}

func TextOneLine(doc string) string {
	if doc == "" {
		return ""
	}

	var buf strings.Builder
	parts := strings.Split(doc, "\n")
	lines := 0
	prevPartDot := false
	for _, part := range parts {
		part = strings.TrimSpace(part)
		size := len(part)
		if size > 0 {
			if lines > 0 {
				if !prevPartDot && unicode.IsUpper(rune(part[0])) {
					buf.WriteString(".")
				}
				buf.WriteString(" ")
			}
			buf.WriteString(part)
			prevPartDot = part[size-1] == '.'
			lines++
		}
	}
	if lines > 0 && !prevPartDot {
		buf.WriteString(".")
	}
	return buf.String()
}
