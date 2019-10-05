#!/bin/bash
source scripts/_functions.sh

function ctrlC() {
    echo ""
    printDivider
}

trap ctrlC INT

printDivider
build
sleep 10
printDivider
