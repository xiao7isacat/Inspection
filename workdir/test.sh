#!/bin/bash

function has_prometheus() {

  has_prometheus=$(ps -ef |grep prometheus |grep -v grep > /dev/null && echo yes || echo no)
  echo $has_prometheus

}

function has_kubelet() {

  has_kubelet=$(ps -ef |grep kubelet |grep -v grep > /dev/null && echo yes || echo no)
  echo $has_kubelet

}

function has_etcd() {

  has_etcd=$(ps -ef |grep etcd |grep -v grep > /dev/null && echo yes || echo no)
  echo $has_etcd

}


function has_node_exporter() {

  has_node_exporter=$(ps -ef |grep node_exporter |grep -v grep > /dev/null && echo yes || echo no)
  echo $has_node_exporter

}

function has_oom_msg() {

  has_oom_msg=$(dmesg |grep -i oom > /dev/null && echo yes || echo no )
  echo $has_oom_msg

}


function mem_size() {

  mem_size=$(free -m | grep Mem | awk '{print $2}')
  echo $mem_size
}


function cpu_size() {
    cpu_size=&(cat /proc/cpuinfo |grep -c "processor")
}


function main() {
  cpu_size=$(echo "\"cpu_size\": \"$(cpu_size)\"")
  mem_size=$(echo "\"mem_size\": \"$(mem_size)\"")
  tmpBody=$(echo ${cpu_size},${mem_size})
  echo '{'$tmpBody'}'
}

main
