# Heimdall
* A lightweight network traffic analysis tool for an SBC

## User Stories
* A user sees a dashboard/interface showing how much data is consumed by devices in the network.
* The dashboard/interface should detail the devices and what they are what they are doing. (WHAT AND WHO IS CONSUMING WHAT, INTERNAL VISON OF THE NETWORK) 
* A user should just be able to introduce the device and it should work with minimal setup. (Mirror/SPAN port should be set-up)
* An administrator should be able to send a command to update or configure the device remotely.
* An administrator should be able to see a topology of a network (implement later, probably with a trace route).

## Hard Requirements
* Needs to be lightweight as it will run on SBCs
* Allow for remote administration
* Relay a compiled analysis of packets
* Capture Network data
* Capture will be a Ring Buffer
* Remote 

## Packet Scope Definition
* Source and Destination IPs
* Measure traffic volume (Packet Size and Count)
* Protocol classification (TCP, UDP, ICMP)
* Application-layer Protocol Detection (HTTP, DNS, TLS)
* Does not do Deep Packet Inspection, only headers.