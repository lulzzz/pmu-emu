# pmu-emu

A Power Management Unit (PMU) Emulator mostly used for execution in a Horizon environment. The system comprises a Golang binary and a JSON-formatted data file. The `pmu-emu` runs a gRPC server (on port `9009`) that continuously streams synchrophasor data to clients connected to this gRPC server.

Note: both below execution examples assume you've stored a sample data file (like https://tor01.objectstorage.softlayer.net/v1/AUTH_bd05f276-e42f-4fa1-b7b3-780e8544769f/pmu-emu-data/_a6_bus1_pmu_merged?temp_url_sig=5cf9855857580c96dd903a013217206abaceec81&temp_url_expires=30001492637154) in `/tmp/_a6_bus1_pmu_merged`.

### Example native invocation

    DATA_PUBLISH_PAUSE_MS="500" DEVICE_ID="15-Zbzvv-09" DATA_FILE="/tmp/_a6_bus1_pmu_merged" pmu-emu -logtostderr -v 5

### Example Docker container invocation

    docker run --rm --name pmu-emu -p 127.0.0.1:8008:8008/tcp -e "DATA_PUBLISH_PAUSE_MS=500" -e "DEVICE_ID=15-Zbzvv-09" -e "DATA_FILE=/tmp/_a6_bus1_pmu_merged" -v /tmp/a6_bus1_pmu_merged:/tmp/a6_bus1_pmu_merged:ro -t summit.hovitos.engineering/$(uname -m)/pmu-emu:latest

## Related Projects

 * `synchrophasor-proto` (https://github.com/michaeldye/synchrophasor-proto): The protocol specifications for all synchrophasor data projects
 * `synchrophasor-publisher` (https://github.com/michaeldye/synchrophasor-publisher): A client that connects to the `pmu-emu`s gRPC server, processes data it gathers, and then publishes it to a gRPC ingest Data Processing Engine (DPE), an instances of `synchrophasor-dpe`
 * `synchrophasor-dpe` (https://github.com/michaeldye/synchrophasor-dpe): A DPE data ingest server that is connected-to by `synchrophasor-publisher` clients

## Development

### Environment setup

 * Install `make`
 * Install Golang v.1.7.x or newer, set up an appropriate `$GOPATH`, etc. (cf. https://golang.org/doc/install)
 * Install `protoc`, the Google protobuf compiler (cf. instructions at https://github.com/michaeldye/synchrophasor-proto)
 * Install Docker Community Edition version 17.04.0-ce or newer (cf. https://www.docker.com/community-edition#/download or use https://get.docker.com/)

## Building

### Considerations

This project requires that you build it from the proper place in your `$GOPATH`. Also note that it will automatically install `govendor` in your `$GOPATH` when executing `make deps`.

### Compiling the executable

    make

### Creating a Docker execution container

    make docker

### Publishing

This project include the make target `publish` that is intended to be executed after a PR has been merged. (Note: this scheme does not have a notion of producing staged development or integration builds, only publishing production stuff. There might be some utility in later producing a `publish-integration` target that is stamped appropriately).

  - Check for an uncommitted files, failing if any exist
  - Clean the project (`make clean`)
  - Build the project (`make all`)
  - Execute all tests (`make test test-integration`)
  - Build a docker container and push it to the repository (`make docker-push`)
  - If the above are successful, tag the `canonical` git repository with the current value in `VERSION`
