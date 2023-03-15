package logger

import (
	"bytes"
	"fmt"
	"io"
	slog "log"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/rs/zerolog"
)

// HCLogger is a HC logger that uses zerolog.
type HCLogger struct {
	log     *zerolog.Logger
	implied []interface{}
	name    string
	level   hclog.Level
}

// NewHCLogger creates a HCLogger.
func NewHCLogger(zlog *zerolog.Logger) hclog.Logger {
	hclLogger := zlog.With().Str("log-source", "hc").Logger()
	return &HCLogger{log: &hclLogger}
}

// Emit a message and key/value pairs at a provided log level.
func (l *HCLogger) Log(level hclog.Level, msg string, args ...interface{}) {

	message := l.formatMesage(msg, args...)

	switch level { //nolint:exhaustive
	case hclog.Trace:
		l.log.Trace().Msg(message)
	case hclog.Debug:
		l.log.Debug().Msg(message)
	case hclog.Info:
		l.log.Info().Msg(message)
	case hclog.Warn:
		l.log.Warn().Msg(message)
	case hclog.Error:
		l.log.Error().Msg(message)
	default:
		// TODO: Handle hclog.NoLevel.
	}
}

// Emit a message and key/value pairs at the TRACE level.
func (l *HCLogger) Trace(msg string, args ...interface{}) {
	l.Log(hclog.Trace, msg, args...)
}

// Emit a message and key/value pairs at the DEBUG level.
func (l *HCLogger) Debug(msg string, args ...interface{}) {
	l.Log(hclog.Debug, msg, args...)
}

// Emit a message and key/value pairs at the INFO level.
func (l *HCLogger) Info(msg string, args ...interface{}) {
	l.Log(hclog.Info, msg, args...)
}

// Emit a message and key/value pairs at the WARN level.
func (l *HCLogger) Warn(msg string, args ...interface{}) {
	l.Log(hclog.Warn, msg, args...)
}

// Emit a message and key/value pairs at the ERROR level.
func (l *HCLogger) Error(msg string, args ...interface{}) {
	l.Log(hclog.Error, msg, args...)
}

// Indicate if TRACE logs would be emitted. This and the other Is* guards
// are used to elide expensive logging code based on the current level.
func (l *HCLogger) IsTrace() bool {
	return l.log.GetLevel() == zerolog.TraceLevel
}

// Indicate if DEBUG logs would be emitted. This and the other Is* guards.
func (l *HCLogger) IsDebug() bool {
	return l.log.GetLevel() == zerolog.DebugLevel
}

// Indicate if INFO logs would be emitted. This and the other Is* guards.
func (l *HCLogger) IsInfo() bool {
	return l.log.GetLevel() == zerolog.InfoLevel
}

// Indicate if WARN logs would be emitted. This and the other Is* guards.
func (l *HCLogger) IsWarn() bool {
	return l.log.GetLevel() == zerolog.WarnLevel
}

// Indicate if ERROR logs would be emitted. This and the other Is* guards.
func (l *HCLogger) IsError() bool {
	return l.log.GetLevel() == zerolog.ErrorLevel
}

// ImpliedArgs returns With key/value pairs.
func (l *HCLogger) ImpliedArgs() []interface{} {
	// Not implemented
	return l.implied
}

// Creates a sublogger that will always have the given key/value pairs.
func (l *HCLogger) With(args ...interface{}) hclog.Logger {
	var stringArgs []string
	for _, args := range args {
		stringArgs = append(stringArgs, args.(string))
	}
	logger := l.log.With().Strs("log-source", stringArgs).Logger()
	return &HCLogger{
		log:     &logger,
		name:    l.name,
		implied: args,
	}
}

// Returns the Name of the logger.
func (l *HCLogger) Name() string {
	return l.name
}

// Create a logger that will prepend the name string on the front of all messages.
// If the logger already has a name, the new value will be appended to the current
// name. That way, a major subsystem can use this to decorate all it's own logs
// without losing context.
func (l *HCLogger) Named(name string) hclog.Logger {
	logger := l.log.With().Str("log-source", name).Logger()
	return &HCLogger{
		log:     &logger,
		name:    name,
		implied: l.implied,
	}
}

// Create a logger that will prepend the name string on the front of all messages.
// This sets the name of the logger to the value directly, unlike Named which honor
// the current name as well.
func (l *HCLogger) ResetNamed(name string) hclog.Logger {
	logger := l.log.With().Logger()
	return &HCLogger{
		log:     &logger,
		name:    name,
		implied: l.implied,
	}
}

// Updates the level. This should affect all related loggers as well,
// unless they were created with IndependentLevels. If an
// implementation cannot update the level on the fly, it should no-op.
func (l *HCLogger) SetLevel(level hclog.Level) {
	l.level = level
}

// Returns the current level.
func (l *HCLogger) GetLevel() hclog.Level {
	return l.level
}

// Return a value that conforms to the stdlib log.Logger interface.
func (l *HCLogger) StandardLogger(opts *hclog.StandardLoggerOptions) *slog.Logger {
	return nil
}

