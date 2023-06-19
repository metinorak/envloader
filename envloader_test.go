package envloader

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/metinorak/envloader/mocks"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("TestLoad_WithoutTags", func(t *testing.T) {
		type DBConfig struct {
			Name     string
			Host     string
			Port     int
			Password string
			MaxConns int
		}

		type ConfigModel struct {
			WebsiteURL string
			Database   DBConfig
		}

		// Create mock EnvReader
		mockEnvReader := mocks.NewMockEnvReader(gomock.NewController(t))

		// Set the expected values for the mock
		mockEnvReader.EXPECT().LookupEnv("WEBSITE_URL").Return("https://example.com", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE").Return("", false)
		mockEnvReader.EXPECT().LookupEnv("DATABASE.NAME").Return("db", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE.HOST").Return("localhost", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE.PORT").Return("3306", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE.PASSWORD").Return("password", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE.MAX_CONNS").Return("", false)

		// Create an instance of the EnvLoader with default options
		loader := envLoader{
			options:   DefaultOptions(),
			envReader: mockEnvReader,
		}

		// Call the Load method
		config := &ConfigModel{}

		err := loader.Load(config)
		if err != nil {
			t.Errorf("Load failed: %s", err)
		}

		expected := &ConfigModel{
			WebsiteURL: "https://example.com",
			Database: DBConfig{
				Name:     "db",
				Host:     "localhost",
				Port:     3306,
				Password: "password",
				MaxConns: 0,
			},
		}

		assert.Equal(t, expected, config)
	})

	t.Run("TestLoad_WithTags", func(t *testing.T) {

		type DBConfig struct {
			Name     string `env:"dbName"`
			Host     string `env:"dbHost"`
			Port     int    `env:"dbPort"`
			Password string `env:"dbPassword"`
			MaxConns int    `env:"dbMaxConns"`
		}

		type ConfigModel struct {
			WebsiteURL string   `env:"websiteUrl"`
			Database   DBConfig `env:"database"`
		}

		// Create mock EnvReader
		mockEnvReader := mocks.NewMockEnvReader(gomock.NewController(t))

		// Set the expected values for the mock
		mockEnvReader.EXPECT().LookupEnv("websiteUrl").Return("https://example.com", true)
		mockEnvReader.EXPECT().LookupEnv("database").Return("", false)
		mockEnvReader.EXPECT().LookupEnv("database.dbName").Return("db", true)
		mockEnvReader.EXPECT().LookupEnv("database.dbHost").Return("localhost", true)
		mockEnvReader.EXPECT().LookupEnv("database.dbPort").Return("3306", true)
		mockEnvReader.EXPECT().LookupEnv("database.dbPassword").Return("password", true)
		mockEnvReader.EXPECT().LookupEnv("database.dbMaxConns").Return("", false)

		// Create an instance of the EnvLoader with default options
		loader := envLoader{
			options:   DefaultOptions(),
			envReader: mockEnvReader,
		}

		// Call the Load method
		config := &ConfigModel{}

		err := loader.Load(config)
		if err != nil {
			t.Errorf("Load failed: %s", err)
		}

		expected := &ConfigModel{
			WebsiteURL: "https://example.com",
			Database: DBConfig{
				Name:     "db",
				Host:     "localhost",
				Port:     3306,
				Password: "password",
				MaxConns: 0,
			},
		}

		assert.Equal(t, expected, config)
	})

	t.Run("TestLoad_WithCustomDelimiter", func(t *testing.T) {
		type DBConfig struct {
			Name     string
			Host     string
			Port     int
			Password string
			MaxConns int
		}

		type ConfigModel struct {
			WebsiteURL string
			Database   DBConfig
		}

		// Create mock EnvReader
		mockEnvReader := mocks.NewMockEnvReader(gomock.NewController(t))

		// Set the expected values for the mock
		mockEnvReader.EXPECT().LookupEnv("WEBSITE_URL").Return("https://example.com", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE").Return("", false)
		mockEnvReader.EXPECT().LookupEnv("DATABASE-NAME").Return("db", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE-HOST").Return("localhost", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE-PORT").Return("3306", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE-PASSWORD").Return("password", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE-MAX_CONNS").Return("5", true)

		// Create an instance of the EnvLoader with default options
		loader := envLoader{
			options: Options{
				EnvFieldDelimiter: "-",
			},
			envReader: mockEnvReader,
		}

		// Call the Load method
		config := &ConfigModel{}

		err := loader.Load(config)
		if err != nil {
			t.Errorf("Load failed: %s", err)
		}

		expected := &ConfigModel{
			WebsiteURL: "https://example.com",
			Database: DBConfig{
				Name:     "db",
				Host:     "localhost",
				Port:     3306,
				Password: "password",
				MaxConns: 5,
			},
		}

		assert.Equal(t, expected, config)
	})
}

func BenchmarkLoad(b *testing.B) {
	type DBConfig struct {
		Name     string
		Host     string
		Port     int
		Password string
		MaxConns int
	}

	type ConfigModel struct {
		WebsiteURL string
		Database   DBConfig
	}

	// Create mock EnvReader
	mockEnvReader := mocks.NewMockEnvReader(gomock.NewController(b))

	// Set the expected values for the mock
	mockEnvReader.EXPECT().LookupEnv("WEBSITE_URL").Return("https://example.com", true).AnyTimes()
	mockEnvReader.EXPECT().LookupEnv("DATABASE").Return("", false).AnyTimes()
	mockEnvReader.EXPECT().LookupEnv("DATABASE.NAME").Return("db", true).AnyTimes()
	mockEnvReader.EXPECT().LookupEnv("DATABASE.HOST").Return("localhost", true).AnyTimes()
	mockEnvReader.EXPECT().LookupEnv("DATABASE.PORT").Return("3306", true).AnyTimes()
	mockEnvReader.EXPECT().LookupEnv("DATABASE.PASSWORD").Return("password", true).AnyTimes()
	mockEnvReader.EXPECT().LookupEnv("DATABASE.MAX_CONNS").Return("", false).AnyTimes()

	for i := 0; i < b.N; i++ {
		// Create an instance of the EnvLoader with default options
		loader := envLoader{
			options:   DefaultOptions(),
			envReader: mockEnvReader,
		}

		// Call the Load method
		config := &ConfigModel{}

		err := loader.Load(config)
		if err != nil {
			b.Errorf("Load failed: %s", err)
		}
	}
}
