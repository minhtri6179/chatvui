package web

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// HelloForm returns a templ component with a simple input form
func HelloForm() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := w.Write([]byte(`
		<form action="/hello" method="POST">
			<label for="name">Your name:</label>
			<input type="text" id="name" name="name" required>
			<button type="submit">Submit</button>
		</form>
		`))
		return err
	})
}
