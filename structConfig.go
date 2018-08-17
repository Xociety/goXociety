package main

type config struct {
	ServerAddrBind             string `json:"server_addr_bind"`
	ServerPort                 int    `json:"server_port"`
	ServerSecretFolderPath     string `json:"server_secret_folder_path"`
	ServerSecretCertFilename   string `json:"server_secret_cert_filename"`
	ServerSecretKeyFilename    string `json:"server_secret_key_filename"`
	GCPBucketRootCloudStorage  string `json:"gcp_bucket_root_cloud_storage"`
	GCPSecretFolderPath        string `json:"gcp_secret_folder_path"`
	GCPSecretFilename          string `json:"gcp_secret_filename"`
	PostgresConfigFolderPath   string `json:"postgres_config_folder_path"`
	PostgresConStr             string `json:"postgres_con_str"`
	PostgresSecretFolderPath   string `json:"postgres_secret_folder_path"`
	PostgresSecretAuthFilename string `json:"postgres_secret_auth_filename"`
	MongoConStr                string `json:"mongo_con_str"`
	MongoSecretFolderPath      string `json:"mongo_secret_folder_path"`
	MongoSecretCertFilename    string `json:"mongo_secret_cert_filename"`
	MongoSecretKeyFilename     string `json:"mongo_secret_key_filename"`
	MongoSecretAuthFilename    string `json:"mongo_secret_auth_filename"`
}
