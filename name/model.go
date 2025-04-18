package name

type Property struct {
	Transliteration string `bson:"transliteration"`
	Meaning         string `bson:"meaning"`
}

type Name struct {
	ID        int                 `bson:"_id"`
	Languages map[string]Property `bson:"languages"`
}

// Language codes
const (
	// list all language codes here
	LangDefault    = "default"
	LangArabic     = "ar" // Arabic
	LangBengali    = "bn" // Bengali
	LangChinese    = "zh" // Chinese
	LangDutch      = "nl" // Dutch
	LangEnglish    = "en" // English
	LangFrench     = "fr" // French
	LangGerman     = "de" // German
	LangGreek      = "el" // Greek
	LangHindi      = "hi" // Hindi
	LangIndonesian = "id" // Indonesian
	LangItalian    = "it" // Italian
	LangJapanese   = "ja" // Japanese
	LangKorean     = "ko" // Korean
	LangMalay      = "ms" // Malay
	LangPersian    = "fa" // Persian/Farsi
	LangPortuguese = "pt" // Portuguese
	LangRussian    = "ru" // Russian
	LangSpanish    = "es" // Spanish
	LangSwedish    = "sv" // Swedish
	LangThai       = "th" // Thai
	LangTurkish    = "tr" // Turkish
	LangUrdu       = "ur" // Urdu
	LangVietnamese = "vi" // Vietnamese
)
