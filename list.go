package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

func List(fname string) (chan net.IP, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	c := make(chan net.IP)
	go func() {
		scanner := bufio.NewScanner(file)
		defer file.Close()
		defer close(c)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") {
				continue
			}
			c <- net.ParseIP(line)
		}
		log.Println("list done")
	}()
	return c, nil
}
