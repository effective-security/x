package print

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/effective-security/x/slices"
	"github.com/effective-security/x/values"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

var (
	printRegistry = make(map[reflect.Type]func(io.Writer, any))
	registryMutex sync.RWMutex
)

type Printer interface {
	Print(w io.Writer)
}

// RegisterType allows registering a custom print function for a specific type.
func RegisterType(typ any, printFunc func(io.Writer, any)) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	printRegistry[reflect.TypeOf(typ)] = printFunc
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
	} else if printFunc, found := printRegistry[reflect.TypeOf(value)]; found {
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

	if pr, ok := value.(Printer); ok {
		pr.Print(w)
		return
	}

	if printFunc, found := printRegistry[reflect.TypeOf(value)]; found {
		printFunc(w, value)
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
	table := tablewriter.NewWriter(w)
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetHeader(header)
	table.SetHeaderLine(true)

	for _, k := range values.OrderedMapKeys(vals) {
		table.Append([]string{k, slices.StringUpto(vals[k], 80)})
	}

	table.Render()
	fmt.Fprintln(w)
}
