#!/bin/sh
# Krypt global command tests

printf "krypt make\n"
krypt make

printf "\nkrypt make --wordy\n"
krypt make --wordy

printf "\nkrypt list\n"
krypt list

printf "\nkrypt current\n"
krypt current

printf "\nkrypt strength qwerty\n"
krypt strength qwerty

printf "\nkrypt strength extremelypowerfulele1\n"
krypt strength extremelypowerfulele1

printf "\nkrypt version\n"
krypt version

printf "\nkrypt license\n"
krypt license

printf "\nTests Completed (8/8)"
