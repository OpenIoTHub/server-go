#!/usr/bin/env bash

echo postinstall.sh
systemctl enable server-go
systemctl start server-go