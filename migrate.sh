#!/bin/bash
migrate -url "mysql://root:ali@tcp(localhost:3306)/gocart" -path ./sql $*