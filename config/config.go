package config

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
	"strings"
)

type File struct {
	Accepted   []string `yaml:"accepted"`
	Unaccepted []string `yaml:"unaccepted"`
	Exceptions []string `yaml:"exceptions"`
}

func (f *File) AllowUnknown() bool {
	return (len(f.Accepted) > 0 && len(f.Unaccepted) > 0) ||
		(len(f.Accepted) == 0 && len(f.Unaccepted) > 0) ||
		(len(f.Accepted) == 0 && len(f.Unaccepted) == 0)
}

type Config struct {
	AllowUnknown bool

	Accepted   map[string]bool
	Unaccepted map[string]bool
	Exceptions map[string]bool

	path       string
	configFile *File
}

func (c *Config) ReadConfig() (err error) {
	err = c.readConfigFile()
	if err != nil {
		err = fmt.Errorf("Can't read config file: %s", err)
		return
	}

	c.AllowUnknown = c.configFile.AllowUnknown()
	c.Accepted = c.mapList(c.configFile.Accepted)
	c.Unaccepted = c.mapList(c.configFile.Unaccepted)
	c.Exceptions = c.mapList(c.configFile.Exceptions)

	err = c.findConflicts()

	return
}

func (c *Config) readConfigFile() (err error) {
	_, err = os.Stat(c.path)
	if os.IsNotExist(err) {
		return
	}

	configBytes, err := ioutil.ReadFile(c.path)
	if err != nil {
		return
	}

	c.configFile = &File{}
	err = yaml.Unmarshal(configBytes, c.configFile)

	return
}

func (c *Config) mapList(list []string) (mapped map[string]bool) {
	mapped = make(map[string]bool)

	for _, value := range list {
		mapped[value] = true
	}

	return
}

func (c *Config) findConflicts() (err error) {
	if len(c.Accepted) == 0 || len(c.Unaccepted) == 0 {
		return
	}

	var initial map[string]bool
	var searched map[string]bool

	if len(c.Accepted) > 0 {
		initial = c.Accepted
		searched = c.Unaccepted
	} else {
		initial = c.Unaccepted
		searched = c.Accepted
	}

	return c.searchForrepeatedEntries(initial, searched)
}

func (c *Config) searchForrepeatedEntries(initial, searched map[string]bool) (err error) {
	var conflicts []string
	for licenseName := range initial {
		if searched[licenseName] {
			conflicts = append(conflicts, licenseName)
		}
	}

	if len(conflicts) > 0 {
		err = fmt.Errorf(
			"Configuration conflict! Following licenses were found in both `accepted` and `unaccepted` lists: %s",
			strings.Join(conflicts, ", "),
		)
	}

	return
}

func NewConfig(path string) (config *Config) {
	if path == "" {
		path = ".licenses.yaml"
	}

	return &Config{
		path: path,
	}
}
