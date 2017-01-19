package requests

func (self *session) Get(urlPath string) (Request, error) {
	return newRequest("GET", urlPath, self)
}

func (self *session) Post(urlPath string) (Request, error) {
	return newRequest("POST", urlPath, self)
}

func (self *session) Put(urlPath string) (Request, error) {
	return newRequest("PUT", urlPath, self)
}

func (self *session) Delete(urlPath string) (Request, error) {
	return newRequest("DELETE", urlPath, self)
}

func (self *session) Head(urlPath string) (Request, error) {
	return newRequest("HEAD", urlPath, self)
}

func (self *session) Options(urlPath string) (Request, error) {
	return newRequest("OPTIONS", urlPath, self)
}

func (self *session) Trace(urlPath string) (Request, error) {
	return newRequest("TRACE", urlPath, self)
}
