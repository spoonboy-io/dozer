package hook

import (
	"context"
	"testing"

	"github.com/spoonboy-io/koan"

	"github.com/spoonboy-io/dozer/internal"
)

func Test_fireWebhook(t *testing.T) {
	ctx := context.Background()
	logger := &koan.Logger{}
	process := &internal.Process{}
	hook := &Hook{}
	if err := fireWebhook(process, hook, logger, ctx); err != nil {
		t.Errorf("fail")
	}
}

func Test_processRequestBody(t *testing.T) {

}
