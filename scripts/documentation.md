### Build for *native* darwin/arm64 (fast path):
`scripts/build_darwin.sh v0.0.1 arm64`

### Build for x86_64 (Rosetta cross-compile):
You can only do this if on a darwin machine (OSX)
`scripts/build_darwin.sh v0.0.1 x86_64`

## Git Procs

Once the related libraries are built, make sure to tag the git repo. `git tag v0.0.1` will allow the go install toolchain to install with this particular version. eg: `go install github.com/tlaceby/rawgo@v0.0.1`