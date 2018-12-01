package models

import (
	"bytes"
	"database/sql"
	"html/template"

	"github.com/alexandrevilain/postgrest-auth/pkg/config"
	"github.com/labstack/gommon/log"
)

// EnsureDBElementsExists ensure that required tables/roles/schemas, exists on the database
func EnsureDBElementsExists(db *sql.DB, config *config.DB, logger *log.Logger) error {
	tmpl, err := template.New("sql").Parse(`
	CREATE SCHEMA IF NOT EXISTS auth;
	CREATE TABLE IF NOT EXISTS auth.users (
		id uuid PRIMARY KEY NOT NULL,
		email text UNIQUE NOT NULL UNIQUE,
		password text NOT NULL,
		confirmed boolean NOT NULL DEFAULT FALSE,
		confirmToken uuid DEFAULT NULL,
		resetPasswordToken text DEFAULT NULL
	);
	DO
	$body$
	BEGIN
	IF NOT EXISTS (
		SELECT
		FROM pg_roles
		WHERE rolname = '{{ .Roles.Anonymous }}') THEN
		CREATE ROLE {{ .Roles.Anonymous }} NOLOGIN;
	END IF;
	END
	$body$;
	DO
	$body$
	BEGIN
	IF NOT EXISTS (
		SELECT
		FROM pg_roles
		WHERE rolname = '{{ .Roles.User }}') THEN
		CREATE ROLE {{ .Roles.User }} NOLOGIN;
	END IF;
	END
	$body$;
	GRANT USAGE ON SCHEMA auth TO {{ .Roles.Anonymous }}, {{ .Roles.User }};
	
	CREATE OR REPLACE FUNCTION auth.current_user_id() RETURNS uuid
	LANGUAGE plpgsql
	AS $$
	BEGIN
		RETURN current_setting('request.jwt.claim.userid', true)::uuid;
	EXCEPTION
		-- handle unrecognized configuration parameter error
		WHEN undefined_object THEN RETURN '';
	END;
	$$;
	GRANT EXECUTE ON FUNCTION current_user_id() TO {{ .Roles.User }};
	`)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, config)
	if err != nil {
		return err
	}
	logger.Debugf("Executing the following query: \n %v \n", buf.String())
	_, err = db.Exec(buf.String())
	return err
}
