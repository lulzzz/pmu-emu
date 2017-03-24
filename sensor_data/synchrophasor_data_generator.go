package data

import (
	"fmt"
	pmu_server "github.com/michaeldye/synchrophasor-proto/pmu_server"
	"time"
)

// simpleSynchroDatum implements the SimpleTsDatum interface with data that can be used to fulfill the outgoing protobuf type
type simpleSynchroDatum pmu_server.SynchrophasorDatum

// GetID returns the unique ID of this Datum
func (s simpleSynchroDatum) ID() string {
	return s.Id
}

// Timestamp returns the point in time this datum was recorded
func (s simpleSynchroDatum) Timestamp() uint64 {
	return s.Ts
}

// Datum returns the values in this datum
func (s simpleSynchroDatum) Datum() interface{} {
	return s.PhaseData
}

// TODO: consider making a translator here to the outgoing protobuf type so it can be more easily served via RPC

// SimpleSynchroDatumGenerator generates simpleSynchroDatum type data.
type SimpleSynchroDatumGenerator struct {
	DataWriter chan<- SimpleTsDatum
}

// NewSimpleSynchroDatumGenerator creates a new SimpleSynchroDatumGenerator which writes at an interval to its DataWriter
func NewSimpleSynchroDatumGenerator(deviceID string) <-chan SimpleTsDatum {
	writer := make(chan SimpleTsDatum)

	go func() {

		//generate infinitely
		for {
			nano := time.Now().UnixNano()

			msg := simpleSynchroDatum{
				Id: fmt.Sprintf("%v-%v", deviceID, nano),
				Ts: uint64(nano),
				PhaseData: &pmu_server.SynchrophasorDatum_PhaseData{
					Phase1CurrentAngle: 1,
				},
			}

			writer <- msg
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return writer
}
