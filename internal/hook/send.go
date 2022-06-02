package hook

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/spoonboy-io/koan"

	"github.com/spoonboy-io/dozer/internal"
)

func fireWebhook(ctx context.Context, process *internal.Process, hook *Hook, logger *koan.Logger) error {
	var data io.Reader
	var err error

	// we WILL parse the hook config

	// parse RequestBody if required
	if hook.Method != "GET" && hook.RequestBody != "" {
		data, err = parseRequestBody(process, hook.RequestBody)
		if err != nil {
			return err
		}
	}

	// form the request, make and return any errors
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

func parseRequestBody(process *internal.Process, body string) (io.Reader, error) {

	fmt.Println("before: ", body)

	var parsed io.ReadWriter
	t := template.Must(template.New("body").Parse(body))
	_ = t
	/*if err := t.Execute(parsed, process); err != nil {
		return nil, err
	}*/

	fmt.Println("after: ", parsed)

	return parsed, nil
}
