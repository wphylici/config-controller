package database

import (
	"database/sql"
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
		sc *models.ServiceConfig
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
				sc: &models.ServiceConfig{
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
					WithArgs(args.sc.Service).WillReturnRows(rows)

				rows = mock.NewRows([]string{"id"}).AddRow(1)
				query = regexp.QuoteMeta("INSERT INTO configs (service) VALUES ($1) RETURNING id")
				mock.ExpectQuery(query).
					WithArgs(args.sc.Service).WillReturnRows(rows)

				data, err := json.Marshal(args.sc.Data)
				if err != nil {
					t.Fatal(err)
				}

				query = regexp.QuoteMeta("INSERT INTO data_configs (config_id, version, data) VALUES ($1, $2, $3)")
				mock.ExpectQuery(query).
					WithArgs(args.sc.ID, args.sc.Version, data).WillReturnRows(&sqlmock.Rows{})

				mock.ExpectCommit()
			},
		},
		{
			name: "ConfigAlreadyBeenCreatedError",
			args: args{
				sc: &models.ServiceConfig{
					ID:      1,
					Version: 1,
					Service: "test1",
					Data: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
			},
			expectsv:  nil,
			wantError: true,
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := mock.NewRows([]string{"exist"}).AddRow(true)
				query := regexp.QuoteMeta("SELECT EXISTS(SELECT service FROM configs WHERE service=$1)")
				mock.ExpectQuery(query).
					WithArgs(args.sc.Service).WillReturnRows(rows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.Create(testCase.args.sc)
			if testCase.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectsv, got)
			}
		})
	}
}

func TestRead(t *testing.T) {

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

	r.Create(&models.ServiceConfig{
		Service: "test1",
		Data: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	})

	type args struct {
		sc *models.ServiceConfig
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
			name: "OK Last Record",
			args: args{
				sc: &models.ServiceConfig{
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

				configData, err := json.Marshal(args.sc.Data)
				if err != nil {
					t.Fatal(err)
				}

				rows := mock.NewRows([]string{"id"}).AddRow(1)
				query := regexp.QuoteMeta("SELECT id FROM configs WHERE service=$1")
				mock.ExpectQuery(query).
					WithArgs(args.sc.Service).WillReturnRows(rows)

				rows = mock.NewRows([]string{"data", "version"}).AddRow(configData, 1)
				query = regexp.QuoteMeta("SELECT data, version FROM data_configs WHERE config_id=$1 ORDER BY version DESC LIMIT 1")
				mock.ExpectQuery(query).
					WithArgs(1).WillReturnRows(rows)
			},
		},
		{
			name: "OK Specific Version",
			args: args{
				sc: &models.ServiceConfig{
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

				configData, err := json.Marshal(args.sc.Data)
				if err != nil {
					t.Fatal(err)
				}

				rows := mock.NewRows([]string{"id"}).AddRow(1)
				query := regexp.QuoteMeta("SELECT id FROM configs WHERE service=$1")
				mock.ExpectQuery(query).
					WithArgs(args.sc.Service).WillReturnRows(rows)

				rows = mock.NewRows([]string{"data"}).AddRow(configData)
				query = regexp.QuoteMeta("SELECT data FROM data_configs WHERE (config_id=$1) AND (version=$2)")
				mock.ExpectQuery(query).
					WithArgs(1, 1).WillReturnRows(rows)
			},
		},
		{
			name: "ConfigForServiceNotFound",
			args: args{
				sc: &models.ServiceConfig{
					Service: "dont-exist",
				},
			},
			wantError: true,
			mockBehavior: func(args args) {
				query := regexp.QuoteMeta("SELECT id FROM configs WHERE service=$1")
				mock.ExpectQuery(query).
					WithArgs(args.sc.Service).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "ConfigVersionNotFound",
			args: args{
				sc: &models.ServiceConfig{
					Version: 2,
					Service: "test1",
				},
			},
			wantError: true,
			mockBehavior: func(args args) {

				rows := mock.NewRows([]string{"id"}).AddRow(1)
				query := regexp.QuoteMeta("SELECT id FROM configs WHERE service=$1")
				mock.ExpectQuery(query).
					WithArgs(args.sc.Service).WillReturnRows(rows)

				query = regexp.QuoteMeta("SELECT data FROM data_configs WHERE (config_id=$1) AND (version=$2)")
				mock.ExpectQuery(query).
					WithArgs(1, 2).WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.Read(testCase.args.sc)
			if testCase.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectsv, got)
			}
		})
	}
}

