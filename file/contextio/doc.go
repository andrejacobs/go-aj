// Package contextio provides [io.Writer] and [io.Reader] that stop accepting/providing
// data when an attached context is canceled.
//
// AJ: I want to minimize 3rd party dependencies on my own go-aj packages
// This was taken from https://github.com/dolmen-go/contextio as is.
package contextio
