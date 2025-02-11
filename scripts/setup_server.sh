#!/bin/bash
sudo apt update 
sudo apt install -y docker.io
docker build -t cdn .
