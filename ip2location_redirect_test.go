package traefik_plugin_ip2location_redirect

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type httpHandlerMock struct{}

func (h *httpHandlerMock) ServeHTTP(http.ResponseWriter, *http.Request){}

func TestIP2LocationRedirect(t *testing.T) {
	var err error
	i := &IP2LocationRedirect{
		next: &httpHandlerMock{},
		config: &Config{
			Regions: []string{"CN", "HK"},
			RedirectUrl: "https://github.com/jiangwennn/traefik_plugin_ip2location_redirect",
		},
	}

	i.db, err = OpenDB("IP2LOCATION-LITE-DB1.IPV6.BIN")
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "http://github.com", nil)
	req.RemoteAddr = "59.38.44.63:80"
	rw := httptest.NewRecorder()

	i.ServeHTTP(rw, req)

	c := rw.Code
	if c != http.StatusFound {
		t.Fatal("unexpected status code", c)
	}


}