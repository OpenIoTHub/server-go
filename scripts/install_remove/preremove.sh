#!/usr/bin/env bash

echo preremove.sh:
systemctl stop server-go
systemctl disable server-go