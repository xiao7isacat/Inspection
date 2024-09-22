# Inspection

## add script(添加脚本)
./checkctl add script -n test -f test.sh

## add destride(添加基线，希望值)
./checkctl add desired -n test -f test.desired

## add job (添加任务)
./checkctl add job -n test --node_addr="10.211.55.6:8093,10.211.55.7:8093"

## 执行任务
./checkctl run

## 获取状态
./checkctl get status

## 服务端配置文件
node-env-check.yaml