func TestUpdate(t *testing.T) {

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

	r.Create(&models.ServiceConfig{
		ID:      1,
		Version: 1,
		Service: "test1",
		Data: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	})

	type args struct {
		sc *models.ServiceConfig
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
				sc: &models.ServiceConfig{
					Service: "test1",
					Data: map[string]string{
						"key1": "value1",
						"key2": "value2",
						"key3": "value3",
					},
				},
			},
			expectsv: &models.ServiceConfig{
				ID:      1,
				Version: 2,
				Service: "test1",
				Data: map[string]string{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
				},
			},
			mockBehavior: func(args args) {

				rows := mock.NewRows([]string{"exist"}).AddRow(true)
				query := regexp.QuoteMeta("SELECT EXISTS(SELECT service FROM configs WHERE service=$1)")
				mock.ExpectQuery(query).
					WithArgs(args.sc.Service).WillReturnRows(rows)

				rows = mock.NewRows([]string{"id"}).AddRow(1)
				query = regexp.QuoteMeta("SELECT id FROM configs WHERE service=$1")
				mock.ExpectQuery(query).
					WithArgs(args.sc.Service).WillReturnRows(rows)

				lastVersionConf := map[string]string{
					"key1": "value1",
					"key2": "value2",
				}

				lastVersionData, err := json.Marshal(lastVersionConf)
				if err != nil {
					t.Fatal(err)
				}

				rows = mock.NewRows([]string{"version", "data"}).AddRow(1, lastVersionData)
				query = regexp.QuoteMeta("SELECT version, data FROM data_configs WHERE config_id=$1 ORDER BY version DESC LIMIT 1")
				mock.ExpectQuery(query).
					WithArgs(1).WillReturnRows(rows)

				data, err := json.Marshal(args.sc.Data)
				if err != nil {
					t.Fatal(err)
				}

				query = regexp.QuoteMeta("INSERT INTO data_configs (config_id, version, data) VALUES ($1, $2, $3)")
				mock.ExpectQuery(query).
					WithArgs(1, 2, data).WillReturnRows(&sqlmock.Rows{})
			},
		},
		{
			name: "ConfigForServiceNotFound",
			args: args{
				sc: &models.ServiceConfig{
					Service: "dont-exist",
				},
			},
			wantError: true,
			mockBehavior: func(args args) {

				rows := mock.NewRows([]string{"exist"}).AddRow(false)
				query := regexp.QuoteMeta("SELECT EXISTS(SELECT service FROM configs WHERE service=$1)")
				mock.ExpectQuery(query).
					WithArgs(args.sc.Service).WillReturnRows(rows)
			},
		},
		{
			name: "NoChangeInConfigError",
			args: args{
				sc: &models.ServiceConfig{
					Service: "test1",
					Data: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
			},
			wantError: true,
			mockBehavior: func(args args) {

				rows := mock.NewRows([]string{"exist"}).AddRow(true)
				query := regexp.QuoteMeta("SELECT EXISTS(SELECT service FROM configs WHERE service=$1)")
				mock.ExpectQuery(query).
					WithArgs(args.sc.Service).WillReturnRows(rows)

				rows = mock.NewRows([]string{"id"}).AddRow(1)
				query = regexp.QuoteMeta("SELECT id FROM configs WHERE service=$1")
				mock.ExpectQuery(query).
					WithArgs(args.sc.Service).WillReturnRows(rows)

				lastVersionConf := map[string]string{
					"key1": "value1",
					"key2": "value2",
				}

				lastVersionData, err := json.Marshal(lastVersionConf)
				if err != nil {
					t.Fatal(err)
				}

				rows = mock.NewRows([]string{"version", "data"}).AddRow(1, lastVersionData)
				query = regexp.QuoteMeta("SELECT version, data FROM data_configs WHERE config_id=$1 ORDER BY version DESC LIMIT 1")
				mock.ExpectQuery(query).
					WithArgs(1).WillReturnRows(rows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.Update(testCase.args.sc)
			if testCase.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectsv, got)
			}
		})
	}
}

