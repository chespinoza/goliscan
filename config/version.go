package config

import (
	"bytes"
	"fmt"
	"runtime"
	"text/template"
	"time"

	"github.com/codegangsta/cli"
)

var VERSION = "dev"
var REVISION = "HEAD"
var BRANCH = "master"
var BUILT = "now"

var extendedInfoTemplate = `Version:      {{.Version}}
Git revision: {{.Revision}}
Git branch:   {{.Branch}}
GO version:   {{.GoVersion}}
OS/Arch:      {{.Os}}/{{.Arch}}
Built:        {{.Built}}`

type Version struct {
	Version   string
	Revision  string
	Branch    string
	GoVersion string
	Os        string
	Arch      string
	Built     string
}

func (v *Version) prepare() (err error) {
	var built time.Time

	if v.Built == "now" {
		built = time.Now()
	} else {
		built, err = time.Parse(time.RFC3339, v.Built)
		if err != nil {
			return
		}
	}

	v.Built = built.Format(time.RFC3339)

	return
}

func (v *Version) Printer(c *cli.Context) {
	info, _ := v.ExtendedInfo()
	fmt.Println(info)
}

func (v *Version) ShortInfo() string {
	return fmt.Sprintf("%s (%s)", v.Version, v.Revision)
}

func (v *Version) ExtendedInfo() (string, error) {
	tmpl, err := template.New("version-info").Parse(extendedInfoTemplate)
	if err != nil {
		return "", err
	}

	result := &bytes.Buffer{}
	err = tmpl.Execute(result, v)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

var instance *Version

func GetVersion() *Version {
	if instance == nil {
		instance = &Version{
			Version:   VERSION,
			Revision:  REVISION,
			Branch:    BRANCH,
			GoVersion: runtime.Version(),
			Os:        runtime.GOOS,
			Arch:      runtime.GOARCH,
			Built:     BUILT,
		}
		instance.prepare()
	}

	return instance
}

func init() {
	GetVersion()
}
