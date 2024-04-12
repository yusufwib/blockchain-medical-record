#!/bin/bash

# Compile Go code
go build -o main

# Reload systemd
sudo systemctl daemon-reload

# Stop, start, enable, and check status of the goweb service
sudo systemctl stop goweb
sudo systemctl start goweb
sudo systemctl enable goweb
sudo systemctl status goweb
