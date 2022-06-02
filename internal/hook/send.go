package hook

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/spoonboy-io/koan"

	"github.com/spoonboy-io/dozer/internal"
)

func fireWebhook(ctx context.Context, process *internal.Process, hook *Hook, logger *koan.Logger) error {
	var data io.Reader

	// we WILL parse the hook config
	// we WILL build HTTP the client
	// we WILL parse the requestBody param as a template and inject internal.Process params where needed
	// we WILL get the parsed requestBody

	// form the request, make it, return any errors
	req, err := http.NewRequest(hook.Method, hook.URL, data)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad response (%d): Hook: %s, URL: %s", res.StatusCode, hook.Description, hook.URL)
	}

	return nil
}
