package models

import (
	"encoding/json"
	"fmt"
)

type Data struct {
}

type ServiceConfig struct {
	ID      int
	Service string
	Version uint32
	Data    map[string]string
}

func (s *ServiceConfig) UnmarshalJSON(bytes []byte) error {
	config := &struct {
		ID      int                 `json:"-"`
		Service string              `json:"service"`
		Version uint32              `json:"-"`
		Data    []map[string]string `json:"data"`
	}{}

	err := json.Unmarshal(bytes, &config)
	if err != nil {
		return err
	}

	m := map[string]string{}

	for _, pairs := range config.Data {
		for k, v := range pairs {
			_, found := m[k]
			if found {
				return fmt.Errorf("duplicate key value in config data")
			}
			m[k] = v
		}
	}

	s.ID = config.ID
	s.Service = config.Service
	s.Version = config.Version
	s.Data = m

	return nil
}
