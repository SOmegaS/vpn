#!/bin/bash

cleanup() {
    echo "Stopping Docker container..."
    docker stop vpn_container
    docker rm vpn_container
    echo "Container stopped."
    exit 0
}

trap cleanup INT TERM

if [[ $1 == "start" ]]; then
    echo "Building Docker image..."
    docker build -f vpn.dockerfile --tag vpn_ubuntu_all_files_include .

    echo "Starting Docker container..."
    docker run --privileged --network host --device /dev/net/tun --name vpn_container -it vpn_ubuntu_all_files_include
    
    docker wait vpn_container
    
    cleanup

else
    echo "Usage: $0 start"
fi

