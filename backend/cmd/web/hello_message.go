package web

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// HelloMessage creates a simple HTML component that displays a greeting
func HelloMessage(name string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := w.Write([]byte("<div>Hello, " + name + "!</div>"))
		return err
	})
}
