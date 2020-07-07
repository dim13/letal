package main

import "testing"

func TestBanFlag(t *testing.T) {
	b := Month
	if b.String() != "month" {
		t.Error("string", b)
	}
	b.Set("day")
	if b.String() != "day" {
		t.Error("set", b)
	}
	if err := b.Set("none"); err != ErrBanValue {
		t.Error("err", err)
	}
	if b.Usage() != "Ban duration: hour day 3month 6month unlim remove custom" {
		t.Error("usage", b.Usage())
	}
}
