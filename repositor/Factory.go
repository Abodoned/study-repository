package repositor

import (
	"context"
	"sync"
	"time"
)

type Factory struct {
	CoalAmount      Coal
	BoughtEquipment map[string]bool
	mtx             sync.RWMutex
	ctx             context.Context
	ctxCancel       context.CancelFunc
	NextId          int
	NowMiners       map[int]*miner
	GameStatistic   GameStatisticDTO
}

func NewFactory() *Factory {
	ctx, cancel := context.WithCancel(context.Background())

	f := &Factory{
		BoughtEquipment: make(map[string]bool),
		ctx:             ctx,
		ctxCancel:       cancel,
		NowMiners:       make(map[int]*miner),
	}

	f.startPassiveMoney()
	f.GameStatistic.CreatedTime = time.Now()
	return f
}

func (f *Factory) startPassiveMoney() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {

			case <-f.ctx.Done():
				return

			case <-ticker.C:
				f.mtx.Lock()
				f.CoalAmount++
				f.mtx.Unlock()
			}
		}
	}()
}

func (f *Factory) BuyMiner(name string) (MinerInfo, error) {
	f.mtx.Lock()
	minerType, err := GetMinerType(name)
	if err != nil {
		f.mtx.Unlock()
		return MinerInfo{}, err

	}
	if minerType.Cost > f.CoalAmount {
		f.mtx.Unlock()
		return MinerInfo{}, ErrorNotEnoughCoal
	}

	id := f.NextId
	f.NextId++

	miner := newMiner(id, minerType)
	f.GameStatistic.TotalMiners = append(f.GameStatistic.TotalMiners, miner.Info())
	f.NowMiners[id] = miner
	ResponseMiner := miner.Info()
	f.CoalAmount = f.CoalAmount - minerType.Cost
	f.mtx.Unlock()
	f.startMiner(miner)
	return ResponseMiner, nil
}
func (f *Factory) startMiner(miner Miner) {
	chn := miner.Run(f.ctx)

	go func() {
		for {
			select {
			case <-f.ctx.Done():
				f.mtx.Lock()
				delete(f.NowMiners, miner.Info().Id)
				f.mtx.Unlock()
				return

			case c, ok := <-chn:
				if !ok {
					f.mtx.Lock()
					delete(f.NowMiners, miner.Info().Id)
					f.mtx.Unlock()
					return
				}
				f.addCoal(c)
			}
		}
	}()
}
func (f *Factory) addCoal(c Coal) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	f.CoalAmount += c
}
func (f *Factory) GetMinerByClass(class string) []MinerInfo {
	f.mtx.RLock()
	defer f.mtx.RUnlock()
	Miners := make([]MinerInfo, 0)
	for _, v := range f.NowMiners {
		if v.Info().Type == class {
			Miners = append(Miners, v.Info())
		}
	}
	return Miners
}

func (f *Factory) GetMinerCost() []CostInfo {
	costs := make([]CostInfo, 0, len(MinerTypesCatalog))
	for _, v := range MinerTypesCatalog {
		temp := CostInfo{
			Name:  v.Name,
			Price: v.Cost,
		}
		costs = append(costs, temp)
	}
	return costs
}
func (f *Factory) GetMinersNow() []MinerInfo {
	minerResponse := []MinerInfo{}
	f.mtx.RLock()
	defer f.mtx.RUnlock()

	for _, v := range f.NowMiners {
		minerResponse = append(minerResponse, v.Info())
	}

	return minerResponse
}
func (f *Factory) GetEquipmentCost() []CostInfo {
	costs := make([]CostInfo, 0, len(EquipmentCatalog))
	for _, v := range EquipmentCatalog {
		costs = append(costs, CostInfo{
			Name:  v.Name,
			Price: v.Cost,
		})
	}
	return costs
}
func (f *Factory) isWin() bool {
	return len(f.BoughtEquipment) == len(EquipmentCatalog)
}
func (f *Factory) BuyEquipment(name string) (EquipmentInfo, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	equipmentType, err := GetEquipmentType(name)
	if err != nil {
		return EquipmentInfo{}, err
	}
	if _, ok := f.BoughtEquipment[equipmentType.Name]; ok {
		return EquipmentInfo{}, ErrorEquipmentAlreadyBuy
	}

	if f.CoalAmount < equipmentType.Cost {
		return EquipmentInfo{}, ErrorNotEnoughCoal
	}

	f.CoalAmount -= equipmentType.Cost
	f.BoughtEquipment[equipmentType.Name] = true
	responseEquipment := EquipmentInfo{
		Name:  equipmentType.Name,
		Cost:  equipmentType.Cost,
		IsBuy: true,
	}
	return responseEquipment, nil
}
func (f *Factory) CheckBalance() Coal {
	f.mtx.RLock()

	defer f.mtx.RUnlock()
	return f.CoalAmount
}
func (f *Factory) CheckEquipment() []EquipmentInfo {
	equipment := make([]EquipmentInfo, 0, len(EquipmentCatalog))
	f.mtx.RLock()
	defer f.mtx.RUnlock()
	for k, v := range EquipmentCatalog {
		if _, ok := f.BoughtEquipment[k]; ok {
			temp := EquipmentInfo{
				Name:  v.Name,
				Cost:  v.Cost,
				IsBuy: true,
			}
			equipment = append(equipment, temp)
		} else {
			temp := EquipmentInfo{
				Name:  v.Name,
				Cost:  v.Cost,
				IsBuy: false,
			}
			equipment = append(equipment, temp)
		}

	}
	return equipment
}

func (f *Factory) StopGame() GameStatisticDTO {
	f.mtx.RLock()
	f.GameStatistic.CompletedTime = time.Now()
	stat := GameStatisticDTO{
		TotalMiners:   f.GameStatistic.TotalMiners,
		Balance:       f.CoalAmount,
		Equipment:     f.CheckEquipment(),
		IsWin:         f.isWin(),
		CreatedTime:   f.GameStatistic.CreatedTime,
		CompletedTime: f.GameStatistic.CompletedTime,
	}
	f.mtx.RUnlock()
	f.ctxCancel()
	return stat
}
