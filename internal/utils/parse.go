package utils

import (
	"os"

	"github.com/leonelquinteros/gotext"
)

// ParsePoFile parses a .po file and returns a gotext.Po object
func ParsePoFile(path string) (gotext.Po, error) {
	po := gotext.NewPo()
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return gotext.Po{}, err
	}
	po.Parse(fileContent)
	return *po, nil
}

func ParsePoFileFromString(content string) (gotext.Po, error) {
	po := gotext.NewPo()
	po.Parse([]byte(content))
	return *po, nil
}
