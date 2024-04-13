package main

import "testing"

func TestSerialToHuman(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"1282630063", "Tue, 24 Aug 2010 06:07:43 UTC"},
	}
	for i, ts := range tests {
		if x := SerialToHuman([]byte(ts.in)); x != "  "+ts.out {
			t.Errorf("test %d, expected %s, got %s", i, ts.out, x)
		}
	}
}
