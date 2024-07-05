#!/bin/bash

project_root="."

for dir in $(find "$project_root" -type d); do
    if [[ "$dir" == *.git* ]] || [[ "$dir" == *.vscode* ]]; then
        continue
    fi
    if [[ "$dir" == *dist* ]]; then
        rm -rf "$dir"
    fi

    if [[ "$dir" == *node_modules* ]]; then
        rm -rf "$dir"
    fi
done