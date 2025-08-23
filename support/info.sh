#!/usr/bin/env bash
# lib/info.sh – Functions to display banners and installation summary

show_about() {
    printf '%s\n\n' "${_ABOUT:-}"
}

show_banner() {
    printf '%s\n\n' "${_BANNER:-}"
}

show_headers() {
    show_banner || return 1
    show_about || return 1
}

show_summary() {
    local install_dir="$_BINARY"
    local _cmd_executed=
    check_path "$install_dir"
}

export -f show_summary

