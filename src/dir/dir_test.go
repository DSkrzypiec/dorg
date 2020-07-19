package dir

import (
	"fmt"
	"testing"
)

// TODO: mock file system
func TestCurrentDirScanner(t *testing.T) {
	config := DirScanConfig{}
	cds := CurrentDirScanner{}

	t1, err := cds.Scan("C:\\Go", config)
	if err != nil {
		t.Errorf("%v", err.Error())
	}

	fmt.Println(t1)
}
