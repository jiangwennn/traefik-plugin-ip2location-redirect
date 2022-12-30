package traefik_plugin_ip2location_redirect

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type httpHandlerMock struct{}

func (h *httpHandlerMock) ServeHTTP(http.ResponseWriter, *http.Request){}

func TestIP2LocationRedirect_1(t *testing.T) {
	var err error
	i := &IP2LocationRedirect{
		next: &httpHandlerMock{},
		config: &Config{
			Filename: "IP2LOCATION-LITE-DB1.IPV6.BIN",
			Regions: []string{"CN", "HK"},
			RedirectUrl: "https://github.com/jiangwennn/traefik_plugin_ip2location_redirect",
			ErrorHeader: "X-IP2LOCATION-REDIRECT-ERROR",
		},
	}

	i.db, err = OpenDB(i.config.Filename)
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
func TestIP2LocationRedirect_2(t *testing.T) {
	var err error
	i := &IP2LocationRedirect{
		next: &httpHandlerMock{},
		config: &Config{
			Filename: "IP2LOCATION-LITE-DB1.IPV6.BIN",
			Regions: []string{"CN"},
			RedirectUrl: "https://github.com/jiangwennn/traefik_plugin_ip2location_redirect",
			FromHeader: "X-Forwarded-For",
			Permanent: true,
			NoMatch: true,
		},
	}

	i.db, err = OpenDB(i.config.Filename)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "http://github.com", nil)
	req.Header.Set("X-Forwarded-For", "78.31.211.32, 47.102.25.42, 47.110.182.208")
	rw := httptest.NewRecorder()

	i.ServeHTTP(rw, req)

	c := rw.Code
	if c != http.StatusMovedPermanently {
		t.Fatal("unexpected status code", c)
	}
}