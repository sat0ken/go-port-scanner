#!/usr/bin/env bash

# rootユーザーが必要
if [ $UID -ne 0 ]; then
  echo "Root privileges are required"
  exit 1;
fi

# 全てのnetnsを削除
ip -all netns delete
ip link del br0 type bridge

# bridgeを作成
ip link add br0 type bridge
ip addr add 192.168.1.1/24 dev br0

# Network Namespaceを作成
ip netns add host0
ip netns add host1
ip netns add attacker

# リンクの作成
ip link add name host0-br0 type veth peer name br0-host0 # host0とbr0のリンク
ip link add name host1-br0 type veth peer name br0-host1 # host1とbr0のリンク
ip link add name attacker-br0 type veth peer name br0-attacker

# ブリッジに接続
ip link set dev br0-host0 master br0
ip link set dev br0-host1 master br0
ip link set dev br0-attacker master br0
ip link set br0-host0 up
ip link set br0-host1 up
ip link set br0-attacker up
ip link set br0 up

# リンクの割り当て
ip link set host0-br0 netns host0
ip link set host1-br0 netns host1
ip link set attacker-br0 netns attacker

# host0のリンクの設定
ip netns exec host0 ip addr add 192.168.1.2/24 dev host0-br0
ip netns exec host0 ip link set host0-br0 up
ip netns exec host0 ethtool -K host0-br0 rx off tx off
ip netns exec host0 ip route add default via 192.168.1.1

# host1のリンクの設定
ip netns exec host1 ip addr add 192.168.1.3/24 dev host1-br0
ip netns exec host1 ip link set host1-br0 up
ip netns exec host1 ethtool -K host1-br0 rx off tx off
ip netns exec host1 ip route add default via 192.168.1.1

# attackerのリンクの設定
ip netns exec attacker ip addr add 192.168.1.4/24 dev attacker-br0
ip netns exec attacker ip link set attacker-br0 up
ip netns exec attacker ethtool -K attacker-br0 rx off tx off
ip netns exec attacker ip route add default via 192.168.1.1
ip netns exec attacker ip link set attacker-br0 address AA:BB:CC:DD:EE:FF

# bridge間でpingが飛ばなかったら
sysctl net.bridge.bridge-nf-call-iptables=0
