package config

import "testing"

func TestConfigHash(t *testing.T) {
	c1 := Config{"file1", "Downloads", "Downloads"}
	c2 := Config{"file2", "Downloads", "Downloads"}
	c3 := Config{"file1", "Downloads_Crap", "Downloads"}
	c4 := Config{"file1", "Downloads", ""}

	if !c1.Equals(c1) {
		t.Errorf("The same config should be Equals to itself")
	}
	if c1.Equals(c2) {
		t.Errorf("Configs with different filepath should not be equal")
	}
	if c1.Equals(c3) {
		t.Errorf("Configs with different DownloadsPath should not be equal")
	}
	if c1.Equals(c4) {
		t.Errorf("Configs with different TargetBasePath should not be equal")
	}
}
