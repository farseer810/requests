package requests

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type Request interface {
	SetHeader(string, ...string) Request
	Headers() map[string][]string

	SetUrlParam(string, ...string) Request
	UrlParams() map[string][]string
	UrlPath() string
	SetJSON(string) Request
	SetBody([]byte) Request

	SetBodyParam(string, ...string) Request
	BodyParams() map[string][]string

	AddFile(string, string, []byte) Request

	Send() (Response, error)
}

type formFile struct {
	filename string
	data     []byte
}

type request struct {
	*session
	method  string
	URL     string
	headers map[string][]string

	isJSON     bool
	body       []byte
	bodyParams map[string][]string
	urlParams  map[string][]string
	files      map[string]*formFile
}

func parseURL(urlPath string) (URL *url.URL, err error) {
	URL, err = url.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	if URL.Scheme != "http" && URL.Scheme != "https" {
		urlPath = "http://" + urlPath
		URL, err = url.Parse(urlPath)
		if err != nil {
			return nil, err
		}
		if URL.Scheme != "http" && URL.Scheme != "https" {
			return nil, errors.New("[package requests] only HTTP and HTTPS are accepted")
		}
	}
	return
}

func newRequest(method string, urlPath string, s *session) (Request, error) {
	URL, err := parseURL(urlPath)
	if err != nil {
		return nil, err
	}
	urlParams := make(map[string][]string)
	for key, values := range URL.Query() {
		urlParams[key] = values
	}
	urlPath = URL.Scheme + "://" + URL.Host + URL.Path
	r := &request{session: s, method: method, URL: urlPath}
	r.headers = make(map[string][]string)
	r.bodyParams = make(map[string][]string)
	r.urlParams = urlParams
	r.files = make(map[string]*formFile)
	return r, nil
}

func (self *request) SetHeader(key string, values ...string) Request {
	if len(values) > 0 {
		self.headers[key] = values[:]
	} else {
		delete(self.headers, key)
	}
	return self
}

func (self *request) Headers() map[string][]string {
	headers := make(map[string][]string)
	for key, values := range self.headers {
		headers[key] = values[:]
	}
	return headers
}

func (self *request) SetUrlParam(key string, values ...string) Request {
	if len(values) > 0 {
		self.urlParams[key] = values[:]
	} else {
		delete(self.urlParams, key)
	}
	return self
}

func (self *request) UrlParams() map[string][]string {
	params := make(map[string][]string)
	for key, values := range self.urlParams {
		params[key] = values[:]
	}
	return params
}

func (self *request) UrlPath() string {
	return self.URL + "?" + parseParams(self.urlParams).Encode()
}

func (self *request) SetJSON(json string) Request {
	self.isJSON = true
	self.body = []byte(json)
	return self
}

func (self *request) SetBody(body []byte) Request {
	self.isJSON = false
	self.body = body
	return self
}

func (self *request) SetBodyParam(key string, values ...string) Request {
	if len(values) > 0 {
		self.bodyParams[key] = values[:]
	} else {
		delete(self.bodyParams, key)
	}
	return self
}

func (self *request) BodyParams() map[string][]string {
	params := make(map[string][]string)
	for key, values := range self.urlParams {
		params[key] = values[:]
	}
	return params
}

func (self *request) AddFile(fieldname string, filename string, data []byte) Request {
	if fieldname != "" && filename != "" && data != nil {
		self.files[fieldname] = &formFile{filename: fieldname, data: data}
	}
	return self
}

func (self *request) parseBody() (req *http.Request, err error) {
	// GET and TRACE request should not have a message body
	if self.method == "GET" || self.method == "TRACE" {
		req, err = http.NewRequest(self.method, self.UrlPath(), nil)
	}

	// Process message body
	if len(self.body) > 0 {
		if self.isJSON {
			self.headers["Content-Type"] = []string{"application/json"}
			req, err = http.NewRequest(self.method, self.UrlPath(),
				strings.NewReader(string(self.body)))
		} else {
			var body *bytes.Buffer
			body = bytes.NewBuffer(self.body)
			req, err = http.NewRequest(self.method, self.UrlPath(), body)
		}
	} else if len(self.files) > 0 {
		// multipart
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		var part io.Writer
		for fieldname, file := range self.files {
			part, err = writer.CreateFormFile(fieldname, file.filename)
			if err != nil {
				return
			}
			_, err = part.Write(file.data)
			if err != nil {
				return
			}
		}
		for fieldname, values := range self.bodyParams {
			temp := make(map[string][]string)
			temp[fieldname] = values

			value := parseParams(temp).Encode()
			err = writer.WriteField(fieldname, value)
			if err != nil {
				return
			}
		}
		err = writer.Close()
		if err != nil {
			return
		}
		self.headers["Content-Type"] = []string{writer.FormDataContentType()}
		req, err = http.NewRequest(self.method, self.UrlPath(), body)
	} else {
		self.headers["Content-Type"] = []string{"application/x-www-form-urlencoded"}
		req, err = http.NewRequest(self.method, self.UrlPath(),
			strings.NewReader(parseParams(self.bodyParams).Encode()))
	}
	return
}

func (self *request) Send() (res Response, err error) {
	req, err := self.parseBody()
	if err != nil {
		return
	}
	self.session.setCookies(req.URL)
	req.Header = parseHeaders(self.headers)
	httpResponse, err := self.session.Do(req)
	if err != nil {
		return
	}
	res, err = newResponse(httpResponse)
	return
}
