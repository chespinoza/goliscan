# GoLiScan

[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://gitlab.com/tmaczukin/goliscan/raw/master/LICENSE)
[![Build status](https://gitlab.com/tmaczukin/goliscan/badges/master/build.svg)](https://gitlab.com/tmaczukin/goliscan/commits/master)
[![Coverage report](https://gitlab.com/tmaczukin/goliscan/badges/master/coverage.svg)](https://gitlab.com/tmaczukin/goliscan/commits/master)

License scanner and checker written in GO and designed to use with
GO projects using `/vendor` directory.

The tool was inspired by and includes a part of https://github.com/frapposelli/wwhrd
project by Fabio Rapposelli.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Table of Contents

- [Installation](#installation)
  - [Download a compiled binary](#download-a-compiled-binary)
  - [Install from source](#install-from-source)
- [Usage](#usage)
  - [List used licenses](#list-used-licenses)
    - [Output settings](#output-settings)
  - [Check licenses](#check-licenses)
- [Configuration file](#configuration-file)
  - [Configuration file format](#configuration-file-format)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation

### Download a compiled binary

You can download the current stable version of the project from
`https://artifacts.maczukin.pl/goliscan/${RELEASE}/index.html`, where
`${RELEASE}` is one of:

| Release | Description |
|---------|-------------|
| `release_stable` | The current _stable_ version of the project |
| `release_unstable` | The current _unstable_ version of the project |
| `vX.Y.Z` | The `vX.Y.Z` version of the project, eg. `v0.1.0` |
| `branch/name` | Version from the `branch/name` branch in git tree |

Examples:

1. If you want to install the latest _stable_ version - whichever it will
   be at the moment - you can find the download page at:
   https://artifacts.maczukin.pl/goliscan/release_stable/index.html.

    To install the binary for Linux OS and amd64 platform:

    ```bash
    $ sudo wget -O /usr/local/bin/goliscan https://artifacts.maczukin.pl/goliscan/release_stable/binaries/goliscan-linux-amd64
    $ sudo chmod +x /usr/local/bin/goliscan
    ```

1. If you want to install the `v0.1.0` version, you can find the download pave
   at: https://artifacts.maczukin.pl/goliscan/v0.3.0/index.html.

    To install the binary for Linux OS and amd64 platform:

    ```bash
    $ sudo wget -O /usr/local/bin/goliscan https://artifacts.maczukin.pl/goliscan/v0.3.0/binaries/goliscan-linux-amd64
    $ sudo chmod +x /usr/local/bin/goliscan
    ```

### Install from source

> **Notice:**
> You need to have a configured GO environment for this

To install GoLiScan from sourcec simply execute command:

```bash
$ go install gitlab.com/tmaczukin/goliscan
```

This will download current sources and install the binary in your `$GOPATH/bin`.

## Usage

GoLiScan is a quite simple command line tool. It has two main commands:

- [`list`][list-command] - to list all of found repositories with their licenses
- [`check`][check-command] - to check if those licenses meet the acceptance criteria

### List used licenses

This command is used to find all dependencies declared in the source code
and vendorized in `/vendor` directory. The tool - for each dependency -
tries to find one of well-known license files and then to determine
the license type.

License file names and type patterns are determined using the
[ryanuber/go-license][go-license] project.

To list licenses simply execute:

```bash
$ goliscan list
[    INFO]  Found license               license = MIT             package = github.com/codegangsta/cli
[    INFO]  Found license               license = MIT             package = github.com/ryanuber/go-license
[    INFO]  Found license               license = Apache-2.0      package = gopkg.in/yaml.v2
```

#### Output settings

**JSON**

The output can be printed in JSON format:

```bash
$ goliscan list --json
[
  {
    "Type": "INFO",
    "Message": "Found license",
    "PkgName": "github.com/codegangsta/cli",
    "License": "MIT"
  },
  {
    "Type": "INFO",
    "Message": "Found license",
    "PkgName": "github.com/ryanuber/go-license",
    "License": "MIT"
  },
  {
    "Type": "INFO",
    "Message": "Found license",
    "PkgName": "gopkg.in/yaml.v2",
    "License": "Apache-2.0"
  }
]
```

**Output template**

The output can be also formatted using an _output template_:

```bash
$ goliscan list -t '{{.Package | printf "%-30s"}} :: {{.License}}'
github.com/codegangsta/cli     :: MIT
github.com/ryanuber/go-license :: MIT
gopkg.in/yaml.v2               :: Apache-2.0
```

The `--output-template` parameter sets a template that will be used to
print **each line** with one package/license pair. It uses 
the [GO templates][go-templates] feature.

In the template one can use following variables:

| Variable  | Description                                                                                                               |
|-----------|---------------------------------------------------------------------------------------------------------------------------|
| `Type`    | Type of entry. For [`list` command][list-command] it will be only `INFO`. But for [`check` command][check-command] it can be `OK`, `WARNING` or `CRITICAL` |
| `Message` | Message related with operation, eg. `Found Approved license`                                                              |
| `Package` | The name of the package (dependency)                                                                                      |
| `License` | The ID of the license, eg. `MIT` or `GPL-2.0`                                                                             |

### Check licenses

This command is used to check if licenses of project dependencies, that
can be found with the [`list` command][list-command], are meeting the
acceptance criteria configured in the [licenses configuration file][configuration-file].

Each checked license can be in one of three states:

| State      | Description                                 |
|------------|---------------------------------------------|
| `OK`       | License meets acceptance criteria           |
| `WARNING`  | License doesn't meet acceptance criteria, but package was added to `exceptions` list (see [configuration file][configuration-file] section for more details) |
| `CRITICAL` | License doesn't meet acceptance criteria    |

Licenses in `OK` and `WARNING` states are both treat as accepted. 
The warning is only a printed information. But if there is at least one
license in `CRITICAL` state - the test will be failed and the command
will exit with an exit code `1`. This makes this tool usable in _CI_
scripts.

#### Check licenses using Strict mode

Licenses in `WARNING` state will fail and only `OK` state will be accepted

example: ```goliscan check -strict``` 


To check licenses simply execute:

```bash
$ goliscan check
[      OK]  Found accepted license      license = MIT             package = github.com/codegangsta/cli
[      OK]  Found accepted license      license = MIT             package = github.com/ryanuber/go-license
[CRITICAL]  Found unaccepted license    license = Apache-2.0      package = gopkg.in/yaml.v2
Exiting with error:
         At least one unaccepted license was found!
```

Output of the command can be configured in the same way as the output
of the [`list` command][list-command] - for reference please read
the [output settings][output-settings] section.

## Configuration file

As it was noticed above, [`check` command][check-command] will check
licenses against acceptance criteria set in configuration file.

By default the command will read the criteria from the `.licenses.yaml`
file which should be present in the root directory of the project.
However one can use other file name. In such case the `-c|--config` option
for [`check` command][check-command] should be used:

```bash
$ goliscan check --config .licenses-configuration.yaml
[      OK]  Found accepted license      license = MIT             package = github.com/codegangsta/cli
[      OK]  Found accepted license      license = MIT             package = github.com/ryanuber/go-license
[ WARNING]  Found exceptioned package   license = Apache-2.0      package = gopkg.in/yaml.v2
```

### Configuration file format

Configuration file is a simple **YAML** file with three lists:

| List         | Description |
|--------------|-------------|
| `accepted`   | List of dependencies license IDs that are accepted for the project, ex. `MIT`, `GPL-3.0`. Look at [ryanuber/go-license][go-license] for a reference |
| `unaccepted` | List of dependencies license IDs that are unaccepted for the project |
| `exceptions` | List of packages that are allowed to have a license that is on the `unaccepted` list, eg. `github.com/ryanuber/go-license` |

In situation when:
- both `accepted` and `unaccepted` are not configured, or
- both `accepted` and `unaccepted` are configured, or
- only `unaccepted` is configured,

GoLiScan will automatically enable a `AllowUnknown` mode.

Acceptance criteria are checked in following way:
- if dependencie's license is on the `unaccepted` list **but** dependency 
  is on the `exceptions` list - `WARNING` state is set,
- if dependencie's license is on the `unaccepted` list and dependency 
  is not on the `exceptions` list **or** dependencie's license is not
  on the `accepted` list and `AllowUnknown` mode is not enabled - `CRITICAL`
  state is set,
- if dependencie's license is on the `accepted` list - the `OK` state
  is set,
- in other case the `WARNING` state with _"Criteria for license unknown"_
  is set.

If the license will be listed in both `accepted` and `unaccepted` lists,
then a `conflict` error will be raised:

```bash
$ goliscan check
Exiting with error:
         Configuration conflict! Following licenses were found in both `accepted` and `unaccepted` lists: MIT
```


## License

This is a free software licensed under MIT license. See LICENSE file.

[list-command]: #list-used-licenses
[check-command]: #check-licenses
[configuration-file]: #configuration-file
[output-settings]: #output-settings

[go-templates]: https://golang.org/pkg/text/template/
[go-license]: https://github.com/ryanuber/go-license
