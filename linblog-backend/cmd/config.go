package cmd

type Config struct {
	Port  string
	Mysql MysqlConfig `yaml:"MYSQL"`
	Redis RedisConfig `yaml:"REDIS"`
}

type MysqlConfig struct {
	Username string
	Password string
	Address  string
	DBName   string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}
