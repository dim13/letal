package main

//go:generate curl -o torlist "https://check.torproject.org/cgi-bin/TorBulkExitList.py?ip=178.248.233.6&port=443"

import (
	"flag"
	"log"
)

const defReason = "Anonymous TOR Coward"

var user, pass, reason, file string

func init() {
	flag.StringVar(&user, "user", "", "Username")
	flag.StringVar(&pass, "pass", "", "Password")
	flag.StringVar(&reason, "reason", defReason, "Ban reason")
	flag.StringVar(&file, "file", "torlist", "IP list file")
}

func main() {
	flag.Parse()
	if days != 0 {
		ban = Custom
	}
	if user == "" || pass == "" {
		flag.PrintDefaults()
		return
	}

	list, err := NewList(file)
	if err != nil {
		log.Fatal(err)
	}

	c := NewClient()
	c.Login(user, pass)
	defer c.Logout()

	v := BanParams(reason, ban, days, true, false)
	for n, ip := range list {
		log.Println(n, "of", len(list), "ban", ip)
		err = c.BanIP(ip, v)
		if err != nil {
			log.Fatal("line", n, err)
		}
	}
}
