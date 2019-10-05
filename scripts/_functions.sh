#!/bin/bash
EXAMPLE_DIR="Todos"

function dim() {
    echo -en ""
}

function reset() {
    echo -en "\033[0m"
}

function printDivider() {
    printf '\033[0m\033[38;5;208m%*s\n' "${COLUMNS:-$(tput cols)}" '' | sed "s/ /─/g"
    reset
}

function printDimDivider() {
    printf '\033[0m\033[2m%*s\n' "${COLUMNS:-$(tput cols)}" '' | sed "s/ /─/g"
    reset
}

function header() {
    echo -e "\033[38;5;208m$1"
    reset
}

function section() {
    echo -e "\033[0m\033[1m$1"
    reset
}

function bullet() {
    echo -e "\033[2m- $1"
    reset
}

################################################################################

function build() {
    set -e
    section "Build"
    bullet "Cleaning"
    [ -f "out/$EXAMPLE_DIR" ] && rm "out/$EXAMPLE_DIR"
    bullet "Building"
    go build -o "out/$EXAMPLE_DIR" \
        "examples/$EXAMPLE_DIR/main.go" \
        "examples/$EXAMPLE_DIR/models.go" \
        "examples/$EXAMPLE_DIR/types.go" \
        "examples/$EXAMPLE_DIR/cache_data.go"
}

function run() {
    set -e
    section "Run"
    bullet "Starting out/$EXAMPLE_DIR"

    reset
    printDimDivider
    echo ""
    ENV_FILE="examples/$EXAMPLE_DIR/.env" ./out/$EXAMPLE_DIR
}
