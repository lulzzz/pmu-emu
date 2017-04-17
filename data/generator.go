package data

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/michaeldye/pmu-emu/source"

	pmu_server "github.com/michaeldye/synchrophasor-proto/pmu_server"
)

// simpleSynchroDatum implements the SimpleTsDatum interface with data that can be used to fulfill the outgoing protobuf type
type simpleSynchroDatum struct {
	pmu_server.SynchrophasorDatum
	DeviceTs int64
}

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

// NewFileBackedSynchroDatumGenerator instantiates a generator that writes to the returned SimpleTsDatum channel
func NewFileBackedSynchroDatumGenerator(filePath string, deviceID string, datumPublishPauseTime int64) <-chan SimpleTsDatum {

	writer := make(chan SimpleTsDatum) // a blocking channel for fully-formed output Ts data

	publishFn := func(phaseData *pmu_server.SynchrophasorDatum_PhaseData, deviceTs int64) {
		nano := time.Now().UnixNano()
		msg := simpleSynchroDatum{
			SynchrophasorDatum: pmu_server.SynchrophasorDatum{
				Id:        fmt.Sprintf("%v-%v", deviceID, nano),
				Ts:        uint64(nano), // N.B. this is our sampling time
				PhaseData: phaseData,
			},
			DeviceTs: deviceTs,
		}

		writer <- msg

		// TODO: implement a back-off if the errors from reading are severe
		time.Sleep(time.Duration(datumPublishPauseTime) * time.Millisecond)
	}

	reader := source.NewContinuousReader(filePath)
	go func() {

		//generate infinitely
		for {
			if err := reader.ReadDatum(publishFn); err != nil {
				glog.Errorf("Error reading datum: %v", err)
			}
		}
	}()

	return writer
}
