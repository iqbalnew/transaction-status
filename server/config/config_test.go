package config

import (
	"context"
	"os"
	"testing"

	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/utils"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *ConfigTestSuite) SetupSuite() {
	s.ctx = context.Background()
}

func TestInitConfig(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) TestConfig_InitConfig() {
	type expectation struct {
		cfg *Config
	}

	tests := map[string]struct {
		expected expectation
	}{
		"Success": {
			expected: expectation{
				cfg: &Config{},
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			cfg := InitConfig()

			if cfg == nil {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.cfg, cfg)
			}
		})
	}
}

func (s *ConfigTestSuite) TestConfig_AsString() {
	type expectation struct {
		out string
	}

	tests := map[string]struct {
		cfg      *Config
		expected expectation
	}{
		"Success": {
			cfg: &Config{
				AppName: "Account Receivable",
			},
			expected: expectation{
				out: `{"ListenAddress":"","CorsAllowedHeaders":null,"CorsAllowedMethods":null,"CorsAllowedOrigins":null,"JwtSecret":"","JwtDuration":"","SignatureSecretKey":"","AppName":"Account Receivable","SwaggerPath":"","TaskService":"","AuthService":"","AccountService":"","BeneficiaryAccountService":"","SystemService":"","WorkflowService":"","TransactionService":"","NotificationService":"","QueueService":"","MenuService":"","CompanyService":"","CutOffService":"","FileService":"","DbHost":"","DbUser":"","DbPassword":"","DbName":"","DbPort":"","DbSslmode":"","DbTimezone":"","DbMaxRetry":"","DbTimeout":"","FluentbitHost":"","FluentbitPort":"","LoggerTag":"","LoggerOutput":"","LoggerLevel":""}`,
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			configJson := tt.cfg.AsString()

			if configJson != tt.expected.out {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, configJson)
			}
		})
	}
}

func (s *ConfigTestSuite) TestConfig_getEnv() {
	type expectation struct {
		out string
	}

	tests := map[string]struct {
		envTemp  string
		expected expectation
	}{
		"Success": {
			envTemp: "UNIT_TEST",
			expected: expectation{
				out: `UNIT_TEST`,
			},
		},
		"Fallback": {
			envTemp: "NOT_FOUND",
			expected: expectation{
				out: `not_found`,
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			if tt.envTemp != "NOT_FOUND" {
				os.Setenv(tt.envTemp, tt.envTemp)
			}
			envData := utils.GetEnv(tt.envTemp, "not_found")

			if envData != tt.expected.out {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, envData)
			}
		})
	}
}
