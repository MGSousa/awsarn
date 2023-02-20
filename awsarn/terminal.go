package awsarn

import (
	"fmt"
	"os"
	"strings"

	"github.com/muesli/termenv"
)

type Terminal struct {
	out *termenv.Output
}

func NewTerminal() *Terminal {
	return &Terminal{
		out: termenv.NewOutput(os.Stdout),
	}
}

func (t *Terminal) highlight(arn []string, part int) bool {
	if t.out == nil {
		return false
	}
	var s termenv.Style

	if expr := arn[part]; expr != "" {
		s = t.out.String(expr)
	} else {
		s = t.out.String("<EMPTY>")
	}

	prefix := t.out.String(strings.Join(arn[0:part], DELIMITER))
	suffix := t.out.String(strings.Join(arn[(part+1):], DELIMITER))

	fmt.Println(
		fmt.Sprintf("%s:%s:%s",
			prefix.Foreground(t.green()),
			s.Foreground(t.out.Color("#ffffff")).Background(t.out.Color("#FF0000")).Underline(),
			suffix.Foreground(t.green())))
	return true
}

func (t *Terminal) green() termenv.Color {
	return t.out.Color("#008000")
}
