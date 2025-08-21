// Package logger provides a logging utility for Go applications.
package logger

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	manifest "github.com/rafa-mori/grompt/internal/module/info"
	l "github.com/rafa-mori/logz"
)

type GLog[T any] interface {
	GetLogger() l.Logger
	GetLogLevel() LogLevel
	GetShowTrace() bool
	GetDebug() bool
	SetLogLevel(string)
	SetDebug(bool)
	SetShowTrace(bool)
	ObjLog(*T, string, ...any)
	Log(string, ...any)
}

type gLog[T any] struct {
	l.Logger
	gLogLevel  LogLevel // Global log level
	gShowTrace bool     // Flag to show trace in logs
	gDebug     bool     // Flag to show debug messages
}

type LogType string
type LogLevel int

var (
	info      manifest.Manifest
	debug     bool
	showTrace bool
	logLevel  string
	g         *gLog[l.Logger] // Global logger instance
	Logger    GLog[l.Logger]
	err       error
)

const (
	// LogTypeDebug is the log type for debug messages.
	LogTypeDebug LogType = "debug"
	// LogTypeNotice is the log type for notice messages.
	LogTypeNotice LogType = "notice"
	// LogTypeInfo is the log type for informational messages.
	LogTypeInfo LogType = "info"
	// LogTypeWarn is the log type for warning messages.
	LogTypeWarn LogType = "warn"
	// LogTypeError is the log type for error messages.
	LogTypeError LogType = "error"
	// LogTypeFatal is the log type for fatal error messages.
	LogTypeFatal LogType = "fatal"
	// LogTypePanic is the log type for panic messages.
	LogTypePanic LogType = "panic"
	// LogTypeSuccess is the log type for success messages.
	LogTypeSuccess LogType = "success"
	// LogTypeAnswer is the log type for answer messages.
	LogTypeAnswer LogType = "answer"
	// LogTypeSilent is the log type for silent messages.
	LogTypeSilent LogType = "silent"
)

const (
	// LogLevelDebug 0
	LogLevelDebug LogLevel = iota
	// LogLevelNotice 1
	LogLevelNotice
	// LogLevelInfo 2
	LogLevelInfo
	// LogLevelSuccess 3
	LogLevelSuccess
	// LogLevelWarn 4
	LogLevelWarn
	// LogLevelError 5
	LogLevelError
	// LogLevelFatal 6
	LogLevelFatal
	// LogLevelPanic 7
	LogLevelPanic
	// LogLevelAnswer 8
	LogLevelAnswer
	// LogLevelSilent 9
	LogLevelSilent
)

func getEnvOrDefault[T string | int | bool](key string, defaultValue T) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	} else {
		valInterface := reflect.ValueOf(value)
		if valInterface.Type().ConvertibleTo(reflect.TypeFor[T]()) {
			return valInterface.Convert(reflect.TypeFor[T]()).Interface().(T)
		}
	}
	return defaultValue
}

func init() {
	if info == nil {
		info, err = manifest.GetManifest()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get info manifest: %v\n", err)
			os.Exit(1)
		}
		l.GetLogger(info.GetBin())
	}
	if Logger == nil {
		Logger = GetLogger[l.Logger](nil)
		if logger, ok := Logger.(*gLog[l.Logger]); ok {
			g = logger
			logLevel = getEnvOrDefault("GOBE_LOG_LEVEL", "error")
			debug = getEnvOrDefault("GOBE_DEBUG", false)
			showTrace = getEnvOrDefault("GOBE_SHOW_TRACE", false)
			//g.gLogLevel = LogLevelError
			g.gLogLevel = LogLevelInfo
			g.gShowTrace = showTrace
			g.gDebug = debug
		}
	}
}

