#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
RUNTIME_DIR="$ROOT_DIR/.runtime"
BINARY_DIR="$RUNTIME_DIR/.bin"
VERSION=$(cat ${ROOT_DIR}/.version)
GLIDE="${BINARY_DIR}/glide"
GOX="gox"
GOMETALINTER="${BINARY_DIR}/gometalinter"
UPX="${BINARY_DIR}/upx"
GOENV_DIR="$RUNTIME_DIR/.goenv"
export GOENV_ROOT="$GOENV_DIR"
export PATH="$GOENV_ROOT/bin:${BINARY_DIR}:$PATH"

function verbose() { echo -e "$*"; }
function error() { echo -e "ERROR: $*" 1>&2; }
function fatal() { echo -e "ERROR: $*" 1>&2; exit 1; }
function pushd () { command pushd "$@" > /dev/null; }
function popd () { command popd > /dev/null; }

function trap_add() {
  local localtrap_add_cmd=$1; shift || fatal "${FUNCNAME[*]} usage error: $?"
  for trap_add_name in "$@"; do
    trap -- "$(
      extract_trap_cmd() { printf '%s\n' "$3"; }
      eval "extract_trap_cmd $(trap -p "${trap_add_name}")"
      printf '%s\n' "${trap_add_cmd}"
    )" "${trap_add_name}" || fatal "unable to add to trap ${trap_add_name}: $?"
  done
}
declare -f -t trap_add

function get_platform() {
  local unameOut
  unameOut="$(uname -s)" || fatal "unable to get platform type: $?"
  case "${unameOut}" in
    Linux*)
      echo "linux"
    ;;
    Darwin*)
      echo "darwin"
    ;;
    *)
      echo "Unsupported machine type :${unameOut}"
      exit 1
    ;;
  esac
}

PLATFORM=$(get_platform)
GLIDE_URL="https://github.com/Masterminds/glide/releases/download/v0.13.1/glide-v0.13.1-$PLATFORM-amd64.tar.gz"
GOMETALINTER_URL="https://github.com/alecthomas/gometalinter/releases/download/v2.0.12/gometalinter-2.0.12-$PLATFORM-amd64.tar.gz"
UPX_URL="https://github.com/kadaan/upx/releases/download/20181231/upx_$PLATFORM"

function get_go_version() {
  local go_version
  go_version="$(cat "$ROOT_DIR/.go-version")" || fatal "failed to read go version: $?"
  echo "$go_version"
}

function download_go() {
  if [[ ! -d "$GOENV_DIR" ]]; then
    git clone https://github.com/syndbg/goenv.git "$GOENV_DIR" || fatal "failed to get goenv: $?"
  fi
  eval "$(goenv init - --no-rehash)"
  local go_version="$(get_go_version)"
  goenv install ${go_version} --skip-existing || fatal "failed to install go ${go_version}: $?"
  goenv rehash 2> /dev/null || fatal "failed to rehash goenv: $?"
}

function activate_go() {
  eval "$(goenv init - --no-rehash)"
  local go_version="$(get_go_version)"
  goenv shell ${go_version} || fatal "Failed to run switch to go ${go_version}: $?"
}

