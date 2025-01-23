package bplparser2

import (
	"fmt"
	"path"
	"strings"
)

func TrimExtension(filename string) string {
	if ext := path.Ext(filename); len(ext) > 0 {
		filename = strings.TrimSuffix(filename, ext)
		filename = strings.TrimSuffix(filename, ".")
	}
	return filename
}

func ReplaceExtension(filename, replacement string) string {
	return fmt.Sprintf("%s%s", TrimExtension(filename), replacement)
}
