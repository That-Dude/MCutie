#!/bin/bash
# updated 29/11/2020

# *** This script un/installs mcutie.

# Set magic variables for current file & dir
__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__file="${__dir}/$(basename "${BASH_SOURCE[0]}")"
__base="$(basename ${__file} .sh)"
__root="$(cd "$(dirname "${__dir}")" && pwd)" # <-- change this as it depends on your app

# *** Global Variables ***
PRETTY_SCRIPT_NAME="MCutie installer"
var_logfile_path="/tmp/mcutie.log"
var_service_name="org.mcutie.com.plist"
var_service_location="$HOME/Library/LaunchAgents"
var_script_install_path="$HOME/.mcutie"
var_executable_name="mcutie"
var_executable_name_path="bin"
var_config_file_name="config.yaml"

func_program_exisits() {
    if ! hash "$1" 2>/dev/null; then
        printf " - $1 not found | You need to install $1 to use this app\n"
        exit 1
    else
        printf " - $1 found\n"
    fi
}

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

# *** Main code execution starts here ***

printf "\n*** MCutie installer ***\n"

func_file_check "$var_service_name.template" #plist
func_file_check "$var_executable_name_path/$var_executable_name" #binary
func_file_check "$var_config_file_name" #config.yaml
func_file_check "osx-uninstall.sh"

sed "s/REPLACEME/$USER/g" org.mcutie.com.plist.template > "$var_service_name"

#echo " - Calling uninstaller to clean up system"
#bash "uninstall.sh"

printf " - making binary executable\n"
chmod +x "$var_executable_name_path/$var_executable_name"

printf " - installing service plist\n"
cp "$var_service_name" "$var_service_location/$var_service_name"

if func_dir_exisits "$var_script_install_path"; then
    echo " - remove exisitng service scripts"
    rm -rf "$var_script_install_path"
fi
echo " - create install folder: $var_script_install_path"
mkdir "$var_script_install_path"

echo " - copying program files..."
cp "$var_executable_name_path/$var_executable_name" "$var_script_install_path/$var_executable_name"
cp "$var_config_file_name" "$var_script_install_path/$var_config_file_name"

echo " - loading service"
launchctl unload "$var_service_location/$var_service_name" # 2>/dev/null
sleep 1
launchctl load -w "$var_service_location/$var_service_name" # 2>/dev/null
sleep 1
echo " - Starting service"
launchctl start "$var_service_name"

echo "done"
