package main

import (
	"flag"
	"log"
	"net"
	"os"
	"sync"
)

const defReason = "Anonymous TOR Coward"

func main() {
	var (
		user, pass, reason, file, target string
		worker, days                     int
		ban                              = Month
	)

	flag.StringVar(&user, "user", os.Getenv("LORUSER"), "Username")
	flag.StringVar(&pass, "pass", os.Getenv("LORPASS"), "Password")
	flag.StringVar(&reason, "reason", defReason, "Ban reason")
	flag.StringVar(&file, "file", "", "IP list file")
	flag.StringVar(&target, "target", "https://linux.org.ru", "Target host")
	flag.IntVar(&worker, "worker", 2, "Concurrency")
	flag.IntVar(&days, "days", 0, "Custom ban duration in days")
	flag.Var(&ban, "ban", ban.Usage())

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
		list, err = Fetch(ExitList, target)
	}
	if err != nil {
		log.Fatal(err)
	}

	c, err := NewClient(target)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Login(user, pass); err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	bp := BanParams(reason, ban, days, true, false)
	wg := sync.WaitGroup{}
	for i := 0; i < worker; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range list {
				log.Println(ip, ban)
				if err := c.BanIP(ip, bp); err != nil {
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
