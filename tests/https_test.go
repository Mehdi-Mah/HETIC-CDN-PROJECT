package security

import "testing"

func TestUseTLS(t *testing.T) {
	if UseTLS() != false {
		t.Errorf("expected false, got true")
	}
}
