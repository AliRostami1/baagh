#!/bin/sh

# This script copies the content of init_folder to
# systemd_path and reloads systemd configuration so 
# the changes apply

init_folder=../init/systemd/
systemd_path=/etc/systemd/system/

# Copy .service files to sysyemd folder
cp -ar ${init_folder}. $systemd_path 

# Load the configuration to systemd
sudo systemctl daemon-reload