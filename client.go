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

const lorURL = "https://www.linux.org.ru/"

type Client struct {
	*http.Client
	*url.URL
	login        string
	loginProcess string
	banip        string
	logout       string
	token        string
	user         string
	pass         string
}

func (c Client) csrf() string {
	for _, cookie := range c.Jar.Cookies(c.URL) {
		if cookie.Name == "CSRF_TOKEN" {
			if i := strings.Index(cookie.Value, "="); i > 0 {
				return cookie.Value[:i]
			}
			return cookie.Value
		}
	}
	return ""
}

func NewClient() (*Client, error) {
	u, err := url.Parse(lorURL)
	if err != nil {
		return nil, err
	}
	j, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &Client{
		Client:       &http.Client{Jar: j, Timeout: time.Second * 15},
		URL:          u,
		login:        fmt.Sprint(u, "login.jsp"),
		loginProcess: fmt.Sprint(u, "login_process"),
		banip:        fmt.Sprint(u, "banip.jsp"),
		logout:       fmt.Sprint(u, "logout"),
	}, nil
}

func check(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v: %v", resp.Request.URL, resp.Status)
	}
	return nil
}

func (c *Client) Login(user, pass string) error {
	resp, err := c.Get(c.login)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if err := check(resp); err != nil {
		return err
	}
	c.token = c.csrf()
	c.user = user
	c.pass = pass
	return c.LoginProcess()
}

func (c *Client) LoginProcess() error {
	v := url.Values{}
	v.Set("nick", c.user)
	v.Set("passwd", c.pass)
	v.Set("csrf", c.token)
	resp, err := c.PostForm(c.loginProcess, v)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return check(resp)
}

func (c *Client) Logout() error {
	v := url.Values{}
	v.Set("csrf", c.token)
	resp, err := c.PostForm(c.logout, v)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return check(resp)
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
	v.Set("csrf", c.token)
	v.Set("ip", fmt.Sprint(ip))
	resp, err := c.PostForm(c.banip, v)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return check(resp)
}