func setLogLevel(logLevel string) {
	if g == nil || Logger == nil {
		_ = GetLogger[l.Logger](nil)
	}
	switch strings.ToLower(logLevel) {
	case "debug":
		g.gLogLevel = LogLevelDebug
		g.SetLevel("debug")
	case "info":
		g.gLogLevel = LogLevelInfo
		g.SetLevel("info")
	case "warn":
		g.gLogLevel = LogLevelWarn
		g.SetLevel("warn")
	case "error":
		g.gLogLevel = LogLevelError
		g.SetLevel("error")
	case "fatal":
		g.gLogLevel = LogLevelFatal
		g.SetLevel("fatal")
	case "panic":
		g.gLogLevel = LogLevelPanic
		g.SetLevel("panic")
	case "notice":
		g.gLogLevel = LogLevelNotice
		g.SetLevel("notice")
	case "success":
		g.gLogLevel = LogLevelSuccess
		g.SetLevel("success")
	case "silent":
		g.gLogLevel = LogLevelSilent
		g.SetLevel("silent")
	case "answer":
		g.gLogLevel = LogLevelAnswer
		g.SetLevel("answer")
	default:
		// logLevel = "error"
		// g.gLogLevel = LogLevelError
		logLevel = "info"
		g.gLogLevel = LogLevelInfo
		g.SetLevel(logLevel)
	}
}
func getShowTrace() bool {
	if debug {
		return true
	} else {
		if !showTrace {
			return false
		} else {
			return true
		}
	}
}
func willPrintLog(logType string) bool {
	if debug {
		return true
	} else {
		lTypeInt := LogLevelError
		switch strings.ToLower(logType) {
		case "debug":
			lTypeInt = LogLevelDebug
		case "info":
			lTypeInt = LogLevelInfo
		case "warn":
			lTypeInt = LogLevelWarn
		case "error":
			lTypeInt = LogLevelError
		case "notice":
			lTypeInt = LogLevelNotice
		case "success":
			lTypeInt = LogLevelSuccess
		case "fatal":
			lTypeInt = LogLevelFatal
		case "panic":
			lTypeInt = LogLevelPanic
		case "silent":
			lTypeInt = LogLevelSilent
		case "answer":
			lTypeInt = LogLevelAnswer
		default:
			lTypeInt = LogLevelError
		}
		return lTypeInt >= g.gLogLevel
	}
}
func getCtxMessageMap(logType, funcName, file string, line int) map[string]any {
	ctxMessageMap := map[string]any{
		"context":   funcName,
		"file":      file,
		"line":      line,
		"logType":   logType,
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   info.GetVersion(),
	}
	if !debug && !showTrace {
		ctxMessageMap["showData"] = false
	} else {
		ctxMessageMap["showData"] = getShowTrace()
	}
	if info != nil {
		ctxMessageMap["appName"] = info.GetName()
		ctxMessageMap["bin"] = info.GetBin()
		ctxMessageMap["version"] = info.GetVersion()
	}
	return ctxMessageMap
}
func getFuncNameMessage(lgr l.Logger) (string, int, string) {
	if lgr == nil {
		return "", 0, ""
	}
	if getShowTrace() {
		pc, file, line, ok := runtime.Caller(3)
		if !ok {
			lgr.ErrorCtx("Log: unable to get caller information", nil)
			return "", 0, ""
		}
		funcName := runtime.FuncForPC(pc).Name()
		if strings.Contains(funcName, "LogObjLogger") {
			pc, file, line, ok = runtime.Caller(4)
			if !ok {
				lgr.ErrorCtx("Log: unable to get caller information", nil)
				return "", 0, ""
			}
			funcName = runtime.FuncForPC(pc).Name()
		}
		return funcName, line, file
	}
	return "", 0, ""
}
func getFullMessage(messages ...any) string {
	fullMessage := ""
	for _, msg := range messages {
		if msg != nil {
			if str, ok := msg.(string); ok {
				fullMessage += str + " "
			} else {
				fullMessage += fmt.Sprintf("%v ", msg)
			}
		}
	}
	return strings.TrimSuffix(
		strings.TrimPrefix(
			strings.TrimSpace(fullMessage),
			" ",
		),
		" ",
	)
}

