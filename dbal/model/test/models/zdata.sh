#!/bin/bash
codeBegin=$codeBegin'''package models

// SchemaFileContents the json shcema files
var SchemaFileContents = map[string][]byte{
'''

codeEnd='''}'''
dirs=$(ls .)
filename="zdata.go"
function content(){
    file=$1
    echo -e "\t\"models/$file\": []byte(\`"
    cat $file
    echo '''`),'''
}
echo "// Package models  $(date)" > $filename
echo "// THIS FILE IS AUTO-GENERATED DO NOT MODIFY MANUALLY" >> $filename
echo "$codeBegin" >> $filename
for d in $dirs; do
    if  echo "$d" | grep '.json'; then
        content $d >> $filename
    fi
done
echo "$codeEnd" >> $filename

