package main

type secret struct {
	Postgres secretPostgres
}

type secretPostgres struct {
	PostgresAuthStr string `json:"postgres_auth_str"`
}
