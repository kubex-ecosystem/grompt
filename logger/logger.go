package logger

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	l "github.com/rafa-mori/logz"
)

type gLog struct {
	l.Logger
	gLogLevel LogType
}

var (
	// debug is a boolean that indicates whether to log debug messages.
	debug bool
	// g is the global logger instance.
	g *gLog = &gLog{
		Logger:    l.GetLogger("GoBE - Test"),
		gLogLevel: LogTypeInfo,
	}
)

func init() {
	// Set the debug flag to true for testing purposes.
	debug = false
	// Initialize the global logger instance with a default logger.
	if g.Logger == nil {
		g = &gLog{
			Logger:    l.GetLogger("GoBE - Test"),
			gLogLevel: LogTypeInfo,
		}
	}
}

type LogType string

const (
	LogTypeNotice  LogType = "notice"
	LogTypeInfo    LogType = "info"
	LogTypeDebug   LogType = "debug"
	LogTypeError   LogType = "error"
	LogTypeWarn    LogType = "warn"
	LogTypeFatal   LogType = "fatal"
	LogTypePanic   LogType = "panic"
	LogTypeSuccess LogType = "success"
)

// SetDebug is a function that sets the debug flag for logging.
func SetDebug(d bool) { debug = d }

// LogObjLogger is a function that logs messages with the specified log type.
func LogObjLogger[T any](obj *T, logType string, messages ...string) {
	if obj == nil {
		g.ErrorCtx(fmt.Sprintf("log object (%s) is nil", reflect.TypeFor[T]()), map[string]any{
			"context":  "Log",
			"logType":  logType,
			"object":   obj,
			"msg":      messages,
			"showData": true,
		})
		return
	}
	var lgr l.Logger
	if objValueLogger := reflect.ValueOf(obj).Elem().MethodByName("GetLogger"); !objValueLogger.IsValid() {
		if objValueLogger = reflect.ValueOf(obj).Elem().FieldByName("Logger"); !objValueLogger.IsValid() {
			g.ErrorCtx(fmt.Sprintf("log object (%s) does not have a logger field", reflect.TypeFor[T]()), map[string]any{
				"context":  "Log",
				"logType":  logType,
				"object":   obj,
				"msg":      messages,
				"showData": true,
			})
			return
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
		//lgrC := objValueLogger.Call(nil)[0].Convert(reflect.TypeFor[l.Logger]())
		//if lgrC.IsNil() {
		//	lgrC = reflect.ValueOf(g.Logger)
		//}
		//if lgr = lgrC.Interface().(l.Logger); lgr == nil {
		lgr = g.Logger
		//}
	}
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		lgr.ErrorCtx("Log: unable to get caller information", nil)
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	ctxMessageMap := map[string]any{
		"context":  funcName,
		"file":     file,
		"line":     line,
		"showData": debug,
	}
	fullMessage := strings.Join(messages, " ")
	logType = strings.ToLower(logType)
	if logType != "" {
		if reflect.TypeOf(logType).ConvertibleTo(reflect.TypeFor[LogType]()) {
			lType := LogType(logType)
			ctxMessageMap["logType"] = logType
			logging(lgr, lType, fullMessage, ctxMessageMap)
		} else {
			lgr.ErrorCtx(fmt.Sprintf("logType (%s) is not valid", logType), ctxMessageMap)
		}
	} else {
		lgr.InfoCtx(fullMessage, ctxMessageMap)
	}
}

// Log is a function that logs messages with the specified log type and caller information.
func Log(logType string, messages ...any) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		g.ErrorCtx("Log: unable to get caller information", nil)
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	ctxMessageMap := map[string]any{
		"context":  funcName,
		"file":     file,
		"line":     line,
		"showData": debug,
	}
	fullMessage := ""
	if len(messages) > 0 {
		fullMessage = fmt.Sprintf("%v", messages[0:])
	}
	logType = strings.ToLower(logType)
	if logType != "" {
		if reflect.TypeOf(logType).ConvertibleTo(reflect.TypeFor[LogType]()) {
			lType := LogType(logType)
			ctxMessageMap["logType"] = logType
			logging(g.Logger, lType, fullMessage, ctxMessageMap)
		} else {
			g.ErrorCtx(fmt.Sprintf("logType (%s) is not valid", logType), ctxMessageMap)
		}
	} else {
		g.InfoCtx(fullMessage, ctxMessageMap)
	}
}

// logging is a helper function that logs messages with the specified log type.
func logging(lgr l.Logger, lType LogType, fullMessage string, ctxMessageMap map[string]interface{}) {
	debugCtx := debug
	if !debugCtx {
		if lType == "error" || lType == "fatal" || lType == "panic" || lType == "debug" {
			// If debug is false, set the debug value based on the logType
			debugCtx = true
		} else {
			debugCtx = false
		}
	}
	ctxMessageMap["showData"] = debugCtx
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
	default:
		lgr.InfoCtx(fullMessage, ctxMessageMap)
	}
	debugCtx = debug
}
