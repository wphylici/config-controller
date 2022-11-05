package app

import "github.com/wphylici/contest-cloud/internal/database"

func StartPostgreSQL(config *database.Config) error {

	database.Psql = database.New(config)
	err := database.Psql.Open()
	if err != nil {
		return err
	}

	return nil
}
