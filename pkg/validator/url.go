package validator

import "regexp"

var (
	urlRegex   = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	aliasRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]{3,20}$`)
)

var reservedAliases = map[string]bool{
	"api":    true,
	"urls":   true,
	"admin":  true,
	"health": true,
}

func IsValidURL(u string) bool {
	return urlRegex.MatchString(u)
}

func IsValidAlias(alias string) bool {
	return aliasRegex.MatchString(alias)
}

func IsReservedAlias(alias string) bool {
	return reservedAliases[alias]
}
