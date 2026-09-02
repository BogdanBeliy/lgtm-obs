package observability

// Configs defines OTLP connection, service identity, logging, and middleware settings.
type Configs struct {
	Path         string
	Protocol     string
	ServiceName  string
	Environment  string
	Version      string
	LogFormat    string
	Insecure     bool
	ExcludePaths []string
	SampleRate   float64
}

// NewConfigs copies cfg and applies defaults to empty fields.
func NewConfigs(cfg *Configs) *Configs {
	c := Configs{}
	if cfg != nil {
		c.Path = cfg.Path
		c.Protocol = cfg.Protocol
		c.ServiceName = cfg.ServiceName
		c.Environment = cfg.Environment
		c.Version = cfg.Version
		c.LogFormat = cfg.LogFormat
		c.Insecure = cfg.Insecure
		c.ExcludePaths = append([]string(nil), cfg.ExcludePaths...)
		c.SampleRate = cfg.SampleRate
	}

	c.applyDefaults()
	return &c
}

func (c *Configs) applyDefaults() {
	c.setProtocol()
	c.setPath()
	c.setServiceName()
	c.setEnvironment()
	c.setVersion()
	c.setLogFormat()
	c.setSampleRate()
}

func (c *Configs) setLogFormat() {
	if c.LogFormat == "" {
		c.LogFormat = "json"
	}
}

func (c *Configs) setProtocol() {
	if c.Protocol == "" {
		c.Protocol = "http"
	}
}

func (c *Configs) setPath() {
	if c.Path != "" {
		return
	}
	switch c.Protocol {
	case "http":
		c.Path = "localhost:4318"
	case "grpc":
		c.Path = "localhost:4317"
	}
}

func (c *Configs) setServiceName() {
	if c.ServiceName == "" {
		c.ServiceName = "example-service"
	}
}

func (c *Configs) setEnvironment() {
	if c.Environment == "" {
		c.Environment = "dev"
	}
}

func (c *Configs) setVersion() {
	if c.Version == "" {
		c.Version = "0.0.1"
	}
}

func (c *Configs) setSampleRate() {
	if c.SampleRate <= 0 || c.SampleRate > 1 {
		c.SampleRate = 1.0
	}
}
