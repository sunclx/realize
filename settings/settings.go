package settings

import (
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

// Settings defines a group of general settings
type Settings struct {
	Colors    `yaml:"-"`
	Config    `yaml:",inline" json:"config"`
	Resources `yaml:"resources" json:"resources"`
	Server    `yaml:"server,omitempty" json:"server,omitempty"`
}

// Config defines structural options
type Config struct {
	Create bool   `yaml:"-" json:"-"`
	Flimit uint64 `yaml:"flimit,omitempty" json:"flimit,omitempty"`
	Legacy `yaml:"legacy,omitempty" json:"legacy,omitempty"`
}

// Polling configuration
type Legacy struct {
	Status   bool          `yaml:"status" json:"status"`
	Interval time.Duration `yaml:"interval" json:"interval"`
}

// Server settings, used for the web panel
type Server struct {
	Status bool   `yaml:"status" json:"status"`
	Open   bool   `yaml:"open" json:"open"`
	Host   string `yaml:"host" json:"host"`
	Port   int    `yaml:"port" json:"port"`
}

// Resources defines the files generated by realize
type Resources struct {
	Config  string `yaml:"-" json:"-"`
	Outputs string `yaml:"outputs" json:"outputs"`
	Logs    string `yaml:"logs" json:"log"`
	Errors  string `yaml:"errors" json:"error"`
}

// Read from config file
func (s *Settings) Read(out interface{}) error {
	localConfigPath := s.Resources.Config
	if _, err := os.Stat(".realize/" + s.Resources.Config); err == nil {
		localConfigPath = ".realize/" + s.Resources.Config
	}
	content, err := s.Stream(localConfigPath)
	if err == nil {
		err = yaml.Unmarshal(content, out)
		return err
	}
	return err
}

// Record create and unmarshal the yaml config file
func (s *Settings) Record(out interface{}) error {
	if s.Config.Create {
		y, err := yaml.Marshal(out)
		if err != nil {
			return err
		}
		if _, err := os.Stat(".realize/"); os.IsNotExist(err) {
			if err = os.Mkdir(".realize/", 0770); err != nil {
				return s.Write(s.Resources.Config, y)
			}
		}
		return s.Write(".realize/"+s.Resources.Config, y)
	}
	return nil
}

// Remove realize folder
func (s *Settings) Remove() error {
	if _, err := os.Stat(".realize/"); !os.IsNotExist(err) {
		return os.RemoveAll(".realize/")
	}
	return nil
}

// Init configuration for general settings
func (s *Settings) Init(p *cli.Context) {
	s.Config = Config{
		Create: !p.Bool("no-config"),
		Flimit: p.Uint64("flimit"),
		Legacy: Legacy{
			Status:   p.Bool("legacy"),
			Interval: p.Duration("legacy-delay"),
		},
	}
	s.Server = Server{
		Status: !p.Bool("no-server"),
		Open:   p.Bool("serv-open"),
		Host:   p.String("serv-host"),
		Port:   p.Int("serv-port"),
	}
}
