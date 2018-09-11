package config

type Secret struct {
	Postgres SecretPostgres
	Mongo    SecretMongo
}

type SecretPostgres struct {
	PostgresAuthStr string `json:"postgres_auth_str"`
}

type SecretMongo struct {
	MongoUsername string `json:"mongo_username"`
	MongoPassword string `json:"mongo_password"`
	MongoDatabase string `json:"mongo_database"`
}

type SecretMap struct {
	Key string `json:"key"`
}
