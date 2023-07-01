package entity

type HoldingType uint32

const (
	HoldingTypeDefault HoldingType = iota
	HoldingTypeCustom
)

type Holding struct{}

type Lot struct{}
