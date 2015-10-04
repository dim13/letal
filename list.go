package main

import (
	"bufio"
	"net"
	"os"
	"strings"
)

type List []net.IP

func NewList(fname string) (List, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var l List
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		ip := net.ParseIP(line)
		if ip != nil {
			l = append(l, ip)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return l, nil
}
