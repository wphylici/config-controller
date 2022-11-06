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

func getConfigForServiceNotFoundError(serviceName string) string {
	return fmt.Sprintf("config for service '%s' not found", serviceName)
}

func getConfigVersionNotFoundError(serviceName string, version uint32) string {
	return fmt.Sprintf("config version '%d' for '%s' service not found", version, serviceName)
}

func getConfigAlreadyBeenCreatedError(serviceName string) string {
	return fmt.Sprintf("config for service '%s' has already been created", serviceName)
}

func getNoChangeInConfigError() string {
	return fmt.Sprintf("no change in config")
}

func (r *ServiceConfigRepository) Create(c *models.ServiceConfig) (*models.ServiceConfig, error) {
	var isServiceExist bool

	if err := r.psql.db.QueryRow("SELECT EXISTS(SELECT service FROM configs WHERE service=$1)",
		c.Service,
	).Scan(&isServiceExist); err != nil {
		return nil, err
	} else if !isServiceExist {
		if err = r.psql.db.QueryRow(
			"INSERT INTO configs (service) VALUES ($1) RETURNING id",
			c.Service,
		).Scan(&c.ID); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf(getConfigAlreadyBeenCreatedError(c.Service))
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
	if err := r.psql.db.QueryRow("SELECT id FROM configs WHERE service=$1",
		c.Service,
	).Scan(&c.ID); err == sql.ErrNoRows {
		return nil, fmt.Errorf(getConfigForServiceNotFoundError(c.Service))
	} else if err != nil {
		return nil, err
	}

	var configData []byte
	if c.Version == 0 {
		if err := r.psql.db.QueryRow("SELECT data FROM data_configs WHERE config_id=$1 ORDER BY version DESC LIMIT 1",
			c.ID,
		).Scan(&configData); err != nil {
			return nil, err
		}
	} else {
		if err := r.psql.db.QueryRow("SELECT data FROM data_configs WHERE (config_id=$1) AND (version=$2)",
			c.ID,
			c.Service,
		).Scan(&configData); err == sql.ErrNoRows {
			return nil, fmt.Errorf(getConfigVersionNotFoundError(c.Service, c.Version))
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

	if err := r.psql.db.QueryRow("SELECT EXISTS(SELECT service FROM configs WHERE service=$1)",
		c.Service,
	).Scan(&isServiceExist); err != nil {
		return nil, err
	} else if !isServiceExist {
		return nil, fmt.Errorf(getConfigForServiceNotFoundError(c.Service))
	} else {
		if err = r.psql.db.QueryRow("SELECT id FROM configs WHERE service=$1",
			c.Service,
		).Scan(&c.ID); err != nil {
			return nil, err
		}
	}

	var lastConfigData []byte
	if err := r.psql.db.QueryRow("SELECT version, data FROM data_configs WHERE config_id=$1 ORDER BY version DESC LIMIT 1",
		c.ID,
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
		return nil, fmt.Errorf(getNoChangeInConfigError())
	}

	return c, nil
}

func (r *ServiceConfigRepository) Delete(c *models.ServiceConfig) (*models.ServiceConfig, error) {
	if err := r.psql.db.QueryRow("SELECT id FROM configs WHERE service=$1",
		c.Service,
	).Scan(&c.ID); err == sql.ErrNoRows {
		return nil, fmt.Errorf(getConfigForServiceNotFoundError(c.Service))
	} else if err != nil {
		return nil, err
	}

	if c.Version == 0 {
		if row := r.psql.db.QueryRow("DELETE FROM data_configs WHERE config_id=$1",
			c.ID,
		); row.Err() != nil {
			return nil, row.Err()
		}

		if row := r.psql.db.QueryRow("DELETE FROM configs WHERE id=$1",
			c.ID,
		); row.Err() != nil {
			return nil, row.Err()
		}
	} else {
		if err := r.psql.db.QueryRow("DELETE FROM data_configs WHERE (config_id=$1) AND (version=$2) RETURNING config_id",
			c.ID,
			c.Version,
		).Scan(&c.ID); err == sql.ErrNoRows {
			return nil, fmt.Errorf(getConfigVersionNotFoundError(c.Service, c.Version))
		} else if err != nil {
			return nil, err
		}
	}

	return c, nil
}
