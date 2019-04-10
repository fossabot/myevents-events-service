package component

import "testing"

func TestComponent(t * testing.T) {
	if testing.Short() {
		t.Skip("Component test")
	}
}
