# Golang + IoT + OpenCensus Tracing and Reporting

Demo project on how to run a Golang program on an embbeded hardware like Onion Omega 2 or Raspberry Pi. In this case this project works as a basic Weather Station with some sensors and sends the data using OpenCensus library to provide metrics and tracking of the device.

## Prerequisites

Before running the setup script, you will need to install some local deployment tools, configure your GCP projects, and gather your API keys. 

### Create a service account credentials and configure its roles

1. Create the Service Account

    * [Open the IAM & admin page](https://console.cloud.google.com/iam-admin/serviceaccounts)
    * Select **+ Create Service Account** at the top
    * Give it a name (e.g. weather-station)

2. Add the following project roles:

    * Monitoring > Monitoring Metric Writer    

3. Check **Furnish a new private key** and, leave key type as "JSON"

4. Click **Save** and download the file. 

5. Copy the file to `/root/go/key.json` on the device.


### How to build for Omega 

This command will generate a binary file compatible with Omega architecture.

`make build`

### Copy to Omega using rsync

Change your omega address on the Makefile, then run the command: 

`make copy`

### Schematic 

Work in Progress

## References

* https://opencensus.io/articles/iot/