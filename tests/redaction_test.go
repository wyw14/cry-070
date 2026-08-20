package tests

import (
	"github.com/wyw14/cry052/internal/service"
	"testing"
)

func TestRedactorTextIsCaseInsensitive(t *testing.T) {
	r := service.NewRedactor([]string{"token"})
	if got := r.Text("TOKEN=secret"); got == "TOKEN=secret" {
		t.Fatal("sensitive token leaked")
	}
}
