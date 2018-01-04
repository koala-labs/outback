package cmd

type Config struct {
	Profile string         `mapstructure:"profile"`
	Region  string         `mapstructure:"region"`
	Repo    string         `mapstructure:"repo_url"`
	Env     []*Environment `mapstructure:"environments"`
}

type Environment struct {
	Branch     string   `mapstructure:"branch"`
	Region     string   `mapstructure:"region"`
	Cluster    string   `mapstructure:"cluster"`
	Services   []string `mapstructure:"services"`
	Dockerfile string   `mapstructure:"dockerfile"`
}

func getSelectedEnv(c string) (*Environment, error) {
	for _, env := range cfg.Env {
		if env.Cluster == c {
			return env, nil
		}
	}

	return nil, ErrClusterNotFound
}

func getSelectedService(services []string, svc string) (*string, error) {
	for _, service := range services {
		if service == svc {
			return &svc, nil
		}
	}

	return nil, ErrServiceNotFound
}
