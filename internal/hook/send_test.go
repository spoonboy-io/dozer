package hook

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spoonboy-io/koan"

	"github.com/spoonboy-io/dozer/internal"
)

func Test_fireWebhook(t *testing.T) {
	ctx := context.Background()
	logger := &koan.Logger{}
	process := &internal.Process{}
	hook := &Hook{}

	// test good response
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`ok`))
	}))
	defer server.Close()
	hook.URL = server.URL
	hook.Description = "test hook"

	if err := fireWebhook(ctx, process, hook, logger); err != nil {
		t.Errorf("fail %v", err)
	}
	server.Close()

	// test bad response (404)
	server2 := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
	}))
	hook.URL = server2.URL
	hook.Description = "test hook - bad should be 404 and error"

	if err := fireWebhook(ctx, process, hook, logger); err == nil {
		t.Errorf("fail expected an error because the server did not return 200")
	}
	server2.Close()
}

func Test_processRequestBody(t *testing.T) {

}
