# Golang + IoT + OpenCensus Tracing and Reporting

Demo project on how to run a Golang program on an embbeded hardware like Onion Omega 2 or Raspberry Pi. In this case this project works as a basic Weather Station with some sensors and sends the data using OpenCensus library to provide metrics and tracking of the device.

## Schematic 

Work in Progress

## How to build for Omega 

This command will generate a binary file compatible with Omega architecture.

`make build`

## Copy to Omega

Change your omega address on the Makefile, then run the command: 

`make copy`