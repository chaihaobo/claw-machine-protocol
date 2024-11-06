package protocol

type (
	Logger interface {
		Printf(msg string, args ...any)
	}
	LoggerFunc func(msg string, args ...any)
)

func (l LoggerFunc) Printf(msg string, args ...any) {
	l(msg, args...)
}
