package config

type ConfigurationReloader interface {
	ListenAndReload(chan<- Configurator, chan<- error)
}
