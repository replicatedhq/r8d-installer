# r8d agent

## Purpose 
The purpose of the agent is similar to EKCO in kURL - to provide some automation capabilities outside of the dashboard for embedded clusters.

## High-level Functions
* Host-level Metrics and Alerting (CPU, Memory, Disk, Network Connectivity, etc)
* Auditing the host
* Rotate Certs (RKE2 does this, but use needs to know to restart)
* Automatic KOTS updates
* Take snapshots when none are configured

