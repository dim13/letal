package main

//go:generate curl -o torlist "https://check.torproject.org/cgi-bin/TorBulkExitList.py?ip=178.248.233.6&port=443"

import (
	"flag"
	"log"
	"sync"
)

const defReason = "Anonymous TOR Coward"

var user, pass, reason, file string
var concurrancy int

func init() {
	flag.StringVar(&user, "user", "", "Username")
	flag.StringVar(&pass, "pass", "", "Password")
	flag.StringVar(&reason, "reason", defReason, "Ban reason")
	flag.StringVar(&file, "file", "torlist", "IP list file")
	flag.IntVar(&concurrancy, "concurrancy", 10, "Concurrancy")
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

	list, err := List(file)
	if err != nil {
		log.Fatal(err)
	}

	c := NewClient()
	c.Login(user, pass)
	defer c.Logout()

	v := BanParams(reason, ban, days, true, false)
	wg := sync.WaitGroup{}
	for i := 0; i < concurrancy; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range list {
				err := c.BanIP(ip, v)
				if err != nil {
					log.Println(ip, err)
					list <- ip // push back
					return
				}
			}
		}()
	}
	wg.Wait()
}
