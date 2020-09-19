package dir

import (
	"fmt"
	"os"
	"testing"
)

func TestDirDiff(t *testing.T) {
	config := DirScanConfig{}
	md1 := NewMockDir("Path1", "f1.go", "f2.cpp", "sub1_g1.go", "sub2_h1.html", "sub3_i3.jpeg")
	md2 := NewMockDir("Path1", "f1.go", "sub1_g1.go", "sub3_i3.jpeg", "sub3_i4.txt")

	tree1, err1 := Scan(md1, config)
	if err1 != nil {
		t.Errorf("Error while scanning first tree: %s", err1.Error())
	}

	tree2, err2 := Scan(md2, config)
	if err2 != nil {
		t.Errorf("Error while scanning second tree: %s", err2.Error())
	}

	diff0, diffTree0 := tree1.Diff(tree1)
	diff1, diffTree1 := tree1.Diff(tree2)
	diff2, diffTree2 := tree2.Diff(tree1)

	fmt.Println("Tree1 \\ Tree1:", diff0)
	fmt.Println(diffTree0)

	fmt.Println("Tree1 \\ Tree2:", diff1)
	fmt.Println(diffTree1)

	fmt.Println("Tree2 \\ Tree1:", diff2)
	fmt.Println(diffTree2)
}

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
