package config

type Config struct {
	App    App
	DB     DB
	JwtKey []byte
}

type App struct {
	Host string
	Port string
}

type DB struct {
	Name     string
	Host     string
	Port     string
	User     string
	Password string
	SslMode  string
}
