package dir

import (
	"os"
	"testing"
)

func TestFilesDiff(t *testing.T) {
	orig := []os.FileInfo{
		NewMockFileInfo("f1", false),
		NewMockFileInfo("f2", false),
		NewMockFileInfo("g3", false)}

	new := []os.FileInfo{
		NewMockFileInfo("f2", false),
		NewMockFileInfo("f1", false)}

	diff := filesDiff(orig, new)

	if len(diff) != 1 {
		t.Errorf("Expected diff to be length 1, got: %d", len(diff))
	}
	if diff[0].Name() != "g3" {
		t.Errorf("Expected 'g3', got: %s", diff[0].Name())
	}
	if diff[0].IsDir() {
		t.Errorf("Excepted to be regular file, not dir")
	}
}

func TestFilesDiffEmpty(t *testing.T) {
	orig := []os.FileInfo{
		NewMockFileInfo("f1", false),
		NewMockFileInfo("f2", false),
		NewMockFileInfo("g3", false)}

	emptyDiff := filesDiff(orig, orig)

	if len(emptyDiff) > 0 {
		t.Errorf("Expected to be empty, but got length = %d", len(emptyDiff))
	}
}

func TestFilesDiffFull(t *testing.T) {
	orig := []os.FileInfo{
		NewMockFileInfo("f1", false),
		NewMockFileInfo("f2", false),
		NewMockFileInfo("g3", false)}

	empty := []os.FileInfo{}
	fullDiff := filesDiff(orig, empty)

	if len(fullDiff) != len(orig) {
		t.Errorf("Expected length %d, got %d", len(orig), len(fullDiff))
	}

	for id, file := range orig {
		if file.Name() != fullDiff[id].Name() {
			t.Errorf("Diff should be identical as orig. Found diff in [%d]: %s vs %s",
				id, file.Name(), fullDiff[id].Name())
		}
	}
}
