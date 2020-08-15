package config

import "io"

// Configurator represents object which handels configuration. In particular it
// can create and open OS file for configuration. Moreover it has
// functionallity for saving actual configuration into file and parsing
// configuration from a file.
type Configurator interface {
	CreateFile() (io.ReadWriteCloser, error)
	OpenFile() (io.ReadWriteCloser, error)
	SaveToFile(io.Writer) error
	TryParseConfig(io.Reader) (Configurator, error)
}
