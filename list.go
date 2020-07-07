package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const ExitList = `https://check.torproject.org/torbulkexitlist`

func Fetch(from, target string) (chan net.IP, error) {
	t, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	ip, err := net.LookupIP(t.Hostname())
	if err != nil {
		return nil, err
	}
	q, err := url.Parse(from)
	if err != nil {
		return nil, err
	}
	q.Query().Add("ip", ip[0].String())
	q.Query().Add("port", t.Port())
	resp, err := http.Get(q.String())
	if err != nil {
		return nil, err
	}
	if err := checkStatus(resp); err != nil {
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
