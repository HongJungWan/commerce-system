package configs

type Config struct {
	DBHost     string `mapstructure:"mariadb_host"`
	DBUserName string `mapstructure:"mariadb_user"`
	DBPassword string `mapstructure:"mariadb_password"`
	DBName     string `mapstructure:"mariadb_db"`
	DBPort     string `mapstructure:"mariadb_port"`

	Environment string `toml:"ENVIRONMENT"`

	Host     string   `toml:"HOST"`
	Scheme   []string `toml:"SCHEME"`
	Version  string   `toml:"VERSION"`
	BasePath string   `toml:"BASEPATH"`
	Title    string   `toml:"TITLE"`
}
