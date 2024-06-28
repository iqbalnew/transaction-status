package config

import (
	"encoding/json"
	"fmt"
	"strconv"

	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/utils"
	"github.com/joho/godotenv"
)

type Config struct {
	// Listen address is an array of IP addresses and port combinations.
	// Listen address is an array so that this service can listen to many interfaces at once.
	// You can use this value for example: []string{"192.168.1.12:80", "25.49.25.73:80"} to listen to
	// listen to interfaces with IP address of 192.168.1.12 and 25.49.25.73, both on port 80.
	ListenAddress string `config:"LISTEN_ADDRESS"`

	CorsAllowedHeaders []string `config:"CORS_ALLOWED_HEADERS"`
	CorsAllowedMethods []string `config:"CORS_ALLOWED_METHODS"`
	CorsAllowedOrigins []string `config:"CORS_ALLOWED_ORIGINS"`
	ExposedHeaders     []string `config:"EXPOSED_HEADERS"`

	AppName        string `config:"APP_NAME"`
	AppPort        int    `config:"APP_PORT"`
	ApiServicePath string `config:"API_SERVICE_PATH"`

	SwaggerPath string `config:"SWAGGER_PATH"`

	AuthService string `config:"AUTH_SERVICE"`

	DbHost     string `config:"DB_HOST"`
	DbUser     string `config:"DB_USER"`
	DbPassword string `config:"DB_PASSWORD"`
	DbName     string `config:"DB_NAME"`
	DbPort     string `config:"DB_PORT"`
	DbSslmode  string `config:"DB_SSLMODE"`
	DbTimezone string `config:"DB_TIMEZONE"`
	DbMaxRetry string `config:"DB_MAX_RETRY"`
	DbTimeout  string `config:"DB_TIMEOUT"`

	FluentbitHost string `config:"FLUENTBIT_HOST"`
	FluentbitPort string `config:"FLUENTBIT_PORT"`
	LoggerTag     string `config:"LOGGER_TAG"`
	LoggerOutput  string `config:"LOGGER_OUTPUT"`
	LoggerLevel   string `config:"LOGGER_LEVEL"`

	NotificationModuleId   uint64 `config:"NOTIFICATION_MODULE_ID"`
	NotificationModuleName string `config:"NOTIFICATION_MODULE_NAME"`

	AmqpHost            string `config:"AMQP_HOST"`
	AmqpUser            string `config:"AMQP_USER"`
	AmqpPassword        string `config:"AMQP_PASSWORD"`
	AmqpAutoReconnect   string `config:"AMQP_AUTO_RECONNECT"`
	AmqpReconnectDelay  string `config:"AMQP_RECONNECT_DELAY"`
	AmqpQueueConsumer   string `config:"AMQP_QUEUE_CONSUMER"`
	AmqpQueueReportPush string `config:"AMQP_QUEUE_REPORT_PUSH"`

	TimeLocation string `config:"TIME_LOCATION"`
}

func InitConfig() *Config {
	// Todo: add env checker

	godotenv.Load(".env")
	appName := utils.GetEnv("APP_NAME", "")
	if appName == "" {
		appName = utils.GetEnv("ELASTIC_APM_SERVICE_NAME", "")
	}
	appPort, _ := strconv.Atoi(utils.GetEnv("APP_PORT", "9090"))

	notificationModuleId, _ := strconv.ParseUint(utils.GetEnv("NOTIFICATION_MODULE_ID", "236"), 10, 64)

	return &Config{
		ListenAddress: fmt.Sprintf("%s:%s", utils.GetEnv("HOST", "0.0.0.0"), utils.GetEnv("PORT", ":80")),
		CorsAllowedHeaders: []string{
			"Connection", "User-Agent", "Referer",
			"Accept", "Accept-Language", "Content-Type",
			"Content-Language", "Content-Disposition", "Origin",
			"Content-Length", "Authorization", "ResponseType",
			"X-Requested-With", "X-Forwarded-For", "Grpc-Metadata-Signature",
		},
		CorsAllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
		CorsAllowedOrigins: []string{"*"},

		AppName:        appName,
		AppPort:        appPort,
		ApiServicePath: utils.GetEnv("API_SERVICE_PATH", "/account_receivable.service.v1.AccountReceivableService/"),

		DbHost:     utils.GetEnv("DB_HOST", ""),
		DbUser:     utils.GetEnv("DB_USER", ""),
		DbPassword: utils.GetEnv("DB_PASSWORD", ""),
		DbName:     utils.GetEnv("DB_NAME", ""),
		DbPort:     utils.GetEnv("DB_PORT", ""),
		DbSslmode:  utils.GetEnv("DB_SSLMODE", ""),
		DbTimezone: utils.GetEnv("DB_TIMEZONE", ""),
		DbMaxRetry: utils.GetEnv("DB_MAX_RETRY", "3"),
		DbTimeout:  utils.GetEnv("DB_TIMEOUT", "120"),

		SwaggerPath: utils.GetEnv("SWAGGER_PATH", "www/"),

		AuthService: utils.GetEnv("AUTH_SERVICE", ":9091"),

		FluentbitHost: utils.GetEnv("FLUENTBIT_HOST", "0.0.0.0"),
		FluentbitPort: utils.GetEnv("FLUENTBIT_PORT", "24223"),
		LoggerTag:     utils.GetEnv("LOGGER_TAG", "addons.account.receivable.service.dev"),
		LoggerOutput:  utils.GetEnv("LOGGER_OUTPUT", "elastic"),
		LoggerLevel:   utils.GetEnv("LOGGER_LEVEL", "debug"),

		NotificationModuleId:   notificationModuleId,
		NotificationModuleName: utils.GetEnv("NOTIFICATION_MODULE_NAME", "Account Receivable"),

		TimeLocation: utils.GetEnv("TIME_LOCATION", "Asia/Jakarta"),
	}
}

func (c *Config) AsString() string {
	data, _ := json.Marshal(c)
	return string(data)
}
