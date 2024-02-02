#!/usr/bin/env bash
echo "start set version"
if [ $# -ne 2 ]; then
    echo "Parameter exception, please execute in the format of $0 [directory] [version number]"
    echo "PS：$0 ./ v2.4.0"
    exit 1
fi

if [ ! -d "$1" ]; then
    echo "Error: Directory does not exist"
    exit 1
fi

if [[ "$2" != v* ]]; then
    echo "Error: Version number must start with v"
    exit 1
fi

workdir=.
newVersion=$2
echo "Prepare to replace the GF library version numbers in all go.mod files in the ${workdir} directory with ${newVersion}"

# check find command support or not
output=$(find "${workdir}" -name go.mod 2>&1)
if [[ $? -ne 0 ]]; then
    echo "Error: please use bash or zsh to run!"
    exit 1
fi

if [[ true ]]; then
    echo "package gb" > version.go
    echo "" >> version.go
    echo "const (" >> version.go
    echo -e "\t// VERSION is the current gb version." >> version.go
    echo -e "\tVERSION = \"${newVersion}\"" >> version.go
    echo ")" >> version.go
fi

if [ -f "go.work" ]; then
    mv go.work go.work.version.bak
    echo "Back up the go.work file to avoid affecting the upgrade"
fi

for file in `find ${workdir} -name go.mod`; do
    goModPath=$(dirname $file)
    echo ""
    echo "processing dir: $goModPath"
    cd $goModPath
    go mod tidy
    # Upgrading only gb related libraries, sometimes even if a version number is specified, it may not be possible to successfully upgrade. Please confirm before submitting the code
    go list -f "{{if and (not .Indirect) (not .Main)}}{{.Path}}@${newVersion}{{end}}" -m all | grep "^ghostbb.io/gb"
    go list -f "{{if and (not .Indirect) (not .Main)}}{{.Path}}@${newVersion}{{end}}" -m all | grep "^ghostbb.io/gb" | xargs -L1 go get -v
    go mod tidy
    cd -
done
