package main

import (
	"flag"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/tap"

	data "github.com/michaeldye/pmu-emu/data"

	pmu_server "github.com/michaeldye/synchrophasor-proto/pmu_server"
)

const (
	deviceIDEnvvarName               = "DEVICE_ID"
	dataFileEnvvarName               = "DATA_FILE"
	dataPublishPauseTimeMSEnvvarName = "DATA_PUBLISH_PAUSE_MS"

	defaultDataPublishPauseTimeMS = 20

	// defaults overridden by envvars
	defaultBind = "0.0.0.0:8008"
)

// pmuServerImpl is an implementation of the protobuf's interface for a PMUServer, an interface for retrieving Synchrophasor data from a PMU.
type pmuServerImpl struct {
	broadcast *data.SimpleTsDatumBroadcastWriter
}

func (s *pmuServerImpl) Sample(samplingFilter *pmu_server.SamplingFilter, stream pmu_server.SynchrophasorData_SampleServer) error {
	id, reader := s.broadcast.NewReader()
	defer s.broadcast.RemReader(id)

	for {
		inter := <-reader

		// translate from our intermediate, generated type to the RPC type
		datum := &pmu_server.SynchrophasorDatum{
			Id:        inter.ID(),
			Ts:        inter.Timestamp(),
			PhaseData: inter.Datum().(*pmu_server.SynchrophasorDatum_PhaseData),
		}

		if err := stream.Send(datum); err != nil {
			return err
		}

		runtime.Gosched()
	}
}

func main() {
	flag.Parse()

	// instantiate a new broadcast writer
	lis, err := net.Listen("tcp", defaultBind)
	if err != nil {
		glog.Fatalf("Failed to listen: %v", err)
		os.Exit(1)
	}

	dataFile := os.Getenv(dataFileEnvvarName)
	if dataFile == "" {
		glog.Fatalf("Unspecified but required envvar %s", dataFileEnvvarName)
		os.Exit(1)
	}
	glog.Infof("Using dataFile %v set by envvar %v", dataFile, dataFileEnvvarName)

	deviceID := os.Getenv(deviceIDEnvvarName)
	if deviceID == "" {
		glog.Fatalf("Unspecified but required envvar %s", deviceIDEnvvarName)
		os.Exit(1)
	}
	glog.Infof("Using deviceID %v set by envvar %v", deviceID, deviceIDEnvvarName)

	var dataPublishPauseTimeMS int64
	if time, err := strconv.ParseInt(os.Getenv(dataPublishPauseTimeMSEnvvarName), 10, 64); err != nil || time < 5 {
		dataPublishPauseTimeMS = defaultDataPublishPauseTimeMS
		glog.Infof("Using default dataPublishPauseTimeMS %v", dataPublishPauseTimeMS)
	} else {
		dataPublishPauseTimeMS = time
		glog.Infof("Using dataPublishPauseTimeMS %v set by envvar %v", dataPublishPauseTimeMS, dataPublishPauseTimeMSEnvvarName)
	}

	glog.Infof("Setting up gRPC server on %v", defaultBind)

	tap := grpc.InTapHandle(func(ctx context.Context, info *tap.Info) (context.Context, error) {
		glog.V(5).Infof("New connection to RPC method: %v", info.FullMethodName)
		glog.V(6).Infof("Context: %v", ctx)

		return ctx, nil
	})

	ep := grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second,
		PermitWithoutStream: true,
	})

	kp := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Second,
		Time:              10 * time.Second,
		Timeout:           10 * time.Second,
	})

	// Creates a new gRPC server
	s := grpc.NewServer(tap, ep, kp)

	pmu_server.RegisterSynchrophasorDataServer(s, &pmuServerImpl{
		broadcast: data.NewSimpleTsDatumBroadcastWriter(data.NewFileBackedSynchroDatumGenerator(dataFile, deviceID, dataPublishPauseTimeMS)),
	})

	// start it
	s.Serve(lis)
}
