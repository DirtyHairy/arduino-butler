package mock

import "net/http"

type MockResponseWriter struct {
	Headers map[string][]string
	Code    int
}

func (rw MockResponseWriter) Header() http.Header {
	return rw.Headers
}

func (rw MockResponseWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func (rw *MockResponseWriter) WriteHeader(code int) {
	rw.Code = code
}

func CreateMockResponseWriter() MockResponseWriter {
	return MockResponseWriter{make(map[string][]string), -1}
}
