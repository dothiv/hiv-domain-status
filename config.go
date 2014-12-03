package hivdomainstatus

import "fmt"

type Config struct {
	Database   struct {
		Host     string
		Name     string
		User     string
		Password string
		Sslmode  string
	}
}

func (c *Config) DSN() (dsn string) {
	dsn = fmt.Sprintf("user=%s dbname=%s sslmode=%s", c.Database.User, c.Database.Name, c.Database.Sslmode)
	if len(c.Database.Host) > 0 {
		dsn = dsn + " host=" + c.Database.Host
	}
	if len(c.Database.Password) > 0 {
		dsn = dsn + " password=" + c.Database.Password
	}
	return
}

func NewDefaultConfig() (c *Config) {
	c = new(Config)
	c.Database.Sslmode = "disable"
	return
}
