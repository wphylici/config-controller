package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/wphylici/contest-cloud/internal/models"
	"reflect"
)

type ServiceConfigRepository struct {
	psql *PostgreSQL
}

func (r *ServiceConfigRepository) Create(c *models.ServiceConfig) (*models.ServiceConfig, error) {
	var isServiceExist bool
	checkServiceExistQuery := fmt.Sprintf("IF EXISTS(SELECT service FROM configs WHERE service = '%s')", c.Service)

	if err := r.psql.db.QueryRow(checkServiceExistQuery).Scan(&isServiceExist); err != nil {
		return nil, err
	} else if !isServiceExist {
		if err = r.psql.db.QueryRow(
			"INSERT INTO configs (service) VALUES ($1) RETURNING id",
			c.Service,
		).Scan(&c.ID); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("config for service '%s' has already been created", c.Service)
	}

	configData, err := json.Marshal(c.Data)
	if err != nil {
		return nil, err
	}

	r.psql.db.QueryRow(
		"INSERT INTO data_configs (config_id, version, data) VALUES ($1, $2, $3)",
		c.ID,
		1,
		configData,
	)

	return c, nil
}

func (r *ServiceConfigRepository) Read(c *models.ServiceConfig) (*models.ServiceConfig, error) {
	checkServiceExistQuery := fmt.Sprintf("SELECT id FROM configs WHERE service = '%s'", c.Service)
	if err := r.psql.db.QueryRow(checkServiceExistQuery).Scan(&c.ID); err == sql.ErrNoRows {
		return nil, fmt.Errorf("config for service '%s' not found", c.Service)
	} else if err != nil {
		return nil, err
	}

	var configData []byte
	if c.Version == 0 {
		lastRecordQuery := fmt.Sprintf("SELECT data FROM data_configs WHERE config_id = %d ORDER BY version DESC LIMIT 1", c.ID)
		if err := r.psql.db.QueryRow(
			lastRecordQuery,
		).Scan(&configData); err != nil {
			return nil, err
		}
	} else {
		configDataQuery := fmt.Sprintf("SELECT data FROM data_configs WHERE (config_id = %d) AND (version = %d)", c.ID, c.Version)
		if err := r.psql.db.QueryRow(configDataQuery).Scan(&configData); err == sql.ErrNoRows {
			return nil, fmt.Errorf("version %d configuration for '%s' service not found", c.Version, c.Service)
		} else if err != nil {
			return nil, err
		}
	}

	if err := json.Unmarshal(configData, &c.Data); err != nil {
		return nil, err
	}

	return c, nil

}

func (r *ServiceConfigRepository) Update(c *models.ServiceConfig) (*models.ServiceConfig, error) {
	var isServiceExist bool
	checkServiceExistQuery := fmt.Sprintf("SELECT EXISTS(SELECT service FROM configs WHERE service = '%s')", c.Service)

	if err := r.psql.db.QueryRow(checkServiceExistQuery).Scan(&isServiceExist); err != nil {
		return nil, err
	} else if !isServiceExist {
		return nil, fmt.Errorf("config for service '%s' not found", c.Service)
	} else {
		checkServiceIDQuery := fmt.Sprintf("SELECT id FROM configs WHERE service = '%s'", c.Service)
		if err = r.psql.db.QueryRow(checkServiceIDQuery).Scan(&c.ID); err != nil {
			return nil, err
		}
	}

	var lastConfigData []byte
	lastRecordQuery := fmt.Sprintf("SELECT version, data FROM data_configs WHERE config_id = %d ORDER BY version DESC LIMIT 1", c.ID)
	if err := r.psql.db.QueryRow(
		lastRecordQuery,
	).Scan(&c.Version, &lastConfigData); err != nil {
		return nil, err
	}
	c.Version++

	configData, err := json.Marshal(c.Data)
	if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(lastConfigData, configData) {
		r.psql.db.QueryRow(
			"INSERT INTO data_configs (config_id, version, data) VALUES ($1, $2, $3)",
			c.ID,
			c.Version,
			configData,
		)
	} else {
		return nil, fmt.Errorf("no change in config")
	}

	return c, nil
}

func (r *ServiceConfigRepository) Delete(c *models.ServiceConfig) (*models.ServiceConfig, error) {
	return c, nil
}
