#!/usr/bin/env bash

go build 
mv madelyne _example/tests/madelyne
cd _example
go build 
mv example tests/example
cd tests 
time ./madelyne conf.yml 
rm example
rm madelyne
rm access.log