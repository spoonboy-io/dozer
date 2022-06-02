package hook

import (
	"context"
	"fmt"

	"github.com/spoonboy-io/koan"

	"github.com/spoonboy-io/dozer/internal"
)

/*

cx, cancel := context.WithCancel(context.Background())
    req, _ := http.NewRequest("GET", "http://google.com", nil)
    req = req.WithContext(cx)
    ch := make(chan error)

    go func() {
        _, err := http.DefaultClient.Do(req)
        select {
        case <-cx.Done():
            // Already timedout
        default:
            ch <- err
        }
    }()
*/

func fireWebhook(process *internal.Process, hook *Hook, logger *koan.Logger, ctx context.Context) error {

	// we WILL parse the hook config
	// we WILL build HTTP the client
	// we WILL parse the requestBody param as a template and inject internal.Process params where needed
	// we WILL get the parsed requestBody
	// we WILL form a HTTP request and pass it to the client
	// we will send it and return for caller to log errors

	fmt.Println("fireWebhook")

	return nil
}
