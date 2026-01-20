package logstash

import (
	"net"
	"os"
	"time"

	"github.com/faisalyudiansah/auth-service-template/pkg/config"
	"github.com/faisalyudiansah/auth-service-template/pkg/logger"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/jsonutils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger
var ServiceName string

func InitLogstash(cfg *config.Config) {
	host := cfg.LogstashConfig.LogstashHost
	port := cfg.LogstashConfig.LogstashPort
	timeout := cfg.LogstashConfig.LogstashTimeout
	tickerHealthCheck := cfg.LogstashConfig.LogstashTickerHealthCheck
	address := net.JoinHostPort(host, string(port))
	ServiceName = cfg.App.AppName

	Logger = logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	Logger.SetOutput(os.Stdout)

	logstashConn := &LogstashConnection{
		address: address,
		timeout: time.Duration(timeout) * time.Second,
	}

	if err := logstashConn.connect(); err != nil {
		logger.Log.Infof("Initial connection failed: %v\n", err)
	}

	startHealthCheck(logstashConn, time.Duration(tickerHealthCheck)*time.Second)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	if fields == nil {
		fields = logrus.Fields{}
	}
	fields["service"] = ServiceName
	return Logger.WithFields(fields)
}

func Error(key string, fields ...logrus.Fields) {
	f := logrus.Fields{}
	if len(fields) > 0 && fields[0] != nil {
		f = fields[0]
	}
	logAsync(logrus.ErrorLevel, "ERROR:"+ServiceName+":"+key, f)
}

func ErrorSync(key string, fields ...logrus.Fields) {
	f := logrus.Fields{}
	if len(fields) > 0 && fields[0] != nil {
		f = fields[0]
	}
	Logger.WithFields(f).Error("ERROR:" + ServiceName + ":" + key)
}

func Info(key string, fields ...logrus.Fields) {
	f := logrus.Fields{}
	if len(fields) > 0 && fields[0] != nil {
		f = fields[0]
	}
	logAsync(logrus.InfoLevel, ServiceName+":"+key, f)
}

func Warn(key string, fields ...logrus.Fields) {
	f := logrus.Fields{}
	if len(fields) > 0 && fields[0] != nil {
		f = fields[0]
	}
	logAsync(logrus.WarnLevel, "WARN:"+ServiceName+":"+key, f)
}

func Debug(key string, fields ...logrus.Fields) {
	f := logrus.Fields{}
	if len(fields) > 0 && fields[0] != nil {
		f = fields[0]
	}
	logAsync(logrus.DebugLevel, "DEBUG:"+ServiceName+":"+key, f)
}

func logAsync(level logrus.Level, message string, fields logrus.Fields) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Infof("Recovered in async log: %v", r)
			}
		}()
		entry := Logger.WithFields(fields)
		switch level {
		case logrus.InfoLevel:
			entry.Info(message)
		case logrus.WarnLevel:
			entry.Warn(message)
		case logrus.ErrorLevel:
			entry.Error(message)
		case logrus.DebugLevel:
			entry.Debug(message)
		}
	}()
}

func LogstashRequestInfo(c *gin.Context, payload interface{}, key string) {
	requestInfo := logrus.Fields{
		"id":         c.GetHeader("X-Request-ID"),
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"query":      c.Request.URL.Query(),
		"ip":         c.ClientIP(),
		"user_agent": c.Request.UserAgent(),
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	if payload != nil {
		requestInfo["payload"] = jsonutils.SafeMarshal(payload)
	}

	Info(key, requestInfo)
}

func LogstashResponseInfo(c *gin.Context, response interface{}, key string) {
	responseInfo := logrus.Fields{
		"id":         c.GetHeader("X-Request-ID"),
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"query":      c.Request.URL.Query(),
		"ip":         c.ClientIP(),
		"user_agent": c.Request.UserAgent(),
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	if response != nil {
		responseInfo["response"] = jsonutils.SafeMarshal(response)
	}

	Info(key, responseInfo)
}

func LogstashError(c *gin.Context, err error, payload interface{}, key string) {
	errorInfo := logrus.Fields{
		"id":         c.GetHeader("X-Request-ID"),
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"query":      c.Request.URL.Query(),
		"ip":         c.ClientIP(),
		"user_agent": c.Request.UserAgent(),
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	if payload != nil {
		errorInfo["payload"] = jsonutils.SafeMarshal(payload)
	}

	if err != nil {
		errorInfo["error"] = err.Error()
	}

	Error(key, errorInfo)
}
