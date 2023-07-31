package ranges

import (
	"errors"
	"fmt"
	"searchbin/interfaces"
	"searchbin/tree"
)

var (
	ErrorEmptyBankCode = errors.New("bank code is empty")
	ErrorInvalidRange  = errors.New("low is greater than or equal to high")
)

type RangeBin struct {
	Low       uint64             `json:"low"`
	High      uint64             `json:"high"`
	BankCode  string             `json:"bank_code"`
	RU        string             `json:"-"`
	EN        string             `json:"-"`
	RangeBig  []interfaces.Range `json:"range_big,omitempty"`
	subRanges *tree.Tree
}

func NewRange(BankCode string, Low uint64, High uint64) (*RangeBin, error) {
	/* if BankCode == "" {[DoRequestV2]
		return nil, ErrorEmptyBankCode
	} */

	if Low >= High {
		return nil, ErrorInvalidRange
	}

	return &RangeBin{
		Low:       Low,
		High:      High,
		BankCode:  BankCode,
		RangeBig:  make([]interfaces.Range, 0),
		subRanges: tree.NewTree(),
	}, nil
}

func (b *RangeBin) Equal(other interfaces.Range) bool {
	// Тот же диапазон если хотя бы одна из границ входит в уже существующий
	if (b.Low <= other.GetLow() && other.GetLow() <= b.High) || (b.Low <= other.GetHigh() && other.GetHigh() <= b.High) {
		b.subRanges.Insert(other)
		return true
	} else if other.GetLow() <= b.Low && b.High <= other.GetHigh() {
		b.RangeBig = append(b.RangeBig, other)
		return true
	}
	return false
}

func (b *RangeBin) Contains(bin uint64) interfaces.Range {

	// Нужный диапазон если находиться внутри границ
	if b.Low <= bin && bin <= b.High {
		if b.subRanges.Hight() != 0 {
			result := b.subRanges.Find(bin)
			if result != nil {
				return result
			}
		}

		return b
	}

	return nil
}

func (b RangeBin) RangeBelow(other interfaces.Range) bool {
	// Ниже если обе границы меньше нижней границы текущего диапазона
	return other.GetLow() < b.Low && other.GetHigh() < b.Low
}

func (b RangeBin) BinBelow(bin uint64) bool {
	// Ниже если бин меньше нижней границы текущего диапазона
	return bin < b.Low
}

func (b RangeBin) RangeHigher(other interfaces.Range) bool {
	// Выше если обе границы больше верхней границы текущего диапазона
	return other.GetLow() > b.High && other.GetHigh() > b.High
}

func (b RangeBin) BinHigher(bin uint64) bool {
	// Выше если бин больше верхней границы текущего диапазона
	return bin > b.High
}

func (b RangeBin) GetLow() uint64 {
	return b.Low
}

func (b RangeBin) GetHigh() uint64 {
	return b.High
}

func (b RangeBin) Code() string {
	return b.BankCode
}

func (b RangeBin) GetEN() string {
	return b.EN
}
func (b RangeBin) GetRU() string {
	return b.RU
}

func (b RangeBin) String() string {
	return fmt.Sprintf("{\n\t\"bank_code\": \"%s\",\n\t\"range_low\": %d,\n\t\"range_high\": %d,\n\t\"range_big\":%v\n}", b.BankCode, b.Low, b.High, b.RangeBig)
}
