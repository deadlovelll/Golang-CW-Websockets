package databaseconfig

// DatabaseConfig holds the database connection details.
type DatabaseConfig struct {
	Type     string
	User     string
	Name     string
	Host     string
	Password string
	SslMode  string
}