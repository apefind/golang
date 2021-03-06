package shellutil

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// getPromptPath shortens a path given a maximum length and a ratio for head and tail
func getPromptPath(path string, length int, ratio float64) string {
	l := len(path)
	if l > length {
		k := length - 3
		i, j := int((1.0-ratio)*float64(k)), l-int(ratio*float64(k))
		path = path[:i] + "..." + path[j:]
	}
	return path
}

func getPrompt(path string, length int, ratio float64) string {
	hostname, _ := os.Hostname()
	user, _ := user.Current()
	path = filepath.Clean(path)
	if strings.HasPrefix(path, user.HomeDir) {
		path = strings.Replace(filepath.Clean(path), user.HomeDir, "~", 1)
	}
	return user.Username + "@" + strings.Split(hostname, ".")[0] + ":" +
		getPromptPath(path, length, ratio) + "> "
}

// GetShellPrompt returns a nice shell prompt for the current directory
func GetShellPrompt(length int, ratio float64) string {
	path, _ := os.Getwd()
	return getPrompt(path, length, ratio)
}
