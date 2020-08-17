package middleware

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/8treenet/freedom"

	"github.com/kataras/golog"
)

var loggerPool sync.Pool

func init() {
	loggerPool = sync.Pool{
		New: func() interface{} {
			return &freedomLogger{}
		},
	}
}

func newFreedomLogger(traceName, traceId string) *freedomLogger {
	logger := loggerPool.New().(*freedomLogger)
	logger.traceId = traceId
	logger.traceName = traceName
	return logger
}

type freedomLogger struct {
	traceId   string
	traceName string
}

// Print prints a log message without levels and colors.
func (l *freedomLogger) Print(v ...interface{}) {
	traceField, traceId := l.traceField()
	v = append(v, traceField, traceId)
	freedom.Logger().Print(v...)
}

// Printf formats according to a format specifier and writes to `Printer#Output` without levels and colors.
func (l *freedomLogger) Printf(format string, args ...interface{}) {
	traceField, traceId := l.traceField()
	append := fmt.Sprintf(" %s:%s", traceField, traceId)
	freedom.Logger().Printf(format+append, args...)
}

// Println prints a log message without levels and colors.
// It adds a new line at the end, it overrides the `NewLine` option.
func (l *freedomLogger) Println(v ...interface{}) {
	traceField, traceId := l.traceField()
	v = append(v, traceField, traceId)
	freedom.Logger().Println(v...)
}

// Log prints a leveled log message to the output.
// This method can be used to use custom log levels if needed.
// It adds a new line in the end.
func (l *freedomLogger) Log(level golog.Level, v ...interface{}) {
	traceField, traceId := l.traceField()
	v = append(v, traceField, traceId)
	freedom.Logger().Log(level, v...)
}

// Logf prints a leveled log message to the output.
// This method can be used to use custom log levels if needed.
// It adds a new line in the end.
func (l *freedomLogger) Logf(level golog.Level, format string, args ...interface{}) {
	traceField, traceId := l.traceField()
	append := fmt.Sprintf(" %s:%s", traceField, traceId)
	freedom.Logger().Logf(level, format+append, args...)
}

// Fatal `os.Exit(1)` exit no matter the level of the freedomLogger.
// If the freedomLogger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func (l *freedomLogger) Fatal(v ...interface{}) {
	callerField, callerValue := l.callerField()
	traceField, traceId := l.traceField()
	v = append(v, callerField, callerValue, traceField, traceId)
	freedom.Logger().Fatal(v...)
}

// Fatalf will `os.Exit(1)` no matter the level of the freedomLogger.
// If the freedomLogger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func (l *freedomLogger) Fatalf(format string, args ...interface{}) {
	callerField, callerValue := l.callerField()
	traceField, traceId := l.traceField()
	append := fmt.Sprintf(" %s:%s %s:%s", callerField, callerValue, traceField, traceId)
	freedom.Logger().Fatalf(format+append, args...)
}

// Error will print only when freedomLogger's Level is error, warn, info or debug.
func (l *freedomLogger) Error(v ...interface{}) {
	callerField, callerValue := l.callerField()
	traceField, traceId := l.traceField()
	v = append(v, callerField, callerValue, traceField, traceId)
	freedom.Logger().Error(v...)
}

// Errorf will print only when freedomLogger's Level is error, warn, info or debug.
func (l *freedomLogger) Errorf(format string, args ...interface{}) {
	callerField, callerValue := l.callerField()
	traceField, traceId := l.traceField()
	append := fmt.Sprintf(" %s:%s %s:%s", callerField, callerValue, traceField, traceId)
	freedom.Logger().Errorf(format+append, args...)
}

// Warn will print when freedomLogger's Level is warn, info or debug.
func (l *freedomLogger) Warn(v ...interface{}) {
	callerField, callerValue := l.callerField()
	traceField, traceId := l.traceField()
	v = append(v, callerField, callerValue, traceField, traceId)
	freedom.Logger().Warn(v...)
}

// Warnf will print when freedomLogger's Level is warn, info or debug.
func (l *freedomLogger) Warnf(format string, args ...interface{}) {
	callerField, callerValue := l.callerField()
	traceField, traceId := l.traceField()
	append := fmt.Sprintf(" %s:%s %s:%s", callerField, callerValue, traceField, traceId)
	freedom.Logger().Warnf(format+append, args...)
}

// Info will print when freedomLogger's Level is info or debug.
func (l *freedomLogger) Info(v ...interface{}) {
	traceField, traceId := l.traceField()
	v = append(v, traceField, traceId)
	freedom.Logger().Info(v...)
}

// Infof will print when freedomLogger's Level is info or debug.
func (l *freedomLogger) Infof(format string, args ...interface{}) {
	traceField, traceId := l.traceField()
	append := fmt.Sprintf(" %s:%s", traceField, traceId)
	freedom.Logger().Infof(format+append, args...)
}

// Debug will print when freedomLogger's Level is debug.
func (l *freedomLogger) Debug(v ...interface{}) {
	callerField, callerValue := l.callerField()
	traceField, traceId := l.traceField()
	v = append(v, callerField, callerValue, traceField, traceId)
	freedom.Logger().Debug(v...)
}

// Debugf will print when freedomLogger's Level is debug.
func (l *freedomLogger) Debugf(format string, args ...interface{}) {
	callerField, callerValue := l.callerField()
	traceField, traceId := l.traceField()
	append := fmt.Sprintf(" %s:%s %s:%s", callerField, callerValue, traceField, traceId)
	freedom.Logger().Debugf(format+append, args...)
}

// traceField
func (l *freedomLogger) traceField() (string, string) {
	return l.traceName, l.traceId
}

// traceField
func (l *freedomLogger) callerField() (string, string) {
	_, file, line, _ := runtime.Caller(2)
	return "caller", fmt.Sprintf("%s:%d", file, line)
}
