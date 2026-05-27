package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// Logger é a instância global do logger
	Logger *logrus.Logger
)

// Init inicializa o logger com configuração padrão
func Init() {
	Logger = logrus.New()
	
	// Configurar output
	Logger.SetOutput(os.Stdout)
	
	// Configurar formato JSON para produção
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})
	
	// Configurar nível de log
	Logger.SetLevel(logrus.InfoLevel)
	
	// Adicionar campos padrão
	Logger = Logger.WithFields(logrus.Fields{
		"service": "mundo-invest",
		"version": "1.0.0",
	}).Logger
}

// InitWithLevel inicializa o logger com nível específico
func InitWithLevel(level string) {
	Init()
	
	switch level {
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "info":
		Logger.SetLevel(logrus.InfoLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}
}

// WithField cria um logger com um campo adicional
func WithField(key string, value interface{}) *logrus.Entry {
	return Logger.WithField(key, value)
}

// WithFields cria um logger com campos adicionais
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Logger.WithFields(fields)
}

// WithError cria um logger com erro
func WithError(err error) *logrus.Entry {
	return Logger.WithError(err)
}

// Debug loga mensagem de debug
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Debugf loga mensagem formatada de debug
func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

// Info loga mensagem de info
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Infof loga mensagem formatada de info
func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

// Warn loga mensagem de warning
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Warnf loga mensagem formatada de warning
func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

// Error loga mensagem de erro
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Errorf loga mensagem formatada de erro
func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

// Fatal loga mensagem fatal e termina o programa
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Fatalf loga mensagem formatada fatal e termina o programa
func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}
