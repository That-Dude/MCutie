#!/bin/bash
# Version 2020-12-23
#
# This script builds mcutie on executable files

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)" # /Users/username/scrips/testscript
__file="${__dir}/$(basename "${BASH_SOURCE[0]}")"     # /Users/username/scrips/testscript/scriptname.sh
__base="$(basename ${__file} .sh)"                    # scriptname
__root="$(cd "$(dirname "${__dir}")" && pwd)"         # /Users/username/scrips

# Global Variables
pretty_script_name="MCutie build"
var_executable_name="mcutie"
var_executable_name_path="bin"


func_file_exisits() { if [ -f "$1" ]; then true; else false; fi }

func_dir_exisits() { if [ -d "$1" ]; then true; else false; fi }

func_file_check()
{
    if ! [ -f "$1" ]; then
        echo "Install file missing: $1"
        echo "Cannot install."
        exit 1
    else
        echo " - $1 found"
    fi
}

#  Main code execution starts here
echo "*** $pretty_script_name ***"
echo ""

if func_dir_exisits "$var_executable_name_path"; then
    echo " - removing old builds folder"
    rm -rf "$var_executable_name_path"
fi

echo " - Creating build folder tree"
mkdir "bin"
mkdir "bin/macos"
mkdir "bin/win64"

# The -ldflags=-s removes the debug info making the executable  smaller
# -H=windowsgui removes the terminal console window on Windows making the agent appear 'servie like'

echo " - Building MacOS"
go build -ldflags="-s -w" -o bin/macos/mcutie mcutie.go

echo " - Building Windows 10 x64"
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o bin/win64/mcutie.exe mcutie.go
#GOOS=windows GOARCH=amd64 go build -o bin/win64/mcutie.exe mcutie.go

#Compress the resulting files with UPX - not currently working with macos binaries under Bigsur
#upx mcutie.exe - windows 10 creates false positives identifying the compressed file as a virus - Nice work windows!
echo " - done."