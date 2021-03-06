package dir

import (
	"fmt"
	"time"

	"dorg/config"

	"github.com/pkg/errors"
)

// Listener represents a service which keeps listening on the directory and in case
// when new files occurs they are send on the channel to be categorized.
type Listener interface {
	Listen(chan<- Dir, chan<- error)
}

type DirListener struct {
	Path           string
	CurrentDir     *Dir
	ScanConfig     DirScanConfig
	ListenInterval time.Duration
}

// NewDirListener creates new DirListener based on confi.Config by scanning
// DownloadsPath catalog.
func NewDirListener(config config.Config) (DirListener, error) {
	// TODO: What about ExcludedDirs?
	excludedDirs := make(map[string]struct{})
	scanConfig := DirScanConfig{ExcludedDirs: excludedDirs}

	dir, err := Scan(FlatDir{config.DownloadsPath}, scanConfig)
	if err != nil {
		msg := fmt.Sprintf("Couldn't create new DirListener, probably couldn't scan [%s]",
			config.DownloadsPath)
		return DirListener{}, errors.Wrap(err, msg)
	}

	dl := DirListener{
		Path:           config.DownloadsPath,
		CurrentDir:     &dir,
		ScanConfig:     scanConfig,
		ListenInterval: 3 * time.Second,
	}

	return dl, nil
}

// Listen keep listening "downloads" directory file tree. When new files occurs
// they are send onto the newDirDiffChan channel in form of Dir to be later
// moved and categorized.
func (ds *DirListener) Listen(newDirDiffChan chan<- Dir, errChan chan<- error) {
	for {
		time.Sleep(ds.ListenInterval)

		newDir, err := Scan(FlatDir{ds.Path}, ds.ScanConfig)
		if err != nil {
			msg := fmt.Sprintf("Cannot scan file tree starting at [%s]", ds.Path)
			errChan <- errors.Wrap(err, msg)
		}

		if !ds.CurrentDir.Equals(newDir) {
			diff, diffDir := newDir.Diff(*ds.CurrentDir)

			if diff && !diffDir.IsEmpty() {
				diffDir.CleanEmptyDir()
				newDirDiffChan <- diffDir
				ds.CurrentDir = &newDir
			}
		}
	}
}
