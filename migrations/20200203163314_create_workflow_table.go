package migrations

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("CREATE TABLE WORKFLOWS ( ID uuid, SERVICE_FLOW JSONB, NAME VARCHAR NOT NULL, PRIMARY KEY(NAME));")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE WORKFLOWS;")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20200203163314_create_workflow_table", up, down, opts)
}
