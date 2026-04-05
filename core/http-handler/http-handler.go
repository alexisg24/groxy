package http_handler

import (
	"io"
	"net/http"
	"time"

	fileloader "github.com/alexisg24/groxy/core/file-loader"
)

type HttpHandler struct {
	Host    string
	Res     http.ResponseWriter
	Http    *http.Request
	Configs *fileloader.FileLoaderItem
}

func HandleProxyRequest(handlerOpts HttpHandler) {
	requestClient := &http.Client{
		Timeout: time.Duration(handlerOpts.Configs.Timeout) * time.Millisecond,
	}
	proxyRequest, err := http.NewRequest(handlerOpts.Http.Method, handlerOpts.Host, handlerOpts.Http.Body)

	// Set headers from the original request to the proxy request
	for header, values := range handlerOpts.Http.Header {
		proxyRequest.Header[header] = values
	}

	// If there are custom request headers in the config, set them on the proxy request
	for header, value := range handlerOpts.Configs.CustomRequestHeaders {
		proxyRequest.Header.Set(header, value)
	}

	response, err := requestClient.Do(proxyRequest)
	if err != nil {
		// Cath error timeout, connection refused
		if err, ok := err.(interface{ Timeout() bool }); ok && err.Timeout() {
			http.Error(handlerOpts.Res, "Timeout reached on: "+handlerOpts.Host, http.StatusGatewayTimeout)
			return
		}
		http.Error(handlerOpts.Res, "Error forwarding request: "+err.Error()+" on "+handlerOpts.Http.URL.String(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Copy the response headers to the original response
	for header, values := range response.Header {
		handlerOpts.Res.Header()[header] = values
	}

	// Set the status code and write the response body
	handlerOpts.Res.WriteHeader(response.StatusCode)
	// copy body to response
	io.Copy(handlerOpts.Res, response.Body)

}