// Return a value that conforms to io.Writer, which can be passed into log.SetOutput().
func (l *HCLogger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return io.Discard
}

// nolint:funlen,gocyclo // imported from https://github.com/hashicorp/go-hclog
func (l *HCLogger) formatMesage(msg string, args ...interface{}) string {
	writer := stringWriter{}

	if l.name != "" {
		writer.WriteString(l.name)
		writer.WriteString(": ")
	}

	writer.WriteString(msg)

	args = append(l.implied, args...)

	var stacktrace hclog.CapturedStacktrace

	if len(args) > 0 {
		if len(args)%2 != 0 {
			cs, ok := args[len(args)-1].(hclog.CapturedStacktrace)
			if ok {
				args = args[:len(args)-1]
				stacktrace = cs
			} else {
				extra := args[len(args)-1]
				args = append(args[:len(args)-1], "EXTRA_VALUE_AT_END", extra)
			}
		}

		writer.WriteByte(':')

	FOR:
		for i := 0; i < len(args); i += 2 {
			var (
				val string
				raw bool
			)

			switch st := args[i+1].(type) {
			case string:
				val = st
				if st == "" {
					val = `""`
				}
			case int:
				val = strconv.FormatInt(int64(st), 10)
			case int64:
				val = strconv.FormatInt(st, 10)
			case int32:
				val = strconv.FormatInt(int64(st), 10)
			case int16:
				val = strconv.FormatInt(int64(st), 10)
			case int8:
				val = strconv.FormatInt(int64(st), 10)
			case uint:
				val = strconv.FormatUint(uint64(st), 10)
			case uint64:
				val = strconv.FormatUint(st, 10)
			case uint32:
				val = strconv.FormatUint(uint64(st), 10)
			case uint16:
				val = strconv.FormatUint(uint64(st), 10)
			case uint8:
				val = strconv.FormatUint(uint64(st), 10)
			case hclog.Hex:
				val = "0x" + strconv.FormatUint(uint64(st), 16)
			case hclog.Octal:
				val = "0" + strconv.FormatUint(uint64(st), 8)
			case hclog.Binary:
				val = "0b" + strconv.FormatUint(uint64(st), 2)
			case hclog.CapturedStacktrace:
				stacktrace = st
				continue FOR
			case hclog.Format:
				val = fmt.Sprintf(st[0].(string), st[1:]...)
			case hclog.Quote:
				raw = true
				val = strconv.Quote(string(st))
			default:
				v := reflect.ValueOf(st)
				if v.Kind() == reflect.Slice {
					val = l.renderSlice(v)
					raw = true
				} else {
					val = fmt.Sprintf("%v", st)
				}
			}

			var key string

			switch st := args[i].(type) {
			case string:
				key = st
			default:
				key = fmt.Sprintf("%s", st)
			}

			if strings.Contains(val, "\n") {
				writer.WriteString("\n  ")
				writer.WriteString(key)
				writer.WriteString("=\n")
				writeIndent(&writer, val, "  | ")
				writer.WriteString("  ")
			} else if !raw && strings.ContainsAny(val, " \t") {
				writer.WriteByte(' ')
				writer.WriteString(key)
				writer.WriteByte('=')
				writer.WriteByte('"')
				writer.WriteString(val)
				writer.WriteByte('"')
			} else {
				writer.WriteByte(' ')
				writer.WriteString(key)
				writer.WriteByte('=')
				writer.WriteString(val)
			}
		}
	}

	if stacktrace != "" {
		writer.WriteString(string(stacktrace))
		writer.WriteString("\n")
	}

	return writer.String()
}

func (l HCLogger) renderSlice(v reflect.Value) string {
	var buf bytes.Buffer

	buf.WriteRune('[')

	for i := 0; i < v.Len(); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}

		sv := v.Index(i)

		var val string

		switch sv.Kind() { //nolint:exhaustive
		case reflect.String:
			val = strconv.Quote(sv.String())
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			val = strconv.FormatInt(sv.Int(), 10)
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val = strconv.FormatUint(sv.Uint(), 10)
		default:
			val = fmt.Sprintf("%v", sv.Interface())
			if strings.ContainsAny(val, " \t\n\r") {
				val = strconv.Quote(val)
			}
		}

		buf.WriteString(val)
	}

	buf.WriteRune(']')

	return buf.String()
}

type stringWriter struct {
	value string
}

func (w *stringWriter) WriteString(s string) {
	w.value = fmt.Sprintf("%s%s", w.value, s)
}

func (w *stringWriter) WriteByte(c byte) {
	w.value = fmt.Sprintf("%s%c", w.value, c)
}

func (w *stringWriter) String() string {
	return w.value
}

func writeIndent(w *stringWriter, str, indent string) {
	for {
		nl := strings.IndexByte(str, "\n"[0])
		if nl == -1 {
			if str != "" {
				w.WriteString(indent)
				w.WriteString(str)
				w.WriteString("\n")
			}
			return
		}

		w.WriteString(indent)
		w.WriteString(str[:nl])
		w.WriteString("\n")
		str = str[nl+1:]
	}
}
