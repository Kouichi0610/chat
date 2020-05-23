package trace

import (
	"bytes"
	"testing"
)

func Test_Tracer(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Errorf("failed new")
		return
	}
	tracer.Trace("Hello")
	if buf.String() != "Hello¥n" {
		t.Errorf("failed [%s]¥n", buf.String())
	}
}

func Test_Empty(t *testing.T) {
	tracer := Empty()
	if tracer == nil {
		t.Error()
	}
	tracer.Trace("Hello")
}
