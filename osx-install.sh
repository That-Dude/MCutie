#!/bin/bash
# Version 2020-12-23
#
# This script installs mcutie on macOSX
# Tested on BigSur

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)" # /Users/username/scrips/testscript
__file="${__dir}/$(basename "${BASH_SOURCE[0]}")"     # /Users/username/scrips/testscript/scriptname.sh
__base="$(basename ${__file} .sh)"                    # scriptname
__root="$(cd "$(dirname "${__dir}")" && pwd)"         # /Users/username/scrips

# Global Variables
pretty_script_name="MCutie installer"
var_service_name="org.mcutie.com.plist"
var_service_location="$HOME/Library/LaunchAgents"
var_script_install_path="$HOME/.mcutie"
var_executable_name="mcutie"
var_executable_name_path="bin/macos"
var_config_file_name="config.yaml"

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

func_remove()
{
    echo " - Stopping service (if started)"
    launchctl stop "$var_service_name"
    echo " - Unloading service (if loaded)"
    launchctl unload "$var_service_location/$var_service_name" 2>/dev/null
    if func_dir_exisits "$var_script_install_path"; then
        echo " - removing install folder"
        rm -rf "$var_script_install_path"
    fi
    if func_file_exisits "$var_service_location/$var_service_name"; then
        echo " - removing launch control plist file"
        rm "$var_service_location/$var_service_name"
    fi
    echo " - uninstall complete"
}

func_install()
{
    # Check that required files exist
    func_file_check "$var_service_name.template" #template plist file
    func_file_check "$var_executable_name_path/$var_executable_name" #bin/mcutie
    func_file_check "$var_config_file_name" #config.yaml
    
    # Use the temple plist file to create a custom version for the current user
    sed "s/REPLACEME/$USER/g" org.mcutie.com.plist.template > "$var_service_name"
    
    printf " - making sure binary is executable\n"
    chmod +x "$var_executable_name_path/$var_executable_name"
    
    printf " - installing service plist\n"
    cp "$var_service_name" "$var_service_location/$var_service_name"
    
    echo " - create install folder: $var_script_install_path"
    mkdir "$var_script_install_path"
    
    echo " - copying program files..."
    cp "$var_executable_name_path/$var_executable_name" "$var_script_install_path/$var_executable_name"
    cp "$var_config_file_name" "$var_script_install_path/$var_config_file_name"
    
    sleep 1
    launchctl load -w "$var_service_location/$var_service_name" # 2>/dev/null
    sleep 1
    echo " - Starting service"
    launchctl start "$var_service_name"
    echo "done"
}

#  Main code execution starts here
echo "\n*** $pretty_script_name ***"
echo ""

if [ -z "$1" ]; then
    printf "\nTo install type:\n"
    echo "./$__base.sh -i"
    printf "\nTo un-install type:\n"
    echo "./$__base.sh -u"
    echo ""
    exit 1
fi

if [ "$1" == "-i" ]; then
    func_remove
    func_install
fi

if [ "$1" == "-u" ]; then
    func_remove
fi