package requests

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Session interface {
	SetTimeout(time.Duration) Session
	Timeout() time.Duration
	SetCookies(map[string]string) Session
	Cookies() map[string]string

	Get(string) (Request, error)
	Post(string) (Request, error)
	Put(string) (Request, error)
	Delete(string) (Request, error)
	Head(string) (Request, error)
	Options(string) (Request, error)
	Trace(string) (Request, error)
}

type session struct {
	*http.Client
	cookies map[string]string
}

func NewSession() Session {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	s := &session{Client: client}
	s.cookies = make(map[string]string)
	return s
}

func (self *session) SetTimeout(timeout time.Duration) Session {
	self.Client.Timeout = timeout
	return self
}

func (self *session) Timeout() time.Duration {
	return self.Client.Timeout
}

func (self *session) SetCookies(cookies map[string]string) Session {
	self.cookies = make(map[string]string)
	for key, value := range cookies {
		self.cookies[key] = value
	}
	return self
}

func (self *session) Cookies() map[string]string {
	cookies := make(map[string]string)
	for key, value := range self.cookies {
		cookies[key] = value
	}
	return cookies
}

func (self *session) setCookies(URL *url.URL) {
	cookies := self.Jar.Cookies(URL)
	for key, value := range self.cookies {
		// only sets the cookie when none corresponding ones are presented
		found := false
		for _, cookie := range cookies {
			if cookie.Name == key {
				found = true
				break
			}
		}
		if !found {
			cookies = append(cookies, &http.Cookie{Name: key, Value: value, MaxAge: 3600 * 24 * 7})
		}
	}
	self.Jar.SetCookies(URL, cookies)
}
