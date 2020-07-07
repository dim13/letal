package main

import "errors"

type Ban string

var ErrBanValue = errors.New("unknown ban value")

const (
	Hour       Ban = "hour"
	Day        Ban = "day"
	Month      Ban = "month"
	ThreeMonth Ban = "3month"
	SixMonth   Ban = "6month"
	Unlim      Ban = "unlim"
	Remove     Ban = "remove"
	Custom     Ban = "custom"
)

var validBans = []Ban{Hour, Day, ThreeMonth, SixMonth, Unlim, Remove, Custom}

func (b Ban) Usage() string {
	s := "Ban duration:"
	for _, v := range validBans {
		s += " " + string(v)
	}
	return s
}

func (b *Ban) Set(s string) error {
	for _, v := range validBans {
		if Ban(s) == v {
			*b = v
			return nil
		}
	}
	return ErrBanValue
}

func (b Ban) String() string {
	return string(b)
}
