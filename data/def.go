package data

// generic ts sensor data types

// SimpleTsDatum is a discrete timestamped datum.
type SimpleTsDatum interface {
	ID() string
	Timestamp() uint64
	DeviceTimestamp() float64
	Datum() interface{}
}
