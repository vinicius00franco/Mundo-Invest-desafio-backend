package config

import (
	"os"
	"strconv"
	"time"
)

// Config centraliza todas as configurações da aplicação
type Config struct {
	// Database
	Database DatabaseConfig

	// HTTP
	HTTP HTTPConfig

	// Timeouts
	Timeouts TimeoutConfig

	// Business
	Business BusinessConfig

	// Pipefy
	Pipefy PipefyConfig
}

// DatabaseConfig configurações do banco de dados
type DatabaseConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// HTTPConfig configurações do servidor HTTP
type HTTPConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// TimeoutConfig configurações de timeout
type TimeoutConfig struct {
	Default     time.Duration
	Database    time.Duration
	ExternalAPI time.Duration
}

// BusinessConfig configurações de negócio
type BusinessConfig struct {
	LimitePrioridadeAlta     float64
	MaxNomeLength            int
	MaxEmailLength           int
	MaxTipoSolicitacaoLength int
	MinValorPatrimonio       float64
	MaxValorPatrimonio       float64
}

// PipefyConfig configurações do Pipefy
type PipefyConfig struct {
	APIToken string
	APIURL   string
	PipeID   string
}

// Load carrega as configurações de environment variables com valores padrão
func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 300*time.Second),
			ConnMaxIdleTime: getDurationEnv("DB_CONN_MAX_IDLE_TIME", 60*time.Second),
		},
		HTTP: HTTPConfig{
			Port:         getStringEnv("HTTP_PORT", "8080"),
			ReadTimeout:  getDurationEnv("HTTP_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("HTTP_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("HTTP_IDLE_TIMEOUT", 60*time.Second),
		},
		Timeouts: TimeoutConfig{
			Default:     getDurationEnv("TIMEOUT_DEFAULT", 30*time.Second),
			Database:    getDurationEnv("TIMEOUT_DATABASE", 10*time.Second),
			ExternalAPI: getDurationEnv("TIMEOUT_EXTERNAL_API", 15*time.Second),
		},
		Business: BusinessConfig{
			LimitePrioridadeAlta:     getFloatEnv("BUSINESS_LIMITE_PRIORIDADE_ALTA", 200000.00),
			MaxNomeLength:            getIntEnv("BUSINESS_MAX_NOME_LENGTH", 255),
			MaxEmailLength:           getIntEnv("BUSINESS_MAX_EMAIL_LENGTH", 255),
			MaxTipoSolicitacaoLength: getIntEnv("BUSINESS_MAX_TIPO_SOLICITACAO_LENGTH", 50),
			MinValorPatrimonio:       getFloatEnv("BUSINESS_MIN_VALOR_PATRIMONIO", 0.00),
			MaxValorPatrimonio:       getFloatEnv("BUSINESS_MAX_VALOR_PATRIMONIO", 999999999.99),
		},
		Pipefy: PipefyConfig{
			APIToken: getStringEnv("PIPEFY_API_TOKEN", ""),
			APIURL:   getStringEnv("PIPEFY_API_URL", "https://api.pipefy.com/graphql"),
			PipeID:   getStringEnv("PIPEFY_PIPE_ID", ""),
		},
	}
}

// getStringEnv retorna o valor da environment variable ou o valor padrão
func getStringEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv retorna o valor da environment variable como int ou o valor padrão
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getFloatEnv retorna o valor da environment variable como float64 ou o valor padrão
func getFloatEnv(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

// getDurationEnv retorna o valor da environment variable como duration ou o valor padrão
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
