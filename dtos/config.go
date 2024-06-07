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

type DatabaseBackupConfigKeys struct {
	Bucket   string
	Location string
	Region   string
}

type EncryptionConfigKeys struct {
	Key     string `mapstructure:"key"`
	Version string `mapstructure:"version"`
}

type Encryption struct {
	Keys []EncryptionConfigKeys `mapstructure:"keys"`
	Auth string                 `mapstructure:"auth"`
}

type CaptchaConfigKeys struct {
	Secret string
}

type ServerConfigKeys struct {
	Port string
	Host string
}

var ConfigKeys = struct {
	Database       DatabaseConfigKeys
	Encryption     string
	Env            string
	Server         ServerConfigKeys
	DatabaseBackup DatabaseBackupConfigKeys
	Captcha        CaptchaConfigKeys
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
	DatabaseBackup: DatabaseBackupConfigKeys{
		Bucket:   "database_backup.bucket",
		Location: "database_backup.location",
		Region:   "database_backup.region",
	},
	Captcha: CaptchaConfigKeys{
		Secret: "captcha.secret",
	},
}
