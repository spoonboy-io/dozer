package hook

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

	// test token is correctly sent & set
	wantToken := "BEARER AB123XXXXXXXXXZZZ"
	gotToken := ""
	server3 := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		gotToken = req.Header.Get("Authorization")
		rw.Write([]byte(`ok`))
	}))
	hook.URL = server3.URL
	hook.Description = "test hook - bad should be 404 and error"
	hook.Token = wantToken

	if err := fireWebhook(ctx, process, hook, logger); err != nil {
		t.Errorf("fail %v", err)
	}

	if gotToken != wantToken {
		t.Errorf("fail wanted %v, got %v", gotToken, wantToken)
	}
	server3.Close()
}

func Test_processRequestBody(t *testing.T) {
	testCases := []struct {
		name       string
		process    *internal.Process
		body       string
		wantOutput string
	}{
		{
			"testing a random selection of properties are interpolated (1)",
			&internal.Process{
				Id:              1,
				Status:          "complete",
				ProcessTypeName: sql.NullString{String: "Test Process Name"},
				TaskName:        sql.NullString{String: "Test Task Name"},
				CreatedBy:       sql.NullString{String: "Test User"},
				AccountId:       sql.NullInt64{Int64: 2},
			},
			"{{.Id}}, {{.Status}}, {{.ProcessTypeName}}, {{.TaskName}}, {{.CreatedBy}}, {{.AccountId}}",
			"1, complete, Test Process Name, Test Task Name, Test User, 2",
		},

		{
			"testing a random selection of properties are interpolated (2)",
			&internal.Process{
				Success:              sql.NullBool{Bool: true},
				CreatedByDisplayName: sql.NullString{String: "Test User"},
				DisplayName:          sql.NullString{String: "Test User"},
				Input:                sql.NullString{String: "Test Input"},
				AppId:                sql.NullInt64{Int64: 2},
				Message:              sql.NullString{String: "Test Message"},
				RefType:              sql.NullString{String: "Test Ref"},
				JobTemplateId:        sql.NullInt64{Int64: 2},
				ContainerName:        sql.NullString{String: "Test Container Name"},
			},
			"{{.Success}}, {{.CreatedByDisplayName}}, {{.DisplayName}}, {{.Input}}, {{.AppId}}, {{.Message}}, {{.RefType}}, {{.JobTemplateId}}, {{.ContainerName}}",
			"true, Test User, Test User, Test Input, 2, Test Message, Test Ref, 2, Test Container Name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// set the package config
			gotReader, err := parseRequestBody(tc.process, tc.body)
			if err != nil {
				t.Fatalf("Unexpected error %v ", err)
			}

			gotOutput := new(strings.Builder)
			_, err = io.Copy(gotOutput, gotReader)
			if err != nil {
				t.Fatalf("Unexpected error %v ", err)
			}

			if gotOutput.String() != tc.wantOutput {
				t.Errorf("Failed\n got: %v\n wanted: %v", gotOutput, tc.wantOutput)
			}
		})
	}

}
