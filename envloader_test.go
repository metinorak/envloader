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
			WebsiteURL      string
			FormulaConstant float64
			Database        DBConfig
		}

		// Create mock EnvReader
		mockEnvReader := mocks.NewMockEnvReader(gomock.NewController(t))

		// Set the expected values for the mock
		mockEnvReader.EXPECT().LookupEnv("WEBSITE_URL").Return("https://example.com", true)
		mockEnvReader.EXPECT().LookupEnv("FORMULA_CONSTANT").Return("3.14", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE").Return("", false)
		mockEnvReader.EXPECT().LookupEnv("DATABASE_NAME").Return("db", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE_HOST").Return("localhost", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE_PORT").Return("3306", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE_PASSWORD").Return("password", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE_MAX_CONNS").Return("", false)

		// Create an instance of the EnvLoader
		loader := envLoader{
			envReader: mockEnvReader,
		}

		// Call the Load method
		config := &ConfigModel{}

		err := loader.Load(config)
		if err != nil {
			t.Errorf("Load failed: %s", err)
		}

		expected := &ConfigModel{
			WebsiteURL:      "https://example.com",
			FormulaConstant: 3.14,
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
			WebsiteURL      string   `env:"websiteUrl"`
			FormulaConstant float64  `env:"formulaConstant"`
			Database        DBConfig `env:"database"`
		}

		// Create mock EnvReader
		mockEnvReader := mocks.NewMockEnvReader(gomock.NewController(t))

		// Set the expected values for the mock
		mockEnvReader.EXPECT().LookupEnv("websiteUrl").Return("https://example.com", true)
		mockEnvReader.EXPECT().LookupEnv("formulaConstant").Return("3.14", true)
		mockEnvReader.EXPECT().LookupEnv("database").Return("", false)
		mockEnvReader.EXPECT().LookupEnv("database_dbName").Return("db", true)
		mockEnvReader.EXPECT().LookupEnv("database_dbHost").Return("localhost", true)
		mockEnvReader.EXPECT().LookupEnv("database_dbPort").Return("3306", true)
		mockEnvReader.EXPECT().LookupEnv("database_dbPassword").Return("password", true)
		mockEnvReader.EXPECT().LookupEnv("database_dbMaxConns").Return("", false)

		// Create an instance of the EnvLoader
		loader := envLoader{
			envReader: mockEnvReader,
		}

		// Call the Load method
		config := &ConfigModel{}

		err := loader.Load(config)
		if err != nil {
			t.Errorf("Load failed: %s", err)
		}

		expected := &ConfigModel{
			WebsiteURL:      "https://example.com",
			FormulaConstant: 3.14,
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

	t.Run("TestLoad_WithSpecialEnvNames", func(t *testing.T) {
		type DBConfig struct {
			Name1     string
			Host1     string
			Port1     int
			Password1 string
			MaxConns1 int
		}

		type ConfigModel struct {
			Website1URL string
			Database1   DBConfig
		}

		// Create mock EnvReader
		mockEnvReader := mocks.NewMockEnvReader(gomock.NewController(t))

		// Set the expected values for the mock
		mockEnvReader.EXPECT().LookupEnv("WEBSITE1_URL").Return("https://example.com", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE1").Return("", false)
		mockEnvReader.EXPECT().LookupEnv("DATABASE1_NAME1").Return("db", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE1_HOST1").Return("localhost", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE1_PORT1").Return("3306", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE1_PASSWORD1").Return("password", true)
		mockEnvReader.EXPECT().LookupEnv("DATABASE1_MAX_CONNS1").Return("15", true)

		// Create an instance of the EnvLoader
		loader := envLoader{
			envReader: mockEnvReader,
		}

		// Call the Load method
		config := &ConfigModel{}

		err := loader.Load(config)
		assert.NoError(t, err)

		expected := &ConfigModel{
			Website1URL: "https://example.com",
			Database1: DBConfig{
				Name1:     "db",
				Host1:     "localhost",
				Port1:     3306,
				Password1: "password",
				MaxConns1: 15,
			},
		}

		assert.Equal(t, expected, config)
	})

	t.Run("TestLoad_WithDefaultValues", func(t *testing.T) {
		type DBConfig struct {
			Name     string `env:"dbName" default:"db"`
			Host     string `env:"dbHost" default:"localhost"`
			Port     int    `env:"dbPort" default:"3306"`
			Password string `env:"dbPassword" default:"password"`
			MaxConns int    `env:"dbMaxConns" default:"15"`
		}

		type ConfigModel struct {
			WebsiteURL       string             `env:"websiteUrl" default:"https://example.com"`
			FormulaConstants map[string]float64 `env:"formulaConstants" default:"pi:3.14,e:2.71828"`
			UserRoles        []string           `env:"userRoles" default:"admin,editor,author"`
			Database         DBConfig           `env:"database"`
		}

		// Create mock EnvReader
		mockEnvReader := mocks.NewMockEnvReader(gomock.NewController(t))

		// Set the expected values for the mock
		mockEnvReader.EXPECT().LookupEnv("websiteUrl").Return("", false)
		mockEnvReader.EXPECT().LookupEnv("formulaConstants").Return("", false)
		mockEnvReader.EXPECT().LookupEnv("userRoles").Return("", false)
		mockEnvReader.EXPECT().LookupEnv("database").Return("", false)
		mockEnvReader.EXPECT().LookupEnv("database_dbName").Return("", false)
		mockEnvReader.EXPECT().LookupEnv("database_dbHost").Return("", false)
		mockEnvReader.EXPECT().LookupEnv("database_dbPort").Return("", false)
		mockEnvReader.EXPECT().LookupEnv("database_dbPassword").Return("", false)
		mockEnvReader.EXPECT().LookupEnv("database_dbMaxConns").Return("", false)

		// Create an instance of the EnvLoader
		loader := envLoader{
			envReader: mockEnvReader,
		}

		// Call the Load method
		config := &ConfigModel{}

		err := loader.Load(config)
		assert.NoError(t, err)

		expected := &ConfigModel{
			WebsiteURL: "https://example.com",
			FormulaConstants: map[string]float64{
				"pi": 3.14,
				"e":  2.71828,
			},
			UserRoles: []string{"admin", "editor", "author"},
			Database: DBConfig{
				Name:     "db",
				Host:     "localhost",
				Port:     3306,
				Password: "password",
				MaxConns: 15,
			},
		}

		assert.Equal(t, expected, config)
	})

	t.Run("TestLoad_WithSlices", func(t *testing.T) {
		type ConfigModel struct {
			Proxies []string
		}

		// Create mock EnvReader
		mockEnvReader := mocks.NewMockEnvReader(gomock.NewController(t))

		// Set the expected values for the mock
		mockEnvReader.EXPECT().LookupEnv("PROXIES").Return("https://example.com,https://example2.com", true)

		// Create an instance of the EnvLoader
		loader := envLoader{
			envReader: mockEnvReader,
		}

		// Call the Load method
		config := &ConfigModel{}

		err := loader.Load(config)
		assert.NoError(t, err)

		expected := &ConfigModel{
			Proxies: []string{"https://example.com", "https://example2.com"},
		}

		assert.Equal(t, expected, config)
	})

	t.Run("TestLoad_WithMaps", func(t *testing.T) {
		type ConfigModel struct {
			FormulaFactors map[string]float64
		}

		// Create mock EnvReader
		mockEnvReader := mocks.NewMockEnvReader(gomock.NewController(t))

		// Set the expected values for the mock
		mockEnvReader.EXPECT().LookupEnv("FORMULA_FACTORS").Return("pi:3.14,e:2.71828", true)

		// Create an instance of the EnvLoader
		loader := envLoader{
			envReader: mockEnvReader,
		}

		// Call the Load method
		config := &ConfigModel{}

		err := loader.Load(config)
		assert.NoError(t, err)

		expected := &ConfigModel{
			FormulaFactors: map[string]float64{
				"pi": 3.14,
				"e":  2.71828,
			},
		}

		assert.Equal(t, expected, config)
	})

	t.Run("TestLoad_WithRequiredFields", func(t *testing.T) {
		type ConfigModel struct {
			WebsiteURL string `required:"true"`
		}

		// Create mock EnvReader
		mockEnvReader := mocks.NewMockEnvReader(gomock.NewController(t))

		// Set the expected values for the mock
		mockEnvReader.EXPECT().LookupEnv("WEBSITE_URL").Return("", false)

		// Create an instance of the EnvLoader
		loader := envLoader{
			envReader: mockEnvReader,
		}

		// Call the Load method
		config := &ConfigModel{}

		err := loader.Load(config)
		assert.Error(t, err)
	})

	t.Run("TestLoad_WhenEnvValueIsNotValid", func(t *testing.T) {
		type ConfigModel struct {
			FormulaConstant float64
		}

		// Create mock EnvReader
		mockEnvReader := mocks.NewMockEnvReader(gomock.NewController(t))

		// Set the expected values for the mock
		mockEnvReader.EXPECT().LookupEnv("FORMULA_CONSTANT").Return("3.14.15", true)

		// Create an instance of the EnvLoader
		loader := envLoader{
			envReader: mockEnvReader,
		}

		// Call the Load method
		config := &ConfigModel{}

		err := loader.Load(config)
		assert.Error(t, err)
	})

	t.Run("TestLoad_WhenMapValueIsNotValid", func(t *testing.T) {
		type ConfigModel struct {
			FormulaFactors map[string]float64
		}

		// Create mock EnvReader
		mockEnvReader := mocks.NewMockEnvReader(gomock.NewController(t))

		// Set the expected values for the mock
		mockEnvReader.EXPECT().LookupEnv("FORMULA_FACTORS").Return("pi:abc,e:2.71828,phi:1.618", true)

		// Create an instance of the EnvLoader
		loader := envLoader{
			envReader: mockEnvReader,
		}

		// Call the Load method
		config := &ConfigModel{}

		err := loader.Load(config)
		assert.Error(t, err)
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
	mockEnvReader.EXPECT().LookupEnv("DATABASE_NAME").Return("db", true).AnyTimes()
	mockEnvReader.EXPECT().LookupEnv("DATABASE_HOST").Return("localhost", true).AnyTimes()
	mockEnvReader.EXPECT().LookupEnv("DATABASE_PORT").Return("3306", true).AnyTimes()
	mockEnvReader.EXPECT().LookupEnv("DATABASE_PASSWORD").Return("password", true).AnyTimes()
	mockEnvReader.EXPECT().LookupEnv("DATABASE_MAX_CONNS").Return("", false).AnyTimes()

	for i := 0; i < b.N; i++ {
		// Create an instance of the EnvLoader with default options
		loader := envLoader{
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
