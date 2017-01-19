package requests

import (
	// "fmt"
	"io/ioutil"
	"net/http"
)

type Response interface {
	StatusCode() int
	Headers() map[string][]string
	Protocol() string
	Body() []byte
	ContentLength() int64
}

type response struct {
	status         int
	protocol       string
	headers        map[string][]string
	body           []byte
	content_length int64
}

func newResponse(httpResponse *http.Response) (Response, error) {
	var headers map[string][]string = httpResponse.Header
	body, err := ioutil.ReadAll(httpResponse.Body)
	content_length := int64(len(body))
	httpResponse.Body.Close()
	if err != nil {
		return nil, err
	}
	res := &response{headers: headers,
		protocol:       httpResponse.Proto,
		status:         httpResponse.StatusCode,
		body:           body,
		content_length: int64(content_length)}
	return res, nil
}

func (res *response) StatusCode() int {
	return res.status
}

func (res *response) Headers() map[string][]string {
	return res.headers
}

func (res *response) Protocol() string {
	return res.protocol
}

func (res *response) Body() []byte {
	return res.body
}

func (res *response) ContentLength() int64 {
	return res.content_length
}
