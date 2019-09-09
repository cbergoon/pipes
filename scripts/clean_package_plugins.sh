#!/bin/bash

# Recursively finds and deletes all compiled plugin files.
#
# Param $1 - path of top level plugin directory

for s in $(find $1 -type f -name "*.so"); 
do
    echo deleting compiled plugin $s
    rm $s
done