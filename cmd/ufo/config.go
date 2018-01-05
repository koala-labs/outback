package cmd

type Config struct {
	Profile  string     `mapstructure:"profile"`
	Region   string     `mapstructure:"region"`
	Repo     string     `mapstructure:"repo"`
	Clusters []*Cluster `mapstructure:"clusters"`
	Tasks    []*Task    `mapstructure:"tasks"`
}

type Cluster struct {
	Name       string   `mapstructure:"name"`
	Branch     string   `mapstructure:"branch"`
	Region     string   `mapstructure:"region"`
	Services   []string `mapstructure:"services"`
	Dockerfile string   `mapstructure:"dockerfile"`
}

type Task struct {
	Name    string `mapstructure:"name"`
	Command string `mapstructure:"command"`
}

func (c *Config) getSelectedCluster(cluster string) (*Cluster, error) {
	for _, c := range c.Clusters {
		if c.Name == cluster {
			return c, nil
		}
	}

	return nil, ErrClusterNotFound
}

func (c *Config) getSelectedService(services []string, service string) (*string, error) {
	for _, s := range services {
		if s == service {
			return &s, nil
		}
	}

	return nil, ErrServiceNotFound
}

func (c *Config) getCommand(name string) (*string, error) {
	for _, t := range c.Tasks {
		if t.Name == name {
			return &t.Command, nil
		}
	}

	return nil, ErrCommandNotFound
}
