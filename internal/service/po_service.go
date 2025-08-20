package service

import (
	"github.com/leonelquinteros/gotext"
)

type UnTranslatedResult struct {
	Language string            `json:"language"`
	Terms    map[string]string `json:"terms"`
}

type PoService struct {
	poFile *gotext.Po
}

func NewPoService(poFile gotext.Po) *PoService {
	return &PoService{poFile: &poFile}
}

// ListAllUntranslated returns a map of untranslated messages up to the specified limit
// Returns map[msgid]msgstr where msgstr is empty or same as msgid
func (ps *PoService) ListAllUntranslated(limit int) UnTranslatedResult {
	result := make(map[string]string)
	count := 0

	// if limit is 0, then set the limit to 10
	if limit == 0 {
		limit = 10
	}
	// Get domain from Po file to access translations
	domain := ps.poFile.GetDomain()
	translations := domain.GetTranslations()

	for msgid, translation := range translations {
		// Skip if we've reached the limit
		if limit > 0 && count >= limit {
			break
		}

		// Check if translation is not translated or empty
		if !translation.IsTranslated() {
			result[msgid] = ""
			count++
		}
	}

	return UnTranslatedResult{
		Language: ps.poFile.Language,
		Terms:    result,
	}
}

// Translate sets a translation for a given key
func (ps *PoService) Translate(key, value string) {
	ps.poFile.Set(key, value)
}

// List returns a slice of translations with pagination support
// Returns map[msgid]msgstr for the specified range
func (ps *PoService) List(skip, take int) map[string]string {
	result := make(map[string]string)

	// Get domain from Po file to access translations
	domain := ps.poFile.GetDomain()
	translations := domain.GetTranslations()

	count := 0
	taken := 0

	for msgid, translation := range translations {
		// Skip the first 'skip' items
		if count < skip {
			count++
			continue
		}

		// Take 'take' items
		if taken >= take {
			break
		}

		result[msgid] = translation.Get()

		taken++
		count++
	}

	return result
}

// ToOutput returns the string representation of the Po file
func (ps *PoService) ToOutput() string {
	// Use MarshalText to get the Po file content
	data, err := ps.poFile.MarshalText()
	if err != nil {
		return ""
	}
	return string(data)
}
