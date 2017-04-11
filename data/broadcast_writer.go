package data

import (
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"runtime"
	"sync"
)

// writers for data types

// SimpleTsDatumBroadcastWriter writes TS Datum to one or more readers in unbuffered channels (so writer blocks until reader(s) consume the Datum).
type SimpleTsDatumBroadcastWriter struct {
	DataReaders map[string]chan<- SimpleTsDatum
	ReadersSync *sync.Mutex
	DataSource  <-chan SimpleTsDatum
}

// NewSimpleTsDatumBroadcastWriter instantiates a new writer with no readers.
func NewSimpleTsDatumBroadcastWriter(writer <-chan SimpleTsDatum) *SimpleTsDatumBroadcastWriter {
	broadcaster := &SimpleTsDatumBroadcastWriter{
		DataReaders: make(map[string]chan<- SimpleTsDatum),
		ReadersSync: &sync.Mutex{},
		DataSource:  writer,
	}

	// crank it up
	go broadcaster.Broadcast()
	return broadcaster
}

// NewReader creates a new SimpleTsDatum read channel to the broadcast writer and an id and returns them both. This is a threadsafe operation.
func (w *SimpleTsDatumBroadcastWriter) NewReader() (string, <-chan SimpleTsDatum) {
	reader := make(chan SimpleTsDatum)
	id := uuid.New().String()

	glog.Infof("Adding reader: %v", id)
	w.ReadersSync.Lock()
	defer w.ReadersSync.Unlock()

	w.DataReaders[id] = reader
	return id, reader
}

// RemReader adds a new SimpleTsDatum read channel to the broadcast writer. This is a threadsafe operation.
func (w *SimpleTsDatumBroadcastWriter) RemReader(id string) error {
	w.ReadersSync.Lock()
	defer w.ReadersSync.Unlock()

	var exists bool
	if _, exists = w.DataReaders[id]; !exists {
		return fmt.Errorf("Unknown reader id: %v", id)
	}

	delete(w.DataReaders, id)
	// do not close reader from our side
	return nil
}

// TODO: consider that this is long-lived right now; it might be better to make short-lived broadcasters that die after they no longer have readers and close the source channel

// Broadcast reads data and writes each datum to all registered readers. N.B. This is a blocking operation and returns only if an error is encountered.
func (w *SimpleTsDatumBroadcastWriter) Broadcast() error {
	for {

		var d SimpleTsDatum
		var ok bool
		select {
		case d, ok = <-w.DataSource:
			if !ok {
				return errors.New("Data source channel closed.")
			}
		}

		if len(w.DataReaders) == 0 {
			glog.Errorf("Dropping message b/c there aren't any readers. Msg: %v", d)
		} else {

			// write to all readers; will block until they read which means they need to be well-behaved
			for _, reader := range w.DataReaders {
				reader <- d
			}
		}

		runtime.Gosched()
	}

}
