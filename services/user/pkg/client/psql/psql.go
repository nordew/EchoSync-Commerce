package psqlClient

import "github.com/jackc/pgx"

func NewPsqlClient(host, database, user, password string, port uint16) (*pgx.Conn, error) {
	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     host,
		Database: database,
		User:     user,
		Password: password,
		Port:     port,
	})
	if err != nil {
		return nil, err
	}

	return conn, nil
}
