package model

// Locale ...
type Locale string

var (
	// LocaleEn ...
	LocaleEn Locale = "EN"
	// LocaleHu ...
	LocaleHu Locale = "HU"
)

// Locales ...
var Locales map[string]Locale = map[string]Locale{
	"EN": LocaleEn,
	"HU": LocaleHu,
}
