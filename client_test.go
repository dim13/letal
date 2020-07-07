package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/login.jsp":
			if r.Method != http.MethodGet {
				http.Error(w, "method", http.StatusMethodNotAllowed)
				return
			}
			http.SetCookie(w, &http.Cookie{Name: "CSRF_TOKEN", Value: "token"})
		case "/login_process":
			if r.Method != http.MethodPost {
				http.Error(w, "method", http.StatusMethodNotAllowed)
				return
			}
			if c, err := r.Cookie("CSRF_TOKEN"); err != nil || c.Value != "token" {
				http.Error(w, "csrf token", http.StatusUnauthorized)
				return
			}
			if r.FormValue("nick") != "user" || r.FormValue("passwd") != "pass" {
				http.Error(w, "invalid user", http.StatusUnauthorized)
				return
			}
		case "/banip.jsp":
			if r.Method != http.MethodPost {
				http.Error(w, "method", http.StatusMethodNotAllowed)
				return
			}
			if c, err := r.Cookie("CSRF_TOKEN"); err != nil || c.Value != "token" {
				http.Error(w, "csrf token", http.StatusUnauthorized)
				return
			}
			// TODO
		case "/logout":
			if r.Method != http.MethodPost {
				http.Error(w, "method", http.StatusMethodNotAllowed)
				return
			}
			if c, err := r.Cookie("CSRF_TOKEN"); err != nil || c.Value != "token" {
				http.Error(w, "csrf token", http.StatusUnauthorized)
				return
			}
			http.SetCookie(w, &http.Cookie{Name: "CSRF_TOKEN", MaxAge: -1})
		default:
			http.Error(w, r.URL.String(), http.StatusBadRequest)
		}
	}))
	defer ts.Close()
	c, err := NewClient(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Login("user", "pass"); err != nil {
		t.Fatal("login", err)
	}
	bp := BanParams("none", Month, 0, true, false)
	if err := c.BanIP(net.IP{}, bp); err != nil {
		t.Error(err)
	}
	if err := c.Logout(); err != nil {
		t.Error("logout", err)
	}

}
