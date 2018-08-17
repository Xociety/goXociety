package main

type secret struct {
	Postgres secretPostgres
	Mongo    secretMongo
}

type secretPostgres struct {
	PostgresAuthStr string `json:"postgres_auth_str"`
}

type secretMongo struct {
	MongoUsername string `json:"mongo_username"`
	MongoPassword string `json:"mongo_password"`
	MongoDatabase string `json:"mongo_database"`
}
