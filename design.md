# Design Thoughts

## What kind of assets?
- executables on the path
- tarball (airgap images for rke2 itself)
- docker images
- manifests (local path provisioner)

## What's the best way to pull in the assets required? Can this be abstracted?
-> No, every repo has a different struct that keeps these from being loaded.
- Go CLI?
- Github Actions?
- Script?
- Makefile?

# Do we store dependencies in the repo?
No, this is probably a bad idea, but the build command needs to generate these to include in the overall binary.
