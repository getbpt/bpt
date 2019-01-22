# bpt  [![Build Status](https://travis-ci.org/getbpt/bpt.svg?branch=master)](https://travis-ci.org/getbpt/bpt) [![Coverage Status](https://img.shields.io/coveralls/github/getbpt/bpt/master.svg)](https://coveralls.io/github/getbpt/bpt) [![Go Report Card](https://goreportcard.com/badge/github.com/getbpt/bpt)](https://goreportcard.com/report/github.com/getbpt/bpt)
bpt (Bash Package Tool) allows you to declaratively retrieve shell scripts and binaries, exposing them to the current shell.

## Installing
```console
$ eval "$(curl -sfNL "https://raw.github.com/getbpt/bpt/master/get")"
```

## Building from source
To build `bpt` from the source code yourself you need to have a working Go environment with version 1.11.3 or greater installed.

```console
$ go get github.com/getbpt/bpt
$ cd $GOPATH/src/github.com/getbpt/bpt
$ ./build.sh
```

## Example
### Package 1: `https://bitbucket.org/organizaton/bpt_package1`
Package1 defines a library with a single shared function.
```shell
#!/usr/bin/env bash

function package1() {
  echo "Package1 executed!"
}
```

### Package 2: `https://bitbucket.org/organizaton/bpt_package2`
Package2 defines a library with a single shared function and code that runs when the library is included.
```shell
#!/usr/bin/env bash

function package2() {
  echo "Package2 executed!"
}

echo "Package2 loaded!"
```

### Binary1: `https://bitbucket.org/organization/bpt_binary1`
Binary package that contains platform specific versions of the binary.
The expected package layout is:
```none
root_
     |_ darwin _
     |          |_ binary
     |
     |_ linux __
                |_ binary
```

### Requirements: `~/bpt_requirements.txt`
The requirements file defines the packages that your script depends on.
I
```none
organizaton/bpt_package1
organizaton/bpt_package2
organizaton/bpt_binary2 kind:binary
```

### Script: `~/test.sh`
This script is depending on libraries defined in external packages and using `bpt` to ensure they are available.
```shell
#!/usr/bin/env bash

function initialize() {
  eval "$(curl -sfNL "https://raw.github.com/getbpt/bpt/master/get")"
  bpt get -r bpt_requirements.txt
}

function run() {
  package1
  package2
  echo "PATH: $PATH"
}

run "$@"
```

### Output
```console
~ $ ./test.sh
Package2 loaded!
Package1 executed!
Package2 executed!
PATH: ...:https-COLON-SLASH-SLASH-bitbucket.org-SLASH-organizaton-SLASH-bpt_binary2/darwin:...
```

## Packages
Packages are specified by URL to their git repository, for example `https://github.com/caarlos0/jvm`.  For `github.com`, `gitlab.com`, and `bitbucket.org` the schema prefix can be excluded, resulting in a package specification like `github.com/caarlos0/jvm`.  Packages stored on Bitbucket can leave off the host and specify the package like `getbpt/bpt`.

### Annotations

#### Kind
The kind annotation can be used to determine how a package should be handled.

##### sh
The default is `kind:sh`, which will source files in the root folder that match these globs:

* `*.plugin.sh`
* `*.sh`

```console
$ bpt get organization/module1 kind:sh
source /Users/user/.runtime/packages/https-COLON--SLASH--SLASH-bitbucket.org-SLASH-organization-SLASH-module1/library1.sh
source /Users/user/.runtime/packages/https-COLON--SLASH--SLASH-bitbucket.org-SLASH-organization-SLASH-module1/library2.sh
```

##### path
When `kind:path` is specified, the root folder of the package is put in the `$PATH`.

```console
$ bpt get organization/module1 kind:path
export PATH="/Users/user/.runtime/packages/https-COLON--SLASH--SLASH-bitbucket.org-SLASH-organization-SLASH-module1:$PATH"
```

