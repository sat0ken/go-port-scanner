FROM golang:1.20-bullseye

RUN apt update -y && apt upgrade -y \
    && apt install -y libpcap0.8 libpcap-dev iproute2 ethtool iputils-ping arping netcat sudo \
    && apt-get clean && rm -rf /var/lib/apt/lists/*
