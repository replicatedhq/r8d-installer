# (EXPERIMENTAL) r8d-installer 

## Philosphoy 
1. Minimal dependencies - we only install what we need to set up a cluster and KOTS.
1. Data Replication is an app concern. Small clusters shouldn't have to manage the overhead of distributed CSI providers.
1. One workflow. Always download all dependencies and go.
1. Allow whitelabeling (Support TBD)
1. No system dependencies EVER.
1. Files not flags. Don't intercept configuration for components, like RKE2, in the installer.
1. Focus on maintainability, not availability.

## Batteries Included, No Assembly Required
1. RKE2
1. OpenEBS Hostpath
1. KOTS
1. Velero
1. Troubleshoot
1. Host Preflights

## Problems this solves
- No system dependency conflicts
- Fast
  - small download size
  - fast startup
- Log to output file and also standard out
- Can be bootstaped easily in a VM
- Minimizes the amount of memorizing flags to enter into the cluster install command.
- Allows the output of the installer to be logged to file while still printing to STDOUT
- Doesn't require maintaining KOTS manifest changes in a separate place.
- Consolidates all dependency updates as part of a single command that can be run nightly

## Problems this doesn't solve
- Learning curve for RKE2 support (also a problem with kURL)

## Current Limitations
- Linux AMD64 architecture only
- Only support for default CNI (can be changed later)
- Single Node (can be changed later)

# Metrics 
- New install console up time
- Baseline CPU, Memory
- Airgap size

## MVP Goals
- [ ] Airgap installs
- [ ] Airgap Updates
- [ ] Host Preflights
- [ ] 100% Update Automation

## Non-goals
- [ ] Test Coverage
## TODO:
- [X] Finish online install
- [ ] Golang command to collect airgap resource
- [ ] Golang command to install the cluster
- [ ] Host Preflights 
- [ ] Troubleshoot bundling
- [ ] KOTS-lite Image (remove unnecessary kubectl binaries)
- [ ] Golang command for update.
- [ ] Github Action for dependency update
- [ ] Enforce conventional commits
- [ ] Golang command for upgrades
- [ ] TestGrid or Github Actions for Testing
- [ ] (Future) Developer Cache - don't download dependencies that are already installed 
- [ ] (Future) EKCO replacement (r8d-installer agent)
- [ ] (Future) CIS Compliance Action
- [ ] (Future) Multi-node support

## Shout Outs
1. Kira Boyle for initial design discussions
1. John Murphy for idea to use go tooling to manage dependency updates
