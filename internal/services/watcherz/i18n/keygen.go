// Package i18n provides utilities for internationalization key generation.
package wi18nast

import (
	"crypto/sha1"
	"encoding/hex"
	"path/filepath"
	"regexp"
	"strings"
)

var nonWord = regexp.MustCompile(`[^a-z0-9]+`)

func GenKey(u Usage) string {
	compPath := componentPath(u.FilePath) // ex: components/dashboard
	elem := inferElementFromJSX(u.JSXCtx) // ex: h1, button, label, span...
	base := slugLower(u.Component)
	if base == "" {
		base = filepath.Base(compPath)
	}
	textSlug := slugLower(extractGuessText(u))

	parts := []string{}
	if compPath != "" {
		parts = append(parts, strings.ReplaceAll(compPath, "/", "."))
	}
	if base != "" {
		parts = append(parts, base)
	}
	if elem != "" {
		parts = append(parts, elem)
	}
	if textSlug != "" {
		parts = append(parts, textSlug)
	}
	key := strings.Join(parts, ".")
	if len(key) > 64 {
		key = key[:64] + "." + shortHash(key)
	}
	return key

}

func componentPath(p string) string {
	p = filepath.ToSlash(p)
	for _, cut := range []string{"/src/", "/app/"} {
		if idx := strings.Index(p, cut); idx >= 0 {
			return strings.TrimSuffix(p[idx+len(cut):], filepath.Ext(p))
		}
	}
	return strings.TrimSuffix(p, filepath.Ext(p))
}

func inferElementFromJSX(jsx string) string {
	s := strings.ToLower(jsx)
	switch {
	case strings.Contains(s, "<h1"):
		return "h1"
	case strings.Contains(s, "<h2"):
		return "h2"
	case strings.Contains(s, "<button"):
		return "button"
	case strings.Contains(s, "<label"):
		return "label"
	case strings.Contains(s, "<input"):
		return "input"
	case strings.Contains(s, "<p"):
		return "p"
	case strings.Contains(s, "<span"):
		return "span"
	default:
		return "text"
	}
}

func extractGuessText(u Usage) string {
	for _, ln := range u.Nearby {
		if i := strings.Index(ln, "  "); i >= 0 && i+2 < len(ln) {
			ln = ln[i+2:]
		}
		ln = strings.TrimSpace(ln)
		if s := between(ln, `"`, `"`); s != "" && likelyUIText(s) {
			return s
		}
		if s := between(ln, ">", "<"); s != "" && likelyUIText(s) {
			return s
		}
	}
	return u.Key
}

func between(s, a, b string) string {
	i := strings.Index(s, a)
	if i < 0 {
		return ""
	}
	j := strings.Index(s[i+len(a):], b)
	if j < 0 {
		return ""
	}
	return s[i+len(a) : i+len(a)+j]
}

func likelyUIText(s string) bool {
	if len(s) < 2 {
		return false
	}
	if strings.Contains(s, "://") {
		return false
	}
	return strings.Contains(s, " ") || regexp.MustCompile(`[A-Z][a-z]`).MatchString(s)
}

func slugLower(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonWord.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func shortHash(s string) string {
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])[:6]
}
