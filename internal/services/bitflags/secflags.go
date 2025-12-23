package bitflags

type SecFlag uint32

const (
	SecAuth SecFlag = 1 << iota
	SecSanitize
	SecSanitizeBody
)

// FromLegacyMap bridges your map[string]bool → flags.
func FromLegacyMap(m map[string]bool) SecFlag {
	var f SecFlag
	if m["secure"] {
		f |= SecAuth
	}
	if m["validateAndSanitize"] {
		f |= SecSanitize
	}
	if m["validateAndSanitizeBody"] {
		f |= SecSanitizeBody
	}
	return f
}

var secNames = map[string]SecFlag{
	"auth":          SecAuth,
	"sanitize":      SecSanitize,
	"sanitize_body": SecSanitizeBody,
}

func (f SecFlag) String() string { return FlagString(f, secNames) }
