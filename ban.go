package main

import (
	"errors"
	"flag"
	"fmt"
)

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

var (
	ban  = Month
	days int
)

func init() {
	flag.Var(&ban, "ban", banUsage())
	flag.IntVar(&days, "days", 0, "Custom ban in days")
}

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
	s := "Ban:"
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
	return errors.New("unknown ban value")
}

func (b Ban) String() string {
	return banNames[b]
}
