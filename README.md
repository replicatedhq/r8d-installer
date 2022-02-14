# r8d-installer

## Philosphoy 
1. Minimal dependencies - we only install what we need to set up a cluster and KOTS.
1. Data Replication is an app concern. Small clusters shouldn't have to manage the overhead of distributed CSI providers.
1. One workflow. Always download all dependencies and go.
1. Allow whitelabeling
1. Files not flags
1. Focus on maintainability, not availability.

## Batteries Included
1. RKE2
1. OpenEBS Hostpath
1. KOTS
1. Velero
1. Troubleshoot
1. Host Preflights

## Problems this solves


## Problems this doesn't solve



## TODO:
- [X] Finish online install
- [ ] Move everytyhing to Golang
- [ ] Github Action to collect all dependencies into a tarball
- [ ] Convert online install to work with asset tarball
- [ ] TestGrid or Github Actions for Testing
- [ ] Dependabot or Updates for dependencies in MANIFEST
- [ ] (Future) EKCO replacement
- [ ] (Future) CIS Compliance Action
