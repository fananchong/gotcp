#!/bin/bash

set -ex

pids=`ps -ux | grep client | grep -v grep  | awk '{print $2}'`
if [[ $pids != "" ]]; then
    kill -9 $pids
fi

cd bin

nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &
nohup ./client > /dev/null 2>&1 &

./server


