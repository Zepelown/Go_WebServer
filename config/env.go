package config

type EnvConfig struct {
	DbUrl              string `env:"DB_URL"`
	FrontUrl           string `env:"FRONT_URL"`
	ServerPortUrl      string `env:"SERVER_PORT_URL"`
	JwtSecretKey       string `env:"JWT_SECRET_KEY"`
	JwtAcessTokenTTL   string `env:"ACCESS_TOKEN_TTL"`
	JwtRefreshTokenTTL string `env:"REFRESH_TOKEN_TTL"`
	JwtIssuer          string `env:"JWT_ISSUER"`
}
