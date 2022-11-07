package app

import "os"

func IsRunningInDockerContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}
