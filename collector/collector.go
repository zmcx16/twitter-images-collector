package collector

// Collector struct
type Collector struct {
	config Config
}

// Init collector
func (collector *Collector) Init(configPath string) {
	collector.config.LoadConfig(configPath)
}

// Hello function
func Hello() string {
	return "Hello"
}
