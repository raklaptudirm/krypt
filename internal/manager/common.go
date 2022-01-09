package manager

import (
	"os"
	"path/filepath"
	"runtime"
)

var Pass *pass
var Auth *auth

func init() {
	dir := dataDir()
	os.Mkdir(dir, 0755)

	Pass = &pass{Dir: dir}
	Auth = &auth{Dir: dir}
}

// Data path precedence
// 1. XDG_DATA_HOME
// 2. LocalAppData (windows only)
// 3. HOME
func dataDir() string {
	a := os.Getenv("XDG_DATA_HOME")
	if a != "" {
		return filepath.Join(a, "krypt")
	}

	a = os.Getenv("LOCAL_APP_DATA")
	if runtime.GOOS == "windows" && a != "" {
		return filepath.Join(a, "Krypt")
	}

	a, _ = os.UserHomeDir()
	return filepath.Join(a, ".krypt")
}
