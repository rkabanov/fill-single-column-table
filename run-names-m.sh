#!/bin/sh

go run main.go -file /home/deploy1/data/names/names-m.txt -table names_m -column name -dbconn postgresql://root:secret@localhost:5432/namesdb?sslmode=disable
