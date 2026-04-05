package server

import (
	"io"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/alexisg24/groxy/core/config"
	http_handler "github.com/alexisg24/groxy/core/http-handler"
)

func TestRootHandler_ProxyCalled(t *testing.T) {

	err := config.GlobalConfig.Load("../file-loader/config.test.yaml")
	if err != nil {
		panic(err)
	}
	hostConfig, exists := config.GlobalConfig.GetItem("backend_url")
	if !exists {
		t.Fatalf("expected hostConfig to exist")
	}
	var wg sync.WaitGroup
	wg.Add(1)

	orig := ProxyRequestFunc
	defer func() { ProxyRequestFunc = orig }()

	ProxyRequestFunc = func(h http_handler.HttpHandler) {
		defer wg.Done()
		h.Res.WriteHeader(200)
		h.Res.Write([]byte("ok"))
	}

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	RouteHandler(rec, req, hostConfig)

	wg.Wait()

	res := rec.Result()
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}
	if string(body) != "ok" {
		t.Fatalf("expected body 'ok', got %q", string(body))
	}
}
