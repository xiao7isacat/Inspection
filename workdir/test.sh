#!/bin/bash


function mem_size() {

  mem_size=$(free -m | grep Mem | awk '{print $2}')
  echo $mem_size"g"
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
