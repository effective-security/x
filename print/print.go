package print

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"unicode"

	"github.com/effective-security/x/slices"
	"github.com/effective-security/x/values"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

var (
	printRegistry = make(map[reflect.Type]CustomFn)
	registryMutex sync.RWMutex
)

// CustomFn is a custom print function for a specific type.
// It takes precedence over the default print functions.
type CustomFn func(io.Writer, any)

type Printer interface {
	Print(w io.Writer)
}

// RegisterType allows registering a custom print function for a specific type.
func RegisterType(typ any, printFunc CustomFn) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	printRegistry[reflect.TypeOf(typ)] = printFunc
}

// FindRegistered finds a custom print function for a specific value.
func FindRegistered(val any) (CustomFn, bool) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	res, ok := printRegistry[reflect.TypeOf(val)]
	return res, ok
}

// FindRegisteredType finds a custom print function for a specific type.
func FindRegisteredType(typ reflect.Type) (CustomFn, bool) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	res, ok := printRegistry[typ]
	return res, ok
}

// JSON prints value to out
func JSON(w io.Writer, value any) {
	json, _ := json.MarshalIndent(value, "", "\t")
	_, _ = w.Write(json)
	_, _ = w.Write([]byte{'\n'})
}

// Yaml prints value  to out
func Yaml(w io.Writer, value any) {
	y, _ := yaml.Marshal(value)
	_, _ = w.Write(y)
}

// Object prints value to out in format
func Object(w io.Writer, format string, value any) {
	if format == "yaml" {
		Yaml(w, value)
	} else if format == "json" {
		JSON(w, value)
	} else if printFunc, found := FindRegistered(value); found {
		printFunc(w, value)
	} else if pr, ok := value.(Printer); ok {
		pr.Print(w)
	} else {
		// Default to JSON
		JSON(w, value)
	}
}

// Print value
func Print(w io.Writer, value any) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	if printFunc, found := FindRegistered(value); found {
		printFunc(w, value)
		return
	}

	if pr, ok := value.(Printer); ok {
		pr.Print(w)
		return
	}

	switch t := value.(type) {
	case []string:
		Strings(w, t)
	case map[string]string:
		Map(w, []string{"Key", "Value"}, t)
	default:
		JSON(w, value)
	}
}

// Strings prints strings
func Strings(w io.Writer, res []string) {
	for _, r := range res {
		fmt.Fprintln(w, r)
	}
}

// Map prints map
func Map(w io.Writer, header []string, vals map[string]string) {
	table := tablewriter.NewTable(w)
	table.Header(header)

	for _, k := range values.OrderedMapKeys(vals) {
		_ = table.Append([]string{k, slices.StringUpto(vals[k], 80)})
	}

	_ = table.Render()
	fmt.Fprintln(w)
}

// Text prints text with indentation,
// keeping the first line without indentation if noFirstIndent is true
func Text(w io.Writer, doc string, indent string, noFirstIndent bool) {
	if doc == "" {
		return
	}
	parts := strings.Split(doc, "\n")
	for idx, part := range parts {
		if idx > 0 || !noFirstIndent {
			fmt.Fprint(w, indent)
		}
		fmt.Fprintln(w, part)
	}
}

// TextOneLine prints documentation text in one line
func TextOneLine(w io.Writer, doc string) {
	if doc == "" {
		return
	}
	parts := strings.Split(doc, "\n")
	lines := 0
	prevPartDot := false
	for _, part := range parts {
		part = strings.TrimSpace(part)
		size := len(part)
		if size > 0 {
			if lines > 0 {
				if !prevPartDot && unicode.IsUpper(rune(part[0])) {
					fmt.Fprint(w, ".")
				}
				fmt.Fprint(w, " ")
			}
			fmt.Fprint(w, part)
			prevPartDot = part[size-1] == '.'
			lines++
		}
	}
	if lines > 0 && !prevPartDot {
		fmt.Fprint(w, ".")
	}
}
