package main

type config struct {
	ServerAddrBind          string `json:"server_addr_bind"`
	ServerPort              int    `json:"server_port"`
	PostgresConStr          string `json:"postgres_con_str"`
	MongoConStr             string `json:"mongo_con_str"`
	ClientAuthGCPFolderPath string `json:"client_auth_gcp_folder_path"`
	ServerCertFolderPath    string `json:"server_cert_folder_path"`
}
