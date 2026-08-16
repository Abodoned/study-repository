package repositor

import "time"

type MinerType struct {
	Name            string
	Cost            Coal
	Power           int
	CoalPerOnePower Coal
	LoopTime        time.Duration
	Upgrade         Coal
}

var MinerTypesCatalog = map[string]MinerType{
	"small": {
		Name:            "Small",
		Cost:            5,
		Power:           30,
		CoalPerOnePower: 1,
		LoopTime:        3 * time.Second,
	},
	"medium": {
		Name:            "Medium",
		Cost:            50,
		Power:           45,
		CoalPerOnePower: 3,
		LoopTime:        2 * time.Second,
	},
	"strong": {
		Name:            "Strong",
		Cost:            450,
		Power:           60,
		CoalPerOnePower: 10,
		LoopTime:        time.Second,
		Upgrade:         3,
	},
}

func GetMinerType(name string) (MinerType, error) {
	miner, ok := MinerTypesCatalog[name]
	if !ok {
		return MinerType{}, ErrorTypeNotFound
	}
	return miner, nil
}