func TestDelete(t *testing.T) {

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

	r.Create(&models.ServiceConfig{
		Service: "test1",
		Data: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	})

	r.Update(&models.ServiceConfig{
		Service: "test1",
		Data: map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		},
	})

	type args struct {
		sc *models.ServiceConfig
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
			name: "ConfigForServiceNotFoundError",
			args: args{
				sc: &models.ServiceConfig{
					Service: "dont-exist",
				},
			},
			wantError: true,
			mockBehavior: func(args args) {

				query := regexp.QuoteMeta("SELECT id FROM configs WHERE service=$1")
				mock.ExpectQuery(query).
					WithArgs(args.sc.Service).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "ConfigVersionNotFoundError",
			args: args{
				sc: &models.ServiceConfig{
					Service: "test1",
					Version: 5,
				},
			},
			wantError: true,
			mockBehavior: func(args args) {

				rows := mock.NewRows([]string{"id"}).AddRow(1)
				query := regexp.QuoteMeta("SELECT id FROM configs WHERE service=$1")
				mock.ExpectQuery(query).
					WithArgs(args.sc.Service).WillReturnRows(rows)

				query = regexp.QuoteMeta("DELETE FROM data_configs WHERE (config_id=$1) AND (version=$2) RETURNING config_id")
				mock.ExpectQuery(query).
					WithArgs(1, args.sc.Version).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "OK Delete specific config version",
			args: args{
				sc: &models.ServiceConfig{
					Version: 2,
					Service: "test1",
				},
			},
			expectsv: &models.ServiceConfig{
				ID:      1,
				Version: 2,
				Service: "test1",
			},
			mockBehavior: func(args args) {

				rows := mock.NewRows([]string{"id"}).AddRow(1)
				query := regexp.QuoteMeta("SELECT id FROM configs WHERE service=$1")
				mock.ExpectQuery(query).
					WithArgs(args.sc.Service).WillReturnRows(rows)

				rows = mock.NewRows([]string{"id"}).AddRow(1)
				query = regexp.QuoteMeta("DELETE FROM data_configs WHERE (config_id=$1) AND (version=$2) RETURNING config_id")
				mock.ExpectQuery(query).
					WithArgs(1, args.sc.Version).WillReturnRows(rows)
			},
		},
		{
			name: "OK Delete all configs",
			args: args{
				sc: &models.ServiceConfig{
					Version: 0,
					Service: "test1",
				},
			},
			expectsv: &models.ServiceConfig{
				ID:      1,
				Version: 0,
				Service: "test1",
			},
			mockBehavior: func(args args) {

				rows := mock.NewRows([]string{"id"}).AddRow(1)
				query := regexp.QuoteMeta("SELECT id FROM configs WHERE service=$1")
				mock.ExpectQuery(query).
					WithArgs(args.sc.Service).WillReturnRows(rows)

				mock.ExpectBegin()

				rows = mock.NewRows([]string{"id"}).AddRow(1)
				query = regexp.QuoteMeta("DELETE FROM data_configs WHERE config_id=$1")
				mock.ExpectQuery(query).
					WithArgs(1).WillReturnRows(rows)

				rows = mock.NewRows([]string{"id"}).AddRow(1)
				query = regexp.QuoteMeta("DELETE FROM configs WHERE id=$1")
				mock.ExpectQuery(query).
					WithArgs(1).WillReturnRows(rows)

				mock.ExpectCommit()
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.Delete(testCase.args.sc)
			if testCase.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectsv, got)
			}
		})
	}
}
