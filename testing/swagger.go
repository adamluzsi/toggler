package testing

import (
	"github.com/adamluzsi/testcase"
	"github.com/toggler-io/toggler/extintf/httpintf/httpapi"
	"github.com/toggler-io/toggler/lib/go/client"
	"github.com/toggler-io/toggler/usecases"
	"net/http"
	"net/http/httptest"
	"net/url"
)

func SetupSpecWithSwaggerClient(s *testcase.Spec) {

	s.Let(`httpapi.Handler`, func(t *testcase.T) interface{} {
		return httpapi.NewHandler(usecases.NewUseCases(ExampleStorage(t)))
	})

	s.Let(`httptest.NewServer/httpapi`, func(t *testcase.T) interface{} {
		m := t.I(`httpapi.Handler`).(*httpapi.Handler)
		return httptest.NewServer(http.StripPrefix(`/api`, m))
	})

	s.After(func(t *testcase.T) {
		t.I(`httptest.NewServer/httpapi`).(*httptest.Server).Close()
	})

	s.Let(`swagger.Client`, func(t *testcase.T) interface{} {
		s := t.I(`httptest.NewServer/httpapi`).(*httptest.Server)
		tc := client.DefaultTransportConfig()
		u, _ := url.Parse(s.URL)
		tc.Host = u.Host
		tc.Schemes = []string{`http`}
		return client.NewHTTPClientWithConfig(nil, tc)
	})
}

func GetSwaggerClient(t *testcase.T) *client.ProvidesAPIOnHTTPLayerToTheTogglerService {
	return t.I(`swagger.Client`).(*client.ProvidesAPIOnHTTPLayerToTheTogglerService)
}