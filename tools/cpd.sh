#!/usr/bin/env bash
# cpd - copy / paste detector

tmp=/tmp/cpd.$$
dirs=$(find backend core frontend -name java -type d | grep /src/main/java)

# shellcheck disable=SC2086
find $dirs -name \*.java > $tmp

pmd cpd --debug --file-list $tmp --format text --minimum-tokens 100 |tee /tmp/cpd.txt

rm -f $tmp
