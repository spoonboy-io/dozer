package morpheus_test

import (
	"reflect"
	"testing"

	"github.com/spoonboy-io/dozer/internal"

	"github.com/spoonboy-io/dozer/internal/morpheus"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetProcessTypes(t *testing.T) {
	// test expectations
	wantProcessTypes := map[string]string{
		"testProcessTypeCode1": "test process type name 1",
		"testProcessTypeCode2": "test process type name 2",
		"testProcessTypeCode3": "test process type name 3",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// add some mock rows
	rows := sqlmock.NewRows([]string{"id", "code", "name", "image_code"}).
		AddRow(1, "testProcessTypeCode1", "test process type name 1", "not used").
		AddRow(2, "testProcessTypeCode2", "test process type name 2", "not used").
		AddRow(3, "testProcessTypeCode3", "test process type name 3", "not used")

	mock.ExpectQuery("SELECT id, code, name, image_code FROM process_type;").WillReturnRows(rows)

	gotProcessTypes := map[string]string{}
	if err := morpheus.GetProcessTypes(db, gotProcessTypes); err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// test to see if we have them locally in map
	if !reflect.DeepEqual(gotProcessTypes, wantProcessTypes) {
		t.Errorf("failed got %v wanted %v", gotProcessTypes, wantProcessTypes)
	}

	// test to see if they are available in the internal namespace map
	if !reflect.DeepEqual(internal.ProcessTypes, wantProcessTypes) {
		t.Errorf("failed got %v wanted %v", internal.ProcessTypes, wantProcessTypes)
	}

	// check expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
