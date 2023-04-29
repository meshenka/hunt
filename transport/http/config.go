package http

type config struct {
	databaseDSN string
	host        string
}

type Option func(*config) error

func WithDatabaseDSN(dsn string) Option {
	return func(c *config) error {
		c.databaseDSN = dsn
		return nil
	}
}

func WithHost(host string) Option {
	return func(c *config) error {
		c.host = host
		return nil
	}
}
