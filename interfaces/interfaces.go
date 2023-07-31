package interfaces

type (
	Range interface {
		Equal(other Range) bool
		Contains(bin uint64) Range
		RangeBelow(other Range) bool
		BinBelow(bin uint64) bool
		RangeHigher(other Range) bool
		BinHigher(bin uint64) bool
		GetLow() uint64
		GetHigh() uint64
		Code() string
		String() string
		GetRU() string
		GetEN() string
	}
)
