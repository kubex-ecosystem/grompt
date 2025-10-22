// Package kbx provides utilities for working with initialization arguments.
package kbx

import (
	"net"
	"os"
	"path/filepath"
	"reflect"
)

var (
	validKindMap = map[string]reflect.Kind{
		reflect.Struct.String():    reflect.Struct,
		reflect.Map.String():       reflect.Map,
		reflect.Slice.String():     reflect.Slice,
		reflect.Array.String():     reflect.Array,
		reflect.Chan.String():      reflect.Chan,
		reflect.Interface.String(): reflect.Interface,
		reflect.Ptr.String():       reflect.Ptr,
		reflect.String.String():    reflect.String,
		reflect.Int.String():       reflect.Int,
		reflect.Float32.String():   reflect.Float32,
		reflect.Float64.String():   reflect.Float64,
		reflect.Bool.String():      reflect.Bool,
		reflect.Uint.String():      reflect.Uint,
		reflect.Uint8.String():     reflect.Uint8,
		reflect.Uint16.String():    reflect.Uint16,
		reflect.Uint32.String():    reflect.Uint32,
		reflect.Uint64.String():    reflect.Uint64,
	}
)

type InitArgs struct {
	ConfigFile     string
	ConfigType     string
	ConfigDBFile   string
	ConfigDBType   string
	EnvFile        string
	LogFile        string
	Name           string
	Debug          bool
	ReleaseMode    bool
	IsConfidential bool
	Port           string
	Bind           string
	Address        string
	PubCertKeyPath string
	PubKeyPath     string
	PrivKeyPath    string
	Pwd            string
}

func NewInitArgs(
	configFile string,
	configType string,
	configDBFile string,
	configDBType string,
	envFile string,
	logFile string,
	name string,
	debug bool,
	releaseMode bool,
	isConfidential bool,
	port string,
	bind string,
	address string,
	pubCertKeyPath string,
	pubKeyPath string,
	pwd string,
) *InitArgs {
	configFile = GetValueOrDefaultSimple(configFile, os.ExpandEnv(DefaultGoBEConfigPath))
	configDBFile = GetValueOrDefaultSimple(configDBFile, "dbconfig.json")
	envFile = GetValueOrDefaultSimple(envFile, os.ExpandEnv(filepath.Join("$PWD", ".env")))
	logFile = GetValueOrDefaultSimple(
		logFile,
		filepath.Join(filepath.Dir(filepath.Dir(os.ExpandEnv(os.ExpandEnv(DefaultGoBEConfigPath)))), "logs", "gobe.log"),
	)
	port = GetValueOrDefaultSimple(port, "8088")
	bind = GetValueOrDefaultSimple(bind, "0.0.0.0")

	return &InitArgs{
		ConfigFile:     configFile,
		ConfigType:     filepath.Ext(configFile)[1:],
		ConfigDBFile:   configDBFile,
		ConfigDBType:   filepath.Ext(configDBFile)[1:],
		EnvFile:        envFile,
		LogFile:        logFile,
		Name:           GetValueOrDefaultSimple(name, "GoBE"),
		Debug:          GetValueOrDefaultSimple(debug, false),
		ReleaseMode:    GetValueOrDefaultSimple(releaseMode, false),
		IsConfidential: GetValueOrDefaultSimple(isConfidential, false),
		Port:           port,
		Bind:           bind,
		Address:        net.JoinHostPort(bind, port),
		PubCertKeyPath: GetValueOrDefaultSimple(pubCertKeyPath, os.ExpandEnv(DefaultGoBEKeyPath)),
		PubKeyPath:     GetValueOrDefaultSimple(pubKeyPath, os.ExpandEnv(DefaultGoBECertPath)),
		Pwd:            GetValueOrDefaultSimple(pwd, ""),
	}
}

func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetValueOrDefault[T any](value T, defaultValue T) (T, reflect.Type) {
	if !IsObjValid(value) {
		return defaultValue, reflect.TypeFor[T]()
	}
	return value, reflect.TypeFor[T]()
}

func GetValueOrDefaultSimple[T any](value T, defaultValue T) T {
	if !IsObjValid(value) {
		return defaultValue
	}
	return value
}

func IsObjValid(obj any) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return false
		}
		if v.Kind() == reflect.Ptr {
			if v.Elem().Kind() == reflect.Ptr && v.Elem().IsNil() {
				return false
			}
			v = v.Elem()
		}
	}
	if _, ok := validKindMap[v.Kind().String()]; !ok {
		return false
	}
	if !v.IsValid() {
		return false
	}
	if v.IsZero() {
		return false
	}
	if v.Kind() == reflect.String && v.Len() == 0 {
		return false
	}
	if (v.Kind() == reflect.Slice || v.Kind() == reflect.Map || v.Kind() == reflect.Array) && v.Len() == 0 {
		return false
	}
	if v.Kind() == reflect.Bool {
		return true
	}
	return true
}
