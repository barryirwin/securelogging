# Build

Folder that contains the build script, and when run, will put the different binaries built in this folder.

Also contains the default config files and keys for both the client and the server, move them to the desired hosts.

*Only the public comms key is required by the client*, the rest is used by the server. *DO NOT* give the client all the keys!

## Dependencies

On the first build, `go mod` should take care of it. If for some reason it fails, run:

```bash
go get $missing_dependency
```
