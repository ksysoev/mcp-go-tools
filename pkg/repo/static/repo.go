package static

type Config struct {
	Resources []Resource `mapstructure:"resources"`
}

type Resource struct {
	Name string `mapstructure:"name"`
	Data string `mapstructure:"data"`
}

type Service struct {
	config *Config
}

func New(cfg *Config) *Service {
	return &Service{
		config: cfg,
	}
}
