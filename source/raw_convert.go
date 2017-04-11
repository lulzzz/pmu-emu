package source

import (
	"errors"
	pmu_server "github.com/michaeldye/synchrophasor-proto/pmu_server"
)

func toPhaseValue(vals []interface{}, arrIdx int) (float64, error) {
	if len(vals) <= arrIdx {
		return 0, errors.New("Illegal index value")
	}

	value := vals[arrIdx].(float64)
	return value, nil
}

func rawToPhaseData(raw map[string]interface{}) (*pmu_server.SynchrophasorDatum_PhaseData, error) {
	d := raw["d"].([]interface{})

	oca, err := toPhaseValue(d, 0)
	if err != nil {
		return nil, err
	}

	ocm, err := toPhaseValue(d, 1)
	if err != nil {
		return nil, err
	}

	tca, err := toPhaseValue(d, 2)
	if err != nil {
		return nil, err
	}

	tcm, err := toPhaseValue(d, 3)
	if err != nil {
		return nil, err
	}

	thca, err := toPhaseValue(d, 4)
	if err != nil {
		return nil, err
	}

	thcm, err := toPhaseValue(d, 5)
	if err != nil {
		return nil, err
	}

	ova, err := toPhaseValue(d, 6)
	if err != nil {
		return nil, err
	}

	ovm, err := toPhaseValue(d, 7)
	if err != nil {
		return nil, err
	}

	tva, err := toPhaseValue(d, 8)
	if err != nil {
		return nil, err
	}

	tvm, err := toPhaseValue(d, 9)
	if err != nil {
		return nil, err
	}

	thva, err := toPhaseValue(d, 10)
	if err != nil {
		return nil, err
	}

	thvm, err := toPhaseValue(d, 11)
	if err != nil {
		return nil, err
	}

	return &pmu_server.SynchrophasorDatum_PhaseData{
		Phase1CurrentAngle:     oca,
		Phase1CurrentMagnitude: ocm,
		Phase2CurrentAngle:     tca,
		Phase2CurrentMagnitude: tcm,
		Phase3CurrentAngle:     thca,
		Phase3CurrentMagnitude: thcm,
		Phase1VoltageAngle:     ova,
		Phase1VoltageMagnitude: ovm,
		Phase2VoltageAngle:     tva,
		Phase2VoltageMagnitude: tvm,
		Phase3VoltageAngle:     thva,
		Phase3VoltageMagnitude: thvm,
	}, nil

}
