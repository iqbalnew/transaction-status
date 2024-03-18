package servicelogger

import (
	"os"
	"strconv"
	"strings"

	"github.com/evalphobia/logrus_fluent"
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
)

type customLogger struct {
	ecslogrus   *ecslogrus.Formatter
	serviceName string
	hostName    string
	formatter   logrus.Formatter
}

type LoggerConfig struct {
	ServiceName   string
	LogOutput     string
	LogLevel      string
	FluentbitHost string
	FluentbitPort string
}

type AddonsLogrus struct {
	GrpcMetadataKey string
	*logrus.Logger
}

func newCustomFormatter(serviceName string, hostname string) *customLogger {

	return &customLogger{
		ecslogrus: &ecslogrus.Formatter{
			DataKey:     "data_details",
			PrettyPrint: true,
		},
		serviceName: serviceName,
		hostName:    hostname,
		formatter:   logrus.StandardLogger().Formatter,
	}
}

func (l *customLogger) Format(entry *logrus.Entry) ([]byte, error) {
	// set ecs format
	entry.Data["service_name"] = l.serviceName
	entry.Data["host_name"] = l.hostName

	return l.ecslogrus.Format(entry)
}

func New(loggerConfig *LoggerConfig) *AddonsLogrus {
	name := strings.ToLower(loggerConfig.ServiceName)
	name = strings.ReplaceAll(name, " ", "_")

	port, _ := strconv.Atoi(loggerConfig.FluentbitPort)

	hostname, _ := os.Hostname()

	log := logrus.New()

	log.SetFormatter(newCustomFormatter(name, hostname))

	// Add the GlobalKeyHook to the logrus.Logger
	log.Hooks.Add(&GlobalKeyHook{
		keys: logrus.Fields{
			"service_name": name,
			"host_name":    hostname,
		},
		esl: &ecslogrus.Formatter{
			DataKey: "data_details",
		},
	})

	logLevel, levelError := logrus.ParseLevel(loggerConfig.LogLevel)
	if levelError != nil {
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(logLevel)
	}

	if strings.ToLower(loggerConfig.LogOutput) == "elastic" {

		hook, err := logrus_fluent.NewWithConfig(logrus_fluent.Config{
			Port: port,
			Host: loggerConfig.FluentbitHost,
		})
		if err != nil {
			panic(err)
		}

		logrus.Info("no error connection for fluentbit")

		logLevels := []logrus.Level{}
		logLevel, levelError := logrus.ParseLevel(loggerConfig.LogLevel)
		if levelError != nil {
			logLevels = append(logLevels, logrus.AllLevels...)
		} else {
			logLevels = append(logLevels, logLevel)
		}
		hook.SetLevels(logLevels)
		hook.SetTag(name)

		hook.SetMessageField("message")

		// ignore field
		// hook.AddIgnore("context")

		// filter func
		hook.AddFilter("error", logrus_fluent.FilterError)

		log.AddHook(hook)

	}

	log.ReportCaller = true

	return &AddonsLogrus{"grpc_metadata", log}
}

type errorObject struct {
	Message string `json:"message,omitempty"`
}

type GlobalKeyHook struct {
	keys logrus.Fields
	esl  *ecslogrus.Formatter
}

func (h *GlobalKeyHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *GlobalKeyHook) Fire(entry *logrus.Entry) error {
	// now := time.Now().Format("2006-01-02T15:04:05.000Z0700")

	for k, v := range h.keys {
		entry.Data[k] = v
	}

	datahint := len(entry.Data)
	if h.esl.DataKey != "" {
		datahint = 2
	}
	data := make(logrus.Fields, datahint)
	if len(entry.Data) > 0 {
		extraData := data
		if h.esl.DataKey != "" {
			extraData = make(logrus.Fields, len(entry.Data))
		}
		for k, v := range entry.Data {
			switch k {
			case logrus.ErrorKey:
				err, ok := v.(error)
				if ok {
					data["error"] = errorObject{
						Message: err.Error(),
					}
					break
				}
				fallthrough // error has unexpected type
			default:
				if k != "service_name" && k != "data_tag" {
					delete(entry.Data, k)
				}
				extraData[k] = v
			}
		}
		if h.esl.DataKey != "" && len(extraData) > 0 {
			data[h.esl.DataKey] = extraData
		}
	}
	if entry.HasCaller() {
		// Logrus has a single configurable field (logrus.FieldKeyFile)
		// for storing a combined filename and line number, but we want
		// to split them apart into two fields. Remove the event's Caller
		// field, and encode the ECS fields explicitly.
		var funcVal, fileVal string
		var lineVal int
		if h.esl.CallerPrettyfier != nil {
			var fileLineVal string
			funcVal, fileLineVal = h.esl.CallerPrettyfier(entry.Caller)
			if sep := strings.IndexRune(fileLineVal, ':'); sep != -1 {
				fileVal = fileLineVal[:sep]
				lineVal, _ = strconv.Atoi(fileLineVal[sep+1:])
			} else {
				fileVal = fileLineVal
				lineVal = 0
			}
		} else {
			funcVal = entry.Caller.Function
			fileVal = entry.Caller.File
			lineVal = entry.Caller.Line
		}
		entry.Caller = nil

		if funcVal != "" {
			data["log.origin.function"] = funcVal
		}
		if fileVal != "" {
			data["log.origin.file.name"] = fileVal
		}
		if lineVal > 0 {
			data["log.origin.file.line"] = lineVal
		}
	}

	for k, v := range data {
		entry.Data[k] = v
	}

	return nil
}

func getTagName(name string, logName string) map[string]string {
	var envName string
	if strings.Contains(name, "dev") {
		envName = "dev"
	} else if strings.Contains(name, "staging") {
		envName = "staging"
	} else if strings.Contains(name, "prestaging") {
		envName = "prestaging"
	} else if strings.Contains(name, "prod") {
		envName = "prod"
	}

	data := make(map[string]string)
	data["debug"] = logName + "." + envName + ".debug"
	data["info"] = logName + "." + envName + ".info"
	data["warn"] = logName + "." + envName + ".warn"
	data["error"] = logName + "." + envName + ".error"
	data["fatal"] = logName + "." + envName + ".fatal"
	data["panic"] = logName + "." + envName + ".panic"

	return data
}

func (al *AddonsLogrus) WithTaskID(taskID string) *logrus.Entry {
	entry := al.WithField("task_id", taskID)
	return entry
}

func (al *AddonsLogrus) WithGrpcMetadata(metadata string) *logrus.Entry {
	entry := al.WithField("grpc_metadata", metadata)
	return entry
}
