#!/usr/bin/env bash

# rootユーザーが必要
if [ $UID -ne 0 ]; then
  echo "Root privileges are required"
  exit 1;
fi

# 全てのnetnsを削除
ip -all netns delete

# 4つのnetnsを作成
ip netns add host0
ip netns add host1

# リンクの作成
ip link add name host0-host1 type veth peer name host1-host0

# リンクの割り当て
ip link set host0-host1 netns host0
ip link set host1-host0 netns host1

# host0のリンクの設定
ip netns exec host0 ip addr add 192.168.1.2/24 dev host0-host1
ip netns exec host0 ip link set host0-host1 up
ip netns exec host0 ethtool -K host0-host1 rx off tx off

# host1のリンクの設定
ip netns exec host1 ip addr add 192.168.1.3/24 dev host1-host0
ip netns exec host1 ip link set host1-host0 up
ip netns exec host1 ethtool -K host1-host0 rx off tx off
