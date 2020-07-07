package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	*http.Client
	base *url.URL
}

func (c Client) csrf() string {
	for _, cookie := range c.Jar.Cookies(c.base) {
		if cookie.Name == "CSRF_TOKEN" {
			if i := strings.Index(cookie.Value, "="); i > 0 {
				return cookie.Value[:i]
			}
			return cookie.Value
		}
	}
	return ""
}

func (c Client) path(p string) string {
	return c.base.ResolveReference(&url.URL{Path: p}).String()
}

func NewClient(uri string) (*Client, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	j, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &Client{
		Client: &http.Client{Jar: j, Timeout: time.Second * 15},
		base:   u,
	}, nil
}

func checkStatus(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v: %v", resp.Request.URL, resp.Status)
	}
	return nil
}

func (c *Client) Login(user, pass string) error {
	resp, err := c.Get(c.path("login.jsp"))
	if err != nil {
		return err
	}
	resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return err
	}
	return c.LoginProcess(user, pass)
}

func (c *Client) LoginProcess(user, pass string) error {
	v := url.Values{}
	v.Set("nick", user)
	v.Set("passwd", pass)
	v.Set("csrf", c.csrf())
	resp, err := c.PostForm(c.path("login_process"), v)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return checkStatus(resp)
}

func (c *Client) Logout() error {
	v := url.Values{}
	v.Set("csrf", c.csrf())
	resp, err := c.PostForm(c.path("logout"), v)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return checkStatus(resp)
}

func BanParams(reason string, ban Ban, days int, posting, captcha bool) url.Values {
	v := url.Values{}
	v.Set("reason", reason)
	v.Set("time", fmt.Sprint(ban))
	if ban == Custom {
		v.Set("ban_days", fmt.Sprint(days))
	}
	v.Set("allow_posting", fmt.Sprint(posting))
	v.Set("captcha_required", fmt.Sprint(captcha))
	return v
}

func (c *Client) BanIP(ip net.IP, v url.Values) error {
	if ip == nil {
		return fmt.Errorf("empty IP")
	}
	v.Set("csrf", c.csrf())
	v.Set("ip", fmt.Sprint(ip))
	resp, err := c.PostForm(c.path("banip.jsp"), v)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return checkStatus(resp)
}
