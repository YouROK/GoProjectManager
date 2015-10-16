package project

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	ProjectName string
	MainPkgPath string
	GoPath      string
	GoRoot      string
	GoBin       string
	GoArch      string
	GoOS        string
	BuildOpt    struct {
		ShowCommand string
		Race        string
		Jobs        string
		BuildMode   string
		LDFlags     string

		BuildCMD     string
		BuildOutDir  string
		BuildOutFile string
	}
}

func loadConfig(fileName string) (*Config, error) {
	c := &Config{}
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func saveConfig(fileName string, cfg *Config) error {
	buf, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	os.Remove(fileName)
	return ioutil.WriteFile(fileName, buf, 0644)
}

func (p *Project) GetConfig() *Config {
	return p.cfg
}

func (p *Project) SetConfig(cfg *Config) {
	*p.cfg = *cfg
}
