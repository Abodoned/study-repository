package repositor

import (
	"context"
	"sync"
	"time"
)

type Coal int

type MinerInfo struct {
	Id        int
	Type      string
	PowerLeft int
}

type miner struct {
	mtx sync.RWMutex

	id                     int
	minerType              MinerType
	powerLeft              int
	currentCoalPerOnePower Coal
}

func newMiner(id int, t MinerType) *miner {

	return &miner{
		id:                     id,
		minerType:              t,
		powerLeft:              t.Power,
		currentCoalPerOnePower: t.CoalPerOnePower,
	}
}

type Miner interface {
	Run(ctx context.Context) <-chan Coal
	Info() MinerInfo
}

func (m *miner) Run(ctx context.Context) <-chan Coal {
	out := make(chan Coal)

	go func() {
		defer close(out)

		ticker := time.NewTicker(m.minerType.LoopTime)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				m.mtx.Lock()

				if m.powerLeft <= 0 {
					m.mtx.Unlock()
					return
				}

				coal := m.currentCoalPerOnePower

				m.powerLeft--
				m.currentCoalPerOnePower += m.minerType.Upgrade

				m.mtx.Unlock()

				select {
				case out <- coal:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out
}
func (m *miner) Info() MinerInfo {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return MinerInfo{
		Id:        m.id,
		Type:      m.minerType.Name,
		PowerLeft: m.powerLeft,
	}
}