func SetDebug(d bool) {
	if g == nil || Logger == nil {
		_ = GetLogger[l.Logger](nil)
	}
	g.gDebug = d
	if d {
		g.SetLevel("debug")
	} else {
		switch g.gLogLevel {
		case LogLevelDebug:
			g.SetLevel("debug")
		case LogLevelInfo:
			g.SetLevel("info")
		case LogLevelWarn:
			g.SetLevel("warn")
		case LogLevelError:
			g.SetLevel("error")
		case LogLevelFatal:
			g.SetLevel("fatal")
		case LogLevelPanic:
			g.SetLevel("panic")
		case LogLevelNotice:
			g.SetLevel("notice")
		case LogLevelSuccess:
			g.SetLevel("success")
		case LogLevelSilent:
			g.SetLevel("silent")
		case LogLevelAnswer:
			g.SetLevel("answer")
		default:
			g.SetLevel("info")
		}
	}
}
func GetLogger[T any](obj *T) GLog[l.Logger] {
	if g == nil || Logger == nil {
		g = &gLog[l.Logger]{
			Logger:     l.GetLogger(info.GetBin()),
			gLogLevel:  LogLevelInfo,
			gShowTrace: showTrace,
			gDebug:     debug,
		}
		Logger = g
	}
	if obj == nil {
		return Logger
	}
	var lgr l.Logger
	if objValueLogger := reflect.ValueOf(obj).Elem().MethodByName("GetLogger"); !objValueLogger.IsValid() {
		if objValueLogger = reflect.ValueOf(obj).Elem().FieldByName("Logger"); !objValueLogger.IsValid() {
			g.ErrorCtx(fmt.Sprintf("log object (%s) does not have a logger field", reflect.TypeFor[T]()), map[string]any{
				"context":  "Log",
				"logType":  "error",
				"object":   obj,
				"msg":      "object does not have a logger field",
				"showData": getShowTrace(),
			})
			return g
		} else {
			lgrC := objValueLogger.Convert(reflect.TypeFor[l.Logger]())
			if lgrC.IsNil() {
				lgrC = reflect.ValueOf(g.Logger)
			}
			if lgr = lgrC.Interface().(l.Logger); lgr == nil {
				lgr = g.Logger
			}
		}
	} else {
		lgr = g
	}
	if lgr == nil {
		g.ErrorCtx(fmt.Sprintf("log object (%s) does not have a logger field", reflect.TypeFor[T]()), map[string]any{
			"context":  "Log",
			"logType":  "error",
			"object":   obj,
			"msg":      "object does not have a logger field",
			"showData": getShowTrace(),
		})
		return Logger
	}
	return &gLog[l.Logger]{
		Logger:     lgr,
		gLogLevel:  g.gLogLevel,
		gShowTrace: g.gShowTrace,
		gDebug:     g.gDebug,
	}
}
func LogObjLogger[T any](obj *T, logType string, messages ...any) {
	lgr := GetLogger(obj)
	if lgr == nil {
		g.ErrorCtx(fmt.Sprintf("log object (%s) does not have a logger field", reflect.TypeFor[T]()), map[string]any{
			"context":  "Log",
			"logType":  logType,
			"object":   obj,
			"msg":      messages,
			"showData": getShowTrace(),
		})
		return
	}

	fullMessage := getFullMessage(messages...)
	logType = strings.ToLower(logType)
	funcName, line, file := getFuncNameMessage(lgr.GetLogger())

	ctxMessageMap := getCtxMessageMap(logType, funcName, file, line)
	if logType != "" {
		if reflect.TypeOf(logType).ConvertibleTo(reflect.TypeFor[LogType]()) {
			lType := LogType(logType)
			logging(lgr.GetLogger(), lType, fullMessage, ctxMessageMap)
		} else {
			lgr.GetLogger().ErrorCtx(fmt.Sprintf("logType (%s) is not valid", logType), ctxMessageMap)
		}
	} else {
		lgr.GetLogger().InfoCtx(fullMessage, ctxMessageMap)
	}
}
func Log(logType string, messages ...any) {
	funcName, line, file := getFuncNameMessage(g.Logger)
	fullMessage := getFullMessage(messages...)
	logType = strings.ToLower(logType)
	ctxMessageMap := getCtxMessageMap(logType, funcName, file, line)
	if logType != "" {
		if reflect.TypeOf(logType).ConvertibleTo(reflect.TypeFor[LogType]()) {
			lType := LogType(logType)
			ctxMessageMap["logType"] = logType
			logging(g.Logger, lType, fullMessage, ctxMessageMap)
		} else {
			g.ErrorCtx(fmt.Sprintf("logType (%s) is not valid", logType), ctxMessageMap)
		}
	} else {
		logging(g.Logger, LogTypeInfo, fullMessage, ctxMessageMap)
	}
}
func logging(lgr l.Logger, lType LogType, fullMessage string, ctxMessageMap map[string]any) {
	lt := strings.ToLower(string(lType))
	if _, exist := ctxMessageMap["showData"]; !exist {
		ctxMessageMap["showData"] = getShowTrace()
	}
	if willPrintLog(lt) {
		switch lType {
		case LogTypeInfo:
			lgr.InfoCtx(fullMessage, ctxMessageMap)
		case LogTypeDebug:
			lgr.DebugCtx(fullMessage, ctxMessageMap)
		case LogTypeError:
			lgr.ErrorCtx(fullMessage, ctxMessageMap)
		case LogTypeWarn:
			lgr.WarnCtx(fullMessage, ctxMessageMap)
		case LogTypeNotice:
			lgr.NoticeCtx(fullMessage, ctxMessageMap)
		case LogTypeSuccess:
			lgr.SuccessCtx(fullMessage, ctxMessageMap)
		case LogTypeFatal:
			lgr.FatalCtx(fullMessage, ctxMessageMap)
		case LogTypePanic:
			lgr.FatalCtx(fullMessage, ctxMessageMap)
		case LogTypeSilent:
			lgr.SilentCtx(fullMessage, ctxMessageMap)
		case LogTypeAnswer:
			lgr.AnswerCtx(fullMessage, ctxMessageMap)
		default:
			lgr.InfoCtx(fullMessage, ctxMessageMap)
		}
	} else {
		ctxMessageMap["msg"] = fullMessage
		ctxMessageMap["showData"] = false
		lgr.DebugCtx("Log: message not printed due to log level", ctxMessageMap)
	}
}

