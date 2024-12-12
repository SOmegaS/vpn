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

    read -p "Сhoose a port to open later: " port
    port=${port:-12345}
    
    export PORT=$port
    echo "Starting Docker container..."
    docker run --cap-add=NET_ADMIN --cap-add=NET_RAW  --network host --device /dev/net/tun -e PORT=$PORT --name vpn_container -it vpn_ubuntu_all_files_include
    
    docker wait vpn_container
    
    cleanup

else
    echo "Usage: $0 start"
fi

