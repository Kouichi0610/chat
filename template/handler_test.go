package template

import (
	"testing"
)

func Test_New(t *testing.T) {
	h := New("test")
	if h == nil {
		t.Error()
	}
}
