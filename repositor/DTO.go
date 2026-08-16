package repositor

import "time"

type CostInfo struct {
	Name  string
	Price Coal
}
type EquipmentInfo struct {
	Name  string
	Cost  Coal
	IsBuy bool
}

type GameStatisticDTO struct {
	TotalMiners   []MinerInfo
	Balance       Coal
	Equipment     []EquipmentInfo
	IsWin         bool
	CreatedTime   time.Time
	CompletedTime time.Time
}
