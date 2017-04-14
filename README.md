# pmu-emu

## Publishing

This project include the make target `publish` that is intended to be executed after a PR has been merged. (Note: this scheme does not have a notion of producing staged development or integration builds, only publishing production stuff. There might be some utility in later producing a `publish-integration` target that is stamped appropriately).

  - Check for an uncommitted files, failing if any exist
  - Clean the project (`make clean`)
  - Build the project (`make all`)
  - Execute all tests (`make test test-integration`)
  - Build a docker container and push it to the repository (`make docker-push`)
  - If the above are successful, tag the `canonical` git repository with the current value in `VERSION`

## Building

#### Considerations

This project requires that you build it from the proper place in your $GOPATH. Also note that it will automatically install `govendor` in your `$GOPATH` when executing `make deps`.

  make

## Execution

Example invocation:

    DATA_PUBLISH_PAUSE_MS="500" DEVICE_ID="15-Zbzvv-09" DATA_FILE="/home/mdye/tmp/_a6_bus1_pmu_merged" pmu-emu -logtostderr -v 5
