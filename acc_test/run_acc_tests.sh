#!/bin/bash

# All data sources
TF_ACC=1 go test -v -run "DataSource" ./internal/provider
# Resources that do not leave dirty hosts
TF_ACC=1 go test -v -run "User|Password|Autorotate|NetworkPool|Host" ./internal/provider
# Basic compute cluster creation
TF_ACC=1 go test -v -run "TestAccResourceVcfClusterCreate" ./internal/provider -timeout 120m
# Basic workload domain creation
TF_ACC=1 go test -v -run "TestAccResourceVcfDomainCreate" ./internal/provider -timeout 180m