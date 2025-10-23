// Package kbx has default configuration values
package kbx

const (
	KeyringService        = "kubex"
	DefaultKubexConfigDir = "$HOME/.kubex"

	DefaultGoBEKeyPath    = "$HOME/.kubex/gobe/gobe-key.pem"
	DefaultGoBECertPath   = "$HOME/.kubex/gobe/gobe-cert.pem"
	DefaultGoBECAPath     = "$HOME/.kubex/gobe/ca-cert.pem"
	DefaultGoBEConfigPath = "$HOME/.kubex/gobe/config/config.json"

	DefaultConfigDir        = "$HOME/.kubex/gdbase/config"
	DefaultConfigFile       = "$HOME/.kubex/gdbase/config.json"
	DefaultGDBaseConfigPath = "$HOME/.kubex/gdbase/config/config.json"
)

const (
	DefaultVolumesDir     = "$HOME/.kubex/volumes"
	DefaultRedisVolume    = "$HOME/.kubex/volumes/redis"
	DefaultPostgresVolume = "$HOME/.kubex/volumes/postgresql"
	DefaultMongoVolume    = "$HOME/.kubex/volumes/mongo"
	DefaultRabbitMQVolume = "$HOME/.kubex/volumes/rabbitmq"
)

const (
	DefaultRateLimitLimit  = 100
	DefaultRateLimitBurst  = 100
	DefaultRequestWindow   = 1 * 60 * 1000 // 1 minute
	DefaultRateLimitJitter = 0.1
)

const (
	DefaultMaxRetries = 3
	DefaultRetryDelay = 1 * 1000 // 1 second
)

const (
	DefaultMaxIdleConns          = 100
	DefaultMaxIdleConnsPerHost   = 100
	DefaultIdleConnTimeout       = 90 * 1000 // 90 seconds
	DefaultTLSHandshakeTimeout   = 10 * 1000 // 10 seconds
	DefaultExpectContinueTimeout = 1 * 1000  // 1 second
	DefaultResponseHeaderTimeout = 5 * 1000  // 5 seconds
	DefaultTimeout               = 30 * 1000 // 30 seconds
	DefaultKeepAlive             = 30 * 1000 // 30 seconds
	DefaultMaxConnsPerHost       = 100
)

const (
	DefaultLLMProvider    = "gemini"
	DefaultLLMModel       = "gemini-2.0-flash"
	DefaultLLMMaxTokens   = 1024
	DefaultLLMTemperature = 0.3
)

const (
	DefaultApprovalRequireForResponses = false
	DefaultApprovalTimeoutMinutes      = 15
)

const (
	DefaultServerPort = "8088"
	DefaultServerHost = "0.0.0.0"
)

type DBNameKey string

const (
	ContextDBNameKey = DBNameKey("dbName")
)

const (
	HeaderRequestIDKey = "X-Request-ID"
)

const (
	CookieSessionIDKey = "session_id"
)

const (
	AuthTypeNone   = "none"
	AuthTypeBasic  = "basic"
	AuthTypeBearer = "bearer"
	AuthTypeAPIKey = "api_key" // pragma: allowlist secret
	AuthTypeOIDC   = "oidc"
)
