package trace

import (
	"fmt"
	"io"
)

type Tracer interface {
	Trace(...interface{})
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

func Empty() Tracer {
	return &empty{}
}

type tracer struct {
	out io.Writer // 出力先
}

func (t *tracer) Trace(a ...interface{}) {
	s := fmt.Sprint(a...)
	t.out.Write([]byte(s))
	t.out.Write([]byte("\n"))
}

type empty struct {
}

func (t *empty) Trace(a ...interface{}) {
}
