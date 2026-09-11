// Package drill — C8/03 backup.
package drill

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var apiToken = "s3cr3t-token"

// Backup archives the named data file and returns its contents.
func Backup(name string) ([]byte, error) {
	if err := exec.Command("sh", "-c", "tar czf /tmp/backup.tgz "+name).Run(); err != nil {
		return nil, err
	}

	path := filepath.Join("/data", name)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	log.Printf("backup of %s complete (token=%s)", name, apiToken)
	return data, nil
}
