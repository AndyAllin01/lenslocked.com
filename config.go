package main

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"post"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	/*	if c.Password == "" {
			return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", c.Host,c.Port,c.User,c.Name)
		}
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host,c.Port,c.User,c.Password,c.Name)*/
	//dummy for dev database
	return "postgres://bond:password@localhost/lenslocked_dev?sslmode=disable"

}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "andya",
		Password: "password",
		Name:     "lenslocked_test",
	}
}

type Config struct {
	Port    int    `json:"port"`
	Env     string `json:"env"`
	Pepper  string `json:"pepper"`
	HMACKey string `json:"hmac_key"`
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func DefaultConfig() Config {
	return Config{
		Port:    8080,
		Env:     "dev",
		Pepper:  "secret-random-string-this-project",
		HMACKey: "secret-hmac-key",
	}
}
