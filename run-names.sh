#!/bin/sh

go run main.go -file /home/deploy1/data/names/names.txt -table names -column name -dbconn postgresql://root:secret@localhost:5432/namesdb?sslmode=disable
