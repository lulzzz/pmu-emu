package source

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	pmu_server "github.com/michaeldye/synchrophasor-proto/pmu_server"
	"io"
	"os"
)

const (
	keyName = "samples"
)

// ContinuousReader implement a reader that continually outputs data from its source file
type ContinuousReader struct {
	file         *os.File
	jsonFilePath string
	decoder      *json.Decoder
	count        uint64
}

// NewContinuousReader instantiates a ContinuousReader. When all of the source file content has been streamed, it will stream records again from the beginning of the file.
func NewContinuousReader(filePath string) *ContinuousReader {
	glog.V(2).Infof("Instantiated ContinuousReader with filepath: %v", filePath)

	return &ContinuousReader{
		jsonFilePath: filePath,
	}
}

// expected that the desired key is at the top level of the structure
func consumeToKey(decoder *json.Decoder, key string, openingValueDelim json.Delim) (bool, error) {
	// just gonna eat all of the shit in the file until we get to given key and then return

	valFound := false

	for {
		t, err := decoder.Token()
		if err == io.EOF {
			return false, fmt.Errorf("Reached end of content before expected key: '%s'", key)
		}
		if err != nil {
			return false, err
		}

		switch t.(type) {
		case json.Delim:
			if valFound && openingValueDelim == t {
				return decoder.More(), nil
			}

			continue
		case string:
			if t == key {
				if openingValueDelim != ' ' {
					valFound = true
				}
			}
		}
	}
}

func (c *ContinuousReader) init() error {
	if c.file != nil {
		if err := c.file.Close(); err != nil {
			// uh oh
			return err
		}

		glog.V(3).Infof("Closed file: %v", c.jsonFilePath)
		c.decoder = nil
		c.file = nil
	}

	f, err := os.Open(c.jsonFilePath)
	if err != nil {
		// let this fly
		return err
	}

	glog.V(3).Infof("Opened file for reading: %v", c.jsonFilePath)

	c.file = f
	c.decoder = json.NewDecoder(bufio.NewReader(c.file))

	if more, err := consumeToKey(c.decoder, keyName, json.Delim('[')); err != nil {
		return err
	} else if !more {
		return fmt.Errorf("Unable to find content at expected top-level key %v", keyName)
	}

	return nil
}

// ReadDatum reads and publishes Synchrophasor Datum records continuously
func (c *ContinuousReader) ReadDatum(publish func(*pmu_server.SynchrophasorDatum_PhaseData, int64)) error {
	for {
		if c.decoder == nil || !c.decoder.More() {
			// init and bail if trouble
			if err := c.init(); err != nil {
				return err
			}
		}

		for c.decoder.More() {
			var raw map[string]interface{}

			if err := c.decoder.Decode(&raw); err != nil {
				glog.Errorf("Error decoding data record. Err: %v", err)
			} else {

				phase, deviceTs, err := rawToPhaseData(raw)
				if err != nil {
					glog.Errorf("Skipping raw datum %v b/c error occured during parsing. Error: %v", raw, err)
				} else {
					// use callback to publish
					publish(phase, deviceTs)

					glog.V(6).Infof("phase data: %v", phase)
					c.count++
				}
			}
		}

		glog.V(3).Infof("Processed: %v records from file: %v", c.count, c.jsonFilePath)
		// done reading samples, hose this so we can start again
		c.decoder = nil
	}
}
