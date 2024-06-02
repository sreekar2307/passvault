package dtos

type DatabaseConfigKeys struct {
	Host                  string
	Port                  string
	Username              string
	Password              string
	Name                  string
	MaxOpenConnections    string
	MaxIdleConnections    string
	MaxIdleConnectionTime string
}

type EncryptionConfigKeys struct {
	Key     string `mapstructure:"key"`
	Version string `mapstructure:"version"`
}

type Encryption struct {
	Keys []EncryptionConfigKeys `mapstructure:"keys"`
	Auth string                 `mapstructure:"auth"`
}

type ServerConfigKeys struct {
	Port string
	Host string
}

var ConfigKeys = struct {
	Database   DatabaseConfigKeys
	Encryption string
	Env        string
	Server     ServerConfigKeys
}{
	Database: DatabaseConfigKeys{
		Host:                  "database.host",
		Port:                  "database.port",
		Username:              "database.username",
		Password:              "database.password",
		Name:                  "database.name",
		MaxOpenConnections:    "database.max_open_connections",
		MaxIdleConnections:    "database.max_idle_connections",
		MaxIdleConnectionTime: "database.max_idle_connection_time",
	},
	Encryption: "encryption",
	Server: ServerConfigKeys{
		Port: "server.port",
		Host: "server.host",
	},
	Env: "env",
}
