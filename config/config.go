package config

var Conf *Config

const (
	ConfigKeyDBUser = "db_user"
	ConfigKeyDBPass = "db_password"
	ConfigKeyDBHost = "db_host"
	ConfigKeyDBName = "db_name"
)

var configMap = map[string]string{
	ConfigKeyDBUser: "root",
	ConfigKeyDBPass: "root123",
	ConfigKeyDBHost: "localhost",
	ConfigKeyDBName: "payment",
}

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBName     string
}

func init() {
	Conf = &Config{
		DBUser:     configMap[ConfigKeyDBUser],
		DBPassword: configMap[ConfigKeyDBPass],
		DBHost:     configMap[ConfigKeyDBHost],
		DBName:     configMap[ConfigKeyDBName],
	}
}
