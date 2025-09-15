package logger

import (
	logrus "github.com/sirupsen/logrus"
)

type Logger struct {
	Key   string
	Value interface{}
}

func WithField(key string, value interface{}) *Logger {
	e := &Logger{
		Key:   key,
		Value: value,
	}
	return e
}

func (e *Logger) Debug(args ...interface{}) {
	logrus.WithField(e.Key, e.Value).Debug(args...)
}

func (e *Logger) Debugf(format string, args ...interface{}) {
	logrus.WithField(e.Key, e.Value).Debugf(format, args...)
}

func (e *Logger) Info(args ...interface{}) {
	logrus.WithField(e.Key, e.Value).Info(args...)
}

func (e *Logger) Infof(format string, args ...interface{}) {
	logrus.WithField(e.Key, e.Value).Infof(format, args...)
}

func (e *Logger) Error(args ...interface{}) {
	logrus.WithField(e.Key, e.Value).Error(args...)
}

func (e *Logger) Errorf(format string, args ...interface{}) {
	logrus.WithField(e.Key, e.Value).Errorf(format, args...)
}

func (e *Logger) Warn(args ...interface{}) {
	logrus.WithField(e.Key, e.Value).Warn(args...)
}

func (e *Logger) Warning(args ...interface{}) {
	logrus.WithField(e.Key, e.Value).Warning(args...)
}

func (e *Logger) Warningf(format string, args ...interface{}) {
	logrus.WithField(e.Key, e.Value).Warningf(format, args...)
}

func (e *Logger) Trace(args ...interface{}) {
	logrus.WithField(e.Key, e.Value).Trace(args...)
}

func (e *Logger) Tracef(format string, args ...interface{}) {
	logrus.WithField(e.Key, e.Value).Tracef(format, args...)
}

func (e *Logger) Fatal(args ...interface{}) {
	logrus.WithField(e.Key, e.Value).Fatal(args...)
}

func Fatal(text ...interface{}) {
	logrus.Fatal(text...)
}

func Println(text ...interface{}) {
	logrus.Println(text...)
}

func SetLogLevel(level logrus.Level) {
	logrus.SetLevel(level)
}

func SetFormatter(formatter logrus.Formatter) {
	logrus.SetFormatter(formatter)
}

