# pmu-emu

## Build / execution

    go build -o pmu-emu && (DATA_PUBLISH_PAUSE_MS="500" DEVICE_ID="15-Zbzvv-09" DATA_FILE="/home/mdye/tmp/_a6_bus1_pmu_merged" ./pmu-emu -logtostderr -v 5 )
