package main

import (
	"flag"
	"net"
	"os"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	data "github.com/michaeldye/pmu-emu/sensor_data"
	pmu_server "github.com/michaeldye/synchrophasor-proto/pmu_server"
)

const (
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

	glog.Infof("Setting up gRPC server on %v", defaultBind)

	// Creates a new gRPC server
	s := grpc.NewServer()
	pmu_server.RegisterSynchrophasorDataServer(s, &pmuServerImpl{
		broadcast: data.NewSimpleTsDatumBroadcastWriter(data.NewSimpleSynchroDatumGenerator(os.Getenv("SERIAL"))),
	})
	s.Serve(lis)
}
