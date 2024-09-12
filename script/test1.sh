#!/bin/bash

function default_route() {
  ip route | grep default > /dev/null 2>&1
  if [[ $? -eq 0 ]];then
    default_route=true
  else
    default_route=false
  fi
  echo $default_route
}



function main() {
  default_route=$(echo "\"default_route\": \"$(default_route)\"")
  tmpBody=$(echo ${default_route})
  echo '{'$tmpBody'}'
}

main
