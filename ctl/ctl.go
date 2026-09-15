package ctl

import (
	"cmp"
	"fmt"

	"github.com/alecthomas/kong"
)

// VersionFlag is a flag to print version
type VersionFlag string

// Decode the flag
func (v VersionFlag) Decode(_ *kong.DecodeContext) error { return nil }

// IsBool returns true for the flag
func (v VersionFlag) IsBool() bool { return true }

// BeforeApply is executed before context is applied
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Fprintln(app.Stdout, cmp.Or(vars["version"], string(v)))
	app.Exit(0)
	return nil
}
