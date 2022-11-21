package models

import (
	"math"

	"github.com/SimonRichardson/colossus/blist/selectors"
)

const (
	FeeTypeFull selectors.FeeType = "full"
	FeeTypeItem selectors.FeeType = "item"
)

const (
	PercentFeeAmountType selectors.FeeAmountType = "percent"
	FixedFeeAmountType   selectors.FeeAmountType = "fixed"
)

// CalculateCost will calculate the total cost for a item that will include the
// total feels for that ticket. The amount passed to the function is the amount
// of items that each item costs, as there is different fee types to apply for
// each item.
// For raw total cost @see CalculateTotalCost
// For raw fees cost @see CalculateTotalFee
//   - cost Cost : cost contains the amount in pences along with a map of all the
//     fees
//   - amount uint64 : amount is the number of items to apply to cost.
func CalculateCost(cost selectors.Cost, amount uint64) uint64 {
	var (
		fees      = CalculateTotalFee(cost, uint64(amount))
		itemsCost = CalculateTotalCost(cost, uint64(amount))
	)
	return fees + itemsCost
}

// CalculateTotalCost gives you the cost minus the fees. So it's what the
// tickets without any fees applied.
func CalculateTotalCost(cost selectors.Cost, amount uint64) uint64 {
	return cost.Amount * amount
}

// CalculateTotalFee gives you the fees minus the cost. It's useful to see what
// the fees are for the whole thing without the need for the total cost.
func CalculateTotalFee(cost selectors.Cost, amount uint64) uint64 {
	var (
		totalAmount = CalculateTotalCost(cost, amount)
		itemFees    = applyFees(cost.Amount, cost.Fees) * amount
	)
	return applyTotalFees(totalAmount, itemFees, cost.Fees)
}

func applyFees(amount uint64, fees map[selectors.FeeType]selectors.Fee) uint64 {
	if fee, ok := fees[FeeTypeItem]; ok {
		switch fee.Type {
		case PercentFeeAmountType:
			return uint64(math.Ceil(float64(amount) * fee.Percent))
		case FixedFeeAmountType:
			return fee.Fixed
		}
	}
	return 0
}

func applyTotalFees(totalAmount, itemFees uint64, fees map[selectors.FeeType]selectors.Fee) uint64 {
	if fee, ok := fees[FeeTypeFull]; ok {
		switch fee.Type {
		case PercentFeeAmountType:
			return uint64(math.Ceil(float64(totalAmount+itemFees) * fee.Percent))
		case FixedFeeAmountType:
			return itemFees + fee.Fixed
		}
	}
	return itemFees
}
