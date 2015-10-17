package main

import "fmt"

type Ban int

const (
	Hour Ban = iota
	Day
	Month
	ThreeMonth
	SixMonth
	Unlim
	Remove
	Custom
)

var banNames = map[Ban]string{
	Hour:       "hour",
	Day:        "day",
	Month:      "month",
	ThreeMonth: "3month",
	SixMonth:   "6month",
	Unlim:      "unlim",
	Remove:     "remove",
	Custom:     "custom",
}

func banUsage() string {
	s := "Ban duration:"
	for i := Hour; i < Custom; i++ {
		s += fmt.Sprint(" ", i)
	}
	return s
}

func (b *Ban) Set(s string) error {
	for k, v := range banNames {
		if v == s {
			*b = k
			return nil
		}
	}
	return fmt.Errorf("unknown ban value")
}

func (b Ban) String() string {
	return banNames[b]
}
