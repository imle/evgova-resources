package resources

import (
	"bytes"
)

type IDType int16

type TextID IDType

type ShortString [16]byte

type ControlBitFunction string
type ControlBitTest string

type Credits int32
type FrameCount int16
type TechLevel int16

type Status int16
type HitPoints int32

type AIType int16

const (
	AITypeInherentAI  AIType = iota // Visits planets and runs away when attacked
	AITypeWimpyTrader               // Visits planets and runs away when attacked
	AITypeBraveTrader               // Visits planets and fights back when attacked, but runs away when his attacker is out of range.
	AITypeWarship                   // Seeks out and attacks his enemies, or jumps out if there aren't any.
	AITypeInterceptor               // Seeks out his enemies, or parks in orbit around a planet if he can't find any. Buzzes incoming ships to scan them for illegal cargo. Also acts as "piracy police" by attacking any ship that fires on or attempts to board another, non-enemy ship while the interceptor is watching.
)

type CommodityType int16

const (
	CommodityTypeFood CommodityType = iota
	CommodityTypeIndustrial
	CommodityTypeMedicalSupplies
	CommodityTypeLuxuryGoods
	CommodityTypeMetal
	CommodityTypeEquipment
)

type FlagMask16 uint16
type FlagMask32 uint32
type FlagMask64 uint64

type Resource interface {
	ID() IDType
	Validate() error
}

func byteString(b []byte, length int) string {
	return string(bytes.Runes(b[:bytes.IndexByte(b[:length], 0)]))
}
