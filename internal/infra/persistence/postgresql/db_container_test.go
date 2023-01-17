package postgresql_test

import (
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

// Config defines the parameters for the dockerized PostgresSQL used for testing
type Config struct {
	Scheme         string
	Driver         string
	Database       string
	User           string
	Password       string
	Host           string
	Port           string
	SSLMode        string
	BinaryParams   string
	MigrationsPath string
}

const defaultDBMigrationsPath = "../migrations"

// NewConfig is a constructor
func NewConfig() Config {
	return Config{
		Scheme:         "postgres",
		Driver:         "postgres",
		Database:       "test_db",
		User:           "postgres",
		Password:       "",
		Host:           "localhost",
		Port:           "54322",
		SSLMode:        "disable",
		BinaryParams:   "yes",
		MigrationsPath: defaultDBMigrationsPath,
	}
}

// DatabaseURL returns the database URL already built with its param values.
func (c Config) DatabaseURL() string {
	u := url.URL{
		Scheme:   c.Scheme,
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(c.Host, c.Port),
		Path:     c.Database,
		RawQuery: fmt.Sprintf("sslmode=%s&binary_parameters=%s", c.SSLMode, c.BinaryParams),
	}

	return u.String()
}

var defaultTestContainerPort = "5432"

// CreateContainerizedPgDB creates an ephemeral PostgreSQL container instance for testing purposes
func CreateContainerizedPgDB(t *testing.T, config Config) {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	options := &dockertest.RunOptions{
		Repository: "postgres",
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", config.Password),
			fmt.Sprintf("POSTGRES_USER=%s", config.User),
			"POSTGRES_HOST_AUTH_METHOD=trust", // only for testing purposes!!!
			fmt.Sprintf("POSTGRES_DB=%s", config.Database),
		},
		ExposedPorts: []string{defaultTestContainerPort},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port(defaultTestContainerPort): {
				{
					HostIP:   "0.0.0.0",
					HostPort: config.Port,
				},
			},
		},
	}
	container, err := pool.RunWithOptions(options)
	require.NoError(t, err)
	checkFunc := func() error { return checkDBConn(config) }
	require.NoError(t, pool.Retry(checkFunc))
	t.Cleanup(func() {
		if err := container.Close(); err != nil {
			t.Fatal(err)
		}
	})
}

func checkDBConn(config Config) error {
	conn, err := sql.Open(config.Driver, config.DatabaseURL())
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Ping()
}