func (g *gLog[T]) GetLogger() l.Logger                 { return g.Logger }
func (g *gLog[T]) GetLogLevel() LogLevel               { return g.gLogLevel }
func (g *gLog[T]) GetShowTrace() bool                  { return g.gShowTrace }
func (g *gLog[T]) GetDebug() bool                      { return g.gDebug }
func (g *gLog[T]) SetLogLevel(logLevel string)         { setLogLevel(logLevel) }
func (g *gLog[T]) SetShowTrace(showTrace bool)         { g.gShowTrace = showTrace }
func (g *gLog[T]) SetDebug(d bool)                     { SetDebug(d); g.gDebug = d }
func (g *gLog[T]) Log(logType string, messages ...any) { Log(logType, messages...) }
func (g *gLog[T]) ObjLog(obj *T, logType string, messages ...any) {
	LogObjLogger(obj, logType, messages...)
}

func (g *gLog[T]) Notice(m ...any) {
	if g == nil {
		_ = GetLogger[l.Logger](nil)
	}
	g.Log("notice", m...)
}
func (g *gLog[T]) Info(m ...any) {
	if g == nil {
		_ = GetLogger[l.Logger](nil)
	}
	g.Log("info", m...)
}
func (g *gLog[T]) Debug(m ...any) {
	if g == nil {
		_ = GetLogger[l.Logger](nil)
	}
	g.Log("debug", m...)
}
func (g *gLog[T]) Warn(m ...any) {
	if g == nil {
		_ = GetLogger[l.Logger](nil)
	}
	g.Log("warn", m...)
}
func (g *gLog[T]) Error(m ...any) {
	if g == nil {
		_ = GetLogger[l.Logger](nil)
	}
	g.Log("error", m...)
}
func (g *gLog[T]) Fatal(m ...any) {
	if g == nil {
		_ = GetLogger[l.Logger](nil)
	}
	g.Log("fatal", m...)
}
func (g *gLog[T]) Panic(m ...any) {
	if g == nil {
		_ = GetLogger[l.Logger](nil)
	}
	g.Log("fatal", m...)
}
func (g *gLog[T]) Success(m ...any) {
	if g == nil {
		_ = GetLogger[l.Logger](nil)
	}
	g.Log("success", m...)
}
func (g *gLog[T]) Silent(m ...any) {
	if g == nil {
		_ = GetLogger[l.Logger](nil)
	}
	g.Log("silent", m...)
}
func (g *gLog[T]) Answer(m ...any) {
	if g == nil {
		_ = GetLogger[l.Logger](nil)
	}
	g.Log("answer", m...)
}

func NewLogger[T any](prefix string) GLog[T] {
	return &gLog[T]{
		Logger:     l.NewLogger(prefix),
		gLogLevel:  LogLevelError,
		gShowTrace: false,
		gDebug:     false,
	}
}
