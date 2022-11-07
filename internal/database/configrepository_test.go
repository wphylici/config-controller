package database

import (
	"encoding/json"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/wphylici/contest-cloud/internal/models"
	"regexp"
	"testing"
)

func TestCreate(t *testing.T) {

	dbmock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer dbmock.Close()

	r := &ServiceConfigRepository{
		psql: &PostgreSQL{
			db: dbmock,
		},
	}

	type args struct {
		sv *models.ServiceConfig
	}
	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		expectsv     *models.ServiceConfig
		wantError    bool
	}{
		{
			name: "OK",
			args: args{
				sv: &models.ServiceConfig{
					ID:      1,
					Version: 1,
					Service: "test1",
					Data: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
			},
			expectsv: &models.ServiceConfig{
				ID:      1,
				Version: 1,
				Service: "test1",
				Data: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := mock.NewRows([]string{"exist"}).AddRow(false)
				query := regexp.QuoteMeta("SELECT EXISTS(SELECT service FROM configs WHERE service=$1)")
				mock.ExpectQuery(query).
					WithArgs(args.sv.Service).WillReturnRows(rows)

				rows = mock.NewRows([]string{"id"}).AddRow(1)
				query = regexp.QuoteMeta("INSERT INTO configs (service) VALUES ($1) RETURNING id")
				mock.ExpectQuery(query).
					WithArgs(args.sv.Service).WillReturnRows(rows)

				data, err := json.Marshal(args.sv.Data)
				if err != nil {
					t.Fatal(err)
				}

				query = regexp.QuoteMeta("INSERT INTO data_configs (config_id, version, data) VALUES ($1, $2, $3)")
				mock.ExpectQuery(query).
					WithArgs(args.sv.ID, args.sv.Version, data).WillReturnRows(&sqlmock.Rows{})

				mock.ExpectCommit()
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.Create(testCase.args.sv)
			if testCase.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectsv, got)
			}
		})
	}
}
