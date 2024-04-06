package main

import (
	"testing"
)

func TestTimeToHuman(t *testing.T) {
	in := Week
	if x := TimeToHuman(&in); x != "1W" {
		t.Errorf("expected %s, got %s", "1W", x)
	}

	in = Week + Hour
	if x := TimeToHuman(&in); x != "1W1H" {
		t.Errorf("expected %s, got %s", "1W!H", x)
	}
}
