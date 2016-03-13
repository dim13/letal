package main

import (
	"flag"
	"log"
	"net"
	"sync"
)

const defReason = "Anonymous TOR Coward"

var (
	user, pass, reason, file, target string
	worker, days, port               int
	ban                              = Month
)

func init() {
	flag.StringVar(&user, "user", "", "Username")
	flag.StringVar(&pass, "pass", "", "Password")
	flag.StringVar(&reason, "reason", defReason, "Ban reason")
	flag.StringVar(&file, "file", "", "IP list file")
	flag.StringVar(&target, "target", "linux.org.ru", "Target host")
	flag.IntVar(&port, "port", 443, "Target port")
	flag.IntVar(&worker, "worker", 4, "Concurrency")
	flag.IntVar(&days, "days", 0, "Custom ban duration in days")
	flag.Var(&ban, "ban", banUsage())
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

	var (
		list chan net.IP
		err  error
	)
	if file != "" {
		list, err = List(file)
	} else {
		list, err = Fetch(target, port)
	}
	if err != nil {
		log.Fatal(err)
	}

	c := NewClient()
	if err := c.Login(user, pass); err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	v := BanParams(reason, ban, days, true, false)
	wg := sync.WaitGroup{}
	for i := 0; i < worker; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range list {
				if err := c.BanIP(ip, v); err != nil {
					log.Println(ip, err)
					list <- ip // push back
					return
				}
			}
		}()
	}
	wg.Wait()

	if _, ok := <-list; ok {
		log.Fatal("run out of worker")
	}
}
