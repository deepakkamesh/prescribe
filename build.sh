#!/bin/sh
go build
tar -cvzf prescribe.tgz ./prescribe ./resources/*
scp prescribe.tgz dkg@192.168.0.131:~/
rm prescribe.tgz
rm prescribe.log
