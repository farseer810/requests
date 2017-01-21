package requests

// Implementations of the seven HTTP methods
func (this *session) Get(urlPath string) (Request, error) {
	return newRequest("GET", urlPath, this)
}

func (this *session) Post(urlPath string) (Request, error) {
	return newRequest("POST", urlPath, this)
}

func (this *session) Put(urlPath string) (Request, error) {
	return newRequest("PUT", urlPath, this)
}

func (this *session) Delete(urlPath string) (Request, error) {
	return newRequest("DELETE", urlPath, this)
}

func (this *session) Head(urlPath string) (Request, error) {
	return newRequest("HEAD", urlPath, this)
}

func (this *session) Options(urlPath string) (Request, error) {
	return newRequest("OPTIONS", urlPath, this)
}

func (this *session) Trace(urlPath string) (Request, error) {
	return newRequest("TRACE", urlPath, this)
}