##### binary
When `kind:binary` is specified, the root folder of the package is search for a platform folder (linux, darwin, etc.).  If one it found it is put in the `$PATH`.  See [syslist.go](https://github.com/golang/go/blob/master/src/go/build/syslist.go) for the full list of allowed values.

```console
$ bpt get organization/binary1 kind:binary
export PATH="/Users/user/.runtime/packages/https-COLON--SLASH--SLASH-bitbucket.org-SLASH-organization-SLASH-binary1/darwin:$PATH"
```

#### Branch
By default `bpt` will get the master branch of the package.  To specify a different branch include the `branch:<BRANCH_NAME>` annotation.

```console
$ bpt get organization/module1 branch:develop
source /Users/user/.runtime/packages/https-COLON--SLASH--SLASH-bitbucket.org-SLASH-organization-SLASH-module1/library1.sh
```

#### Path
You may specify a subfolder or a specific file within a package by including the `path:<PATH>` annotation.  This annotation can be used multiple times to include different paths in the same package.

```console
$ bpt get organization/module1 path:utils/library3.sh
source /Users/user/.runtime/packages/https-COLON--SLASH--SLASH-bitbucket.org-SLASH-organization-SLASH-module1/utils/library3.sh
$
$ bpt get organization/module1 path:library1.sh path:utils/library3.sh
source /Users/user/.runtime/packages/https-COLON--SLASH--SLASH-bitbucket.org-SLASH-organization-SLASH-module1/library1.sh
source /Users/user/.runtime/packages/https-COLON--SLASH--SLASH-bitbucket.org-SLASH-organization-SLASH-module1/utils/library3.sh
```

## Usage

```console
bpt provides a simple way to declaratively retrieve shell scripts, binaries, etc.
and expose them to the current shell.

To use, there are two steps to perform in a script:
  1. Initialize bpt: eval "$(bpt init)"
  2. Get package: bpt get org/repo

Usage:
  bpt [command]

Available Commands:
  get         Downloads a package and prints its source line
  help        Help about any command
  home        Prints the root directory where packages are kept
  init        Configure the shell environment for bpt
  list        Prints the currently installed packages
  purge       Purges a package from bpt
  purgeAll    Purges all packages from bpt
  update      Updates a previously downloaded package
  updateAll   Updates all previously downloaded packages
  version     Prints the bpt version

Flags:
  -h, --help   help for bpt

Use "bpt [command] --help" for more information about a command.
```

### Version

##### Help
```console
$ bpt help version
Prints the bpt version

Usage:
  bpt version [flags]

Flags:
  -h, --help   help for version
```

##### Example
```console
$ bpt version
bpt, version 0.0.1 (branch: master, revision: b78c5bbe6eb5760eb77bae5cffb0ffba0a742077)
  build user:       user@computer.local
  build date:       2019-01-23T19:10:38Z
  go version:       go1.11.3
```

### Get

##### Help
```console
$ bpt help get
Downloads a package and prints its source line

Usage:
  bpt get [<package>] [<options> ...] [flags]

Flags:
  -h, --help                  help for get
      --parallelism int       the max amount of tasks to launch in parallel (default 12)
  -r, --requirements string   the package requirements file
```

##### Example
```console
$ bpt get organization/module1
source /Users/user/.runtime/packages/https-COLON--SLASH--SLASH-bitbucket.org-SLASH-organization-SLASH-module1/library1.sh
source /Users/user/.runtime/packages/https-COLON--SLASH--SLASH-bitbucket.org-SLASH-organization-SLASH-module1/library2.sh
$ 
$ echo "organization/module1" > bash_requirements.txt
$ echo "organization/module2" >> bash_requirements.txt
$ bpt get -r bash_requirements.txt
source /Users/user/.runtime/packages/https-COLON--SLASH--SLASH-bitbucket.org-SLASH-organization-SLASH-module1/library1.sh
source /Users/user/.runtime/packages/https-COLON--SLASH--SLASH-bitbucket.org-SLASH-organization-SLASH-module1/library2.sh
source /Users/user/.runtime/packages/https-COLON--SLASH--SLASH-bitbucket.org-SLASH-organization-SLASH-module2/library3.sh
```

### Home

##### Help
```console
$ bpt help home
Prints the root directory where packages are kept

Usage:
  bpt home [flags]

Flags:
  -h, --help   help for home
```

##### Example
```console
$ bpt home
/Users/user/.runtime/packages
```

### Init


##### Help
```console
$ bpt help init
Configure the shell environment for bpt

Usage:
  bpt init [flags]

Flags:
  -h, --help   help for init
```

##### Example
```console
$ bpt init
#!/usr/bin/env bash
bpt() {
	case "$1" in
	get)
		source /dev/stdin <<<"$(cat <(/Users/user/bpt $@))" || /Users/user/bpt $@
		;;
	*)
		/Users/user/bpt $@
		;;
	esac
}
```

##### Usage
```console
eval "$(bpt init)"
```

### List

##### Help
```console
$ bpt help list
Prints the currently installed packages

Usage:
  bpt list [flags]

Flags:
  -h, --help   help for list
```

##### Example
```console
$ bpt list
https://bitbucket.org/organization/module1    /Users/user/.runtime/packages/https-COLON--SLASH--SLASH-bitbucket.org-SLASH-organization-SLASH-module1
https://bitbucket.org/organization/module2    /Users/user/.runtime/packages/https-COLON--SLASH--SLASH-bitbucket.org-SLASH-organization-SLASH-module2
```

### Purge

##### Help
```console
$ bpt help purge
Purges previously downloaded packages

Usage:
  bpt purge [<package> ...] [flags]

Flags:
  -h, --help   help for purge
      --parallelism int       the max amount of tasks to launch in parallel (default 12)
```

##### Example
```console
$ bpt purge
Removing all packages... 
$
$ bpt purge organization/module1
Removing organization/module1
```

### Update

##### Help
```console
$ bpt help update
Updates previously downloaded packages

Usage:
  bpt update [<package> ...] [flags]

Flags:
  -h, --help   help for update
      --parallelism int       the max amount of tasks to launch in parallel (default 12)
```

##### Example
```console
$ bpt update
Updating all packages in /Users/user/.runtime/packages...
$
$ bpt update organization/module1
Updating organization/module1
$
$ bpt update organization/module1 organization/module3
Updating organization/module1
Updating organization/module3
```

## Changelog

### 0.0.1
* Initial version

## License

The MIT License, see [LICENSE](https://github.com/getbpt/bpt/blob/master/LICENSE.md).