function download_gometalinter() {
  if [[ ! -f "$GOMETALINTER" ]]; then
    verbose "   --> $GOMETALINTER"
    local tmpdir=`mktemp -d`
    trap_add "rm -rf $tmpdir" EXIT
    pushd ${tmpdir}
    curl -L -s -O ${GOMETALINTER_URL} || fatal "failed to download '$GOMETALINTER_URL': $?"
    for i in *.tar.gz; do
      [[ "$i" = "*.tar.gz" ]] && continue
      tar xzf "$i" -C ${tmpdir} --strip-components 1 && rm -r "$i"
    done
    popd
    mkdir -p ${BINARY_DIR}
    cp ${tmpdir}/* ${BINARY_DIR}/
  fi
}

function download_glide() {
  if [[ ! -f "$GLIDE" ]]; then
    verbose "   --> $GLIDE"
    local tmpdir=`mktemp -d`
    trap_add "rm -rf $tmpdir" EXIT
    pushd ${tmpdir}
    curl -L -s -O ${GLIDE_URL} || fatal "failed to download '$GLIDE_URL': $?"
    for i in *.tar.gz; do
      [[ "$i" = "*.tar.gz" ]] && continue
      tar xzf "$i" -C ${tmpdir} --strip-components 1 && rm -r "$i"
    done
    popd
    mkdir -p ${BINARY_DIR}
    cp ${tmpdir}/* ${BINARY_DIR}/
  fi
}

function download_gox() {
  if [[ ! -x "$(command -v ${GOX})" ]]; then
    echo "   --> $GOX"
    go get github.com/mitchellh/gox || fatal "go get 'github.com/mitchellh/gox' failed: $?"
  fi
}

function download_goveralls() {
  if [[ -n "$TRAVIS" ]]; then
    if [[ ! -x "$(command -v goveralls)" ]]; then
      echo "   --> goveralls"
      go get github.com/mattn/goveralls || fatal "go get 'github.com/mattn/goveralls' failed: $?"
    fi
  fi
}

function download_upx() {
  if [[ ! -f "$UPX" ]]; then
    verbose "   --> $UPX"
    mkdir -p ${BINARY_DIR}
    curl -L -s -o "${UPX}" ${UPX_URL} || fatal "failed to download '$UPX_URL': $?"
    chmod +x "${UPX}"
  fi
}

function download_binaries() {
  verbose "Fetching binaries..."
  download_go || fatal "failed to download 'go': $?"
  download_glide || fatal "failed to download 'glide': $?"
  download_gox || fatal "failed to download 'gox': $?"
  download_gometalinter || fatal "failed to download 'gometalinter': $?"
  download_goveralls || fatal "failed to download 'goveralls': $?"
  download_upx || fatal "failed to download 'upx': $?"
}

function cleanup() {
  if [[ -z "$TRAVIS" ]]; then
    verbose "Cleanup dist..."
    rm -rf dist/*
  fi
}

function get_go_dependencies() {
  verbose "Getting dependencies..."
  ${GLIDE} install -v || fatal "glide install failed: $?"
}

function install_dependencies() {
  verbose "Installing dependencies..."
  go install ./... || fatal "go install failed: $?"
  go test -i ./... || fatal "go test install failed: $?"
}

function format_source() {
  local gofiles=$(find . -path ./vendor -prune -o -path ./.runtime -prune -o -print | grep '\.go$')

  verbose "Formatting source..."
  if [[ ${#gofiles[@]} -gt 0 ]]; then
    while read -r gofile; do
      gofmt -s -w $PWD/${gofile}
    done <<< "$gofiles"
  fi

  if [[ -n "$TRAVIS" ]] && [[ -n "$(git status --porcelain)" ]]; then
    fatal "Source not formatted"
  fi
}

function lint_source() {
  verbose "Linting source..."
  ${GOMETALINTER} --disable-all --enable=vet --enable=gocyclo --cyclo-over=15 --enable=golint --min-confidence=.85 --enable=ineffassign --skip=Godeps --skip=vendor --skip=third_party --skip=testdata --vendor ./... || fatal "gometalinter failed: $?"
}

function run_tests() {
  verbose "Running tests..."
  if [[ -n "$TRAVIS" ]]; then
    if [[ ! -x "$(command -v goveralls)" ]]; then
      echo "Getting goveralls..."
      go get github.com/mattn/goveralls || fatal "go get 'github.com/mattn/goveralls' failed: $?"
    fi
    goveralls -v -service=travis-ci -ignore=main.go,testutil/server.go,testutil/golden.go || fatal "goveralls: $?"
  else
    go test -v ./... || fatal "$gopackage tests failed: $?"
  fi
}

function build_binaries() {
  local revision=`git rev-parse HEAD`
  local branch=`git rev-parse --abbrev-ref HEAD`
  local host=`hostname`
  local buildDate=`date -u +"%Y-%m-%dT%H:%M:%SZ"`
  go version | grep -q 'go version go1.11.3 ' || fatal "go version is not 1.11.3"

  local xc_arch=${XC_ARCH:-"386 amd64"}
  local xc_os=${XC_OS:-"darwin linux"}
  if [[ -z "$TRAVIS" ]]; then
    xc_arch=$(go env GOARCH)
    xc_os=$(go env GOOS)
  fi

  verbose "Building binaries..."
  ${GOX} -os="${xc_os}" -arch="${xc_arch}" -osarch="!darwin/arm !darwin/arm64" -ldflags "-s -w -X github.com/getbpt/bpt/version.Version=$VERSION -X github.com/getbpt/bpt/version.Revision=$revision -X github.com/getbpt/bpt/version.Branch=$branch -X github.com/getbpt/bpt/version.BuildUser=$USER@$host -X github.com/getbpt/bpt/version.BuildDate=$buildDate" -output="dist/{{.Dir}}-{{.OS}}-{{.Arch}}"  || fatal "gox failed: $?"

  verbose "Rename binaries..."
  for f in dist/*; do
    local n="$f"
    n="${n//darwin/Darwin}"
    n="${n//linux/Linux}"
    n="${n//386/i386}"
    n="${n//amd64/x86_64}"
    mv "$f" "$n" || fatal "failed to rename '$f' to '$n': $?"
  done

  verbose "Compress binaries..."
  for f in dist/*; do
    ${UPX} -9 ${f} || fatal "failed to compress binary '$f': $?"
  done
}

function create_archives() {
  if [[ -n "$TRAVIS" ]]; then
    verbose "Creating archives..."
    pushd dist
    for f in *; do
      local filename=$(basename "$f")
      local extension="${filename##*.}"
      local filename="${filename%.*}"
      if [[ "$filename" != "$extension" ]] && [[ -n "$extension" ]]; then
        extension=".$extension"
      else
        extension=""
      fi
      local archivename="$filename.tar.gz"
      verbose "   --> $archivename"
      local genericname="bpt$extension"
      mv -f "$f" "$genericname"
      tar -czf ${archivename} "$genericname"
      rm -rf "$genericname"
    done
  fi
}

function prepare() {
  cleanup
  download_binaries
  get_go_dependencies
  install_dependencies
}

function build() {
  activate_go
  format_source
  lint_source
  run_tests
}

function package() {
  activate_go
  build_binaries
  create_archives
}
