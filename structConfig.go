package main

type config struct {
	ServerAddrBind          string `json:"server_addr_bind,omitempty"`
	ServerPort              int    `json:"server_port,omitempty"`
	PostgresConStr          string `json:"postgres_con_str,omitempty"`
	MongoConStr             string `json:"mongo_con_str,omitempty"`
	ClientAuthGCPFolderPath string `json:"client_auth_gcp_folder_path,omitempty"`
	ServerCertFolderPath    string `json:"server_cert_folder_path,omitempty"`
}
