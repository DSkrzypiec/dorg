package dir

// DirScanner interface represents an IO directory scanner. Main method of this
// interface is Scan which takes a path and config and returns scanned content
// of the directory.
type DirScanner interface {
	Scan(string, DirScanConfig) (Dir, error)
	New(string, int) Dir
	String() string
}
