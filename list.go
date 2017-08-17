package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

const exitList = `https://check.torproject.org/cgi-bin/TorBulkExitList.py?ip=%s&port=%d`

func Fetch(host string, port int) (chan net.IP, error) {
	ip, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(exitList, ip[0], port)
	log.Printf("fetch list for %s:%d\n", ip[0], port)
	resp, err := http.Get(query)
	if err != nil {
		return nil, err
	}
	if err := check(resp); err != nil {
		return nil, err
	}
	return list(resp.Body)
}

func List(fname string) (chan net.IP, error) {
	fd, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	return list(fd)
}

func list(r io.ReadCloser) (chan net.IP, error) {
	c := make(chan net.IP)
	go func() {
		scanner := bufio.NewScanner(r)
		defer r.Close()
		defer close(c)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") {
				continue
			}
			c <- net.ParseIP(line)
		}
		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	}()
	return c, nil
}
