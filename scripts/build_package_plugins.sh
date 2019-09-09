#!/bin/bash

# Recursively finds plugin source files and compiles using go build. Expects that the 
# plugin source is self contained in a single file. Allows for multiple plugins to exist 
# at each directory level. Resulting .so file will be at the same location with only the 
# extension changed. 
#
# Param $1 - path of top level plugin directory

for g in $(find $1 -type f -name "*.go"); 
do
    s=`sed 's/.\{3\}$//' <<< $g`.so
    echo building plugin $g as $s
    go build -buildmode=plugin -o $s $g
done