package postgres

import "time"

// Option -.
type Option func(*Postgres)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}

// MinPoolSize -.
func MinPoolSize(size int) Option {
	return func(c *Postgres) {
		c.minPoolSize = size
	}
}

// IdleTimeoutMinutes -.
func IdleTimeoutMinutes(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.idleTimeoutMinutes = timeout * time.Minute
	}
}
