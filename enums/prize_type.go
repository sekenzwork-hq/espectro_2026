package enums

type PrizeTypeEnum string

const (
	CashPrize         PrizeTypeEnum = "cash_prize"
	ForeignTripsPrize PrizeTypeEnum = "foreign_trips"
	GadgetsPrize      PrizeTypeEnum = "gadgets"
	OthersPrize       PrizeTypeEnum = "Others"
)

func (p PrizeTypeEnum) IsValid() bool {

	switch p {
	case CashPrize, ForeignTripsPrize, GadgetsPrize, OthersPrize:
		return true
	default:
		return false
	}
}
