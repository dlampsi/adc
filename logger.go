package adc

// Client logger interface.
type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
}

// New not operational logger.
func newNopLogger() *nopLogger {
	return &nopLogger{}
}

type nopLogger struct{}

func (l *nopLogger) Debug(args ...interface{})                   {}
func (l *nopLogger) Debugf(template string, args ...interface{}) {}
