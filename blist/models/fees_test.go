package models

import (
	"math"
	"math/rand"
	"testing"
	"testing/quick"
	"time"

	"github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

func config() *quick.Config {
	if testing.Short() {
		return &quick.Config{
			MaxCount:      10,
			MaxCountScale: 10,
			Rand:          rand.New(rand.NewSource(time.Now().UnixNano())),
		}
	}
	return &quick.Config{
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func TestCalculateTotalCost(t *testing.T) {
	var (
		f = func(currency string, price, numOfTickets uint64) uint64 {
			return CalculateTotalCost(selectors.Cost{
				Currency: currency,
				Amount:   price,
			}, numOfTickets)
		}
		g = func(currency string, price, numOfTickets uint64) uint64 {
			return price * numOfTickets
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		typex.Fatal(err)
	}
}

func TestCalculateTotalNoFees(t *testing.T) {
	var (
		f = func(currency string, price, numOfTickets uint64) uint64 {
			return CalculateTotalFee(selectors.Cost{
				Currency: currency,
				Amount:   price,
			}, numOfTickets)
		}
		g = func(currency string, price, numOfTickets uint64) uint64 {
			return 0
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		typex.Fatal(err)
	}
}

// Total

func TestCalculateTotalWithFixedTotalFees(t *testing.T) {
	var (
		f = func(currency string, price, numOfTickets, totalFee uint64) uint64 {
			return CalculateTotalFee(selectors.Cost{
				Currency: currency,
				Amount:   price,
				Fees: map[selectors.FeeType]selectors.Fee{
					FeeTypeFull: selectors.Fee{
						Type:  FixedFeeAmountType,
						Fixed: totalFee,
					},
				},
			}, numOfTickets)
		}
		g = func(currency string, price, numOfTickets, totalFee uint64) uint64 {
			return totalFee
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		typex.Fatal(err)
	}
}

func TestCalculateTotalWithPercentTotalFees(t *testing.T) {
	// Manual
	if val := CalculateTotalFee(selectors.Cost{
		Currency: "GBP",
		Amount:   1,
		Fees: map[selectors.FeeType]selectors.Fee{
			FeeTypeFull: selectors.Fee{
				Type:    PercentFeeAmountType,
				Percent: 0.5,
			},
		},
	}, 1); val != 1 {
		t.Errorf("Expected value was incorrect - %d", val)
	}
	if val := CalculateTotalFee(selectors.Cost{
		Currency: "GBP",
		Amount:   1000,
		Fees: map[selectors.FeeType]selectors.Fee{
			FeeTypeFull: selectors.Fee{
				Type:    PercentFeeAmountType,
				Percent: 0.5,
			},
		},
	}, 2); val != 1000 {
		t.Errorf("Expected value was incorrect - %d", val)
	}

	// Quick
	var (
		f = func(currency string, price, numOfTickets uint64, totalFee float64) uint64 {
			return CalculateTotalFee(selectors.Cost{
				Currency: currency,
				Amount:   price,
				Fees: map[selectors.FeeType]selectors.Fee{
					FeeTypeFull: selectors.Fee{
						Type:    PercentFeeAmountType,
						Percent: totalFee,
					},
				},
			}, numOfTickets)
		}
		g = func(currency string, price, numOfTickets uint64, totalFee float64) uint64 {
			return uint64(math.Ceil(float64(price*numOfTickets) * totalFee))
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		typex.Fatal(err)
	}
}

// Item

func TestCalculateItemWithFixedFees(t *testing.T) {
	var (
		f = func(currency string, price, numOfTickets, totalFee uint64) uint64 {
			return CalculateTotalFee(selectors.Cost{
				Currency: currency,
				Amount:   price,
				Fees: map[selectors.FeeType]selectors.Fee{
					FeeTypeItem: selectors.Fee{
						Type:  FixedFeeAmountType,
						Fixed: totalFee,
					},
				},
			}, numOfTickets)
		}
		g = func(currency string, price, numOfTickets, totalFee uint64) uint64 {
			return numOfTickets * totalFee
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		typex.Fatal(err)
	}
}

func TestCalculateItemWithPercentFees(t *testing.T) {
	// Manual
	if val := CalculateTotalFee(selectors.Cost{
		Currency: "GBP",
		Amount:   1000,
		Fees: map[selectors.FeeType]selectors.Fee{
			FeeTypeItem: selectors.Fee{
				Type:    PercentFeeAmountType,
				Percent: 0.5,
			},
		},
	}, 2); val != 1000 {
		t.Errorf("Expected value was incorrect - %d", val)
	}

	// Quick
	var (
		f = func(currency string, price, numOfTickets uint64, totalFee float64) uint64 {
			return CalculateTotalFee(selectors.Cost{
				Currency: currency,
				Amount:   price,
				Fees: map[selectors.FeeType]selectors.Fee{
					FeeTypeItem: selectors.Fee{
						Type:    PercentFeeAmountType,
						Percent: totalFee,
					},
				},
			}, numOfTickets)
		}
		g = func(currency string, price, numOfTickets uint64, totalFee float64) uint64 {
			return uint64(math.Ceil(float64(price)*totalFee)) * numOfTickets
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		typex.Fatal(err)
	}
}

// Total + Item

func TestCalculateTotalAndItemWithFixedFees(t *testing.T) {
	var (
		f = func(currency string, price, numOfTickets, totalFee, itemFee uint64) uint64 {
			return CalculateTotalFee(selectors.Cost{
				Currency: currency,
				Amount:   price,
				Fees: map[selectors.FeeType]selectors.Fee{
					FeeTypeFull: selectors.Fee{
						Type:  FixedFeeAmountType,
						Fixed: totalFee,
					},
					FeeTypeItem: selectors.Fee{
						Type:  FixedFeeAmountType,
						Fixed: itemFee,
					},
				},
			}, numOfTickets)
		}
		g = func(currency string, price, numOfTickets, totalFee, itemFee uint64) uint64 {
			return (numOfTickets * itemFee) + totalFee
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		typex.Fatal(err)
	}
}

func TestCalculateTotalAndItemWithPercentFees(t *testing.T) {
	// Manual
	if val := CalculateTotalFee(selectors.Cost{
		Currency: "GBP",
		Amount:   1000,
		Fees: map[selectors.FeeType]selectors.Fee{
			FeeTypeFull: selectors.Fee{
				Type:    PercentFeeAmountType,
				Percent: 0.5,
			},
			FeeTypeItem: selectors.Fee{
				Type:    PercentFeeAmountType,
				Percent: 0.5,
			},
		},
	}, 2); val != 1500 {
		t.Errorf("Expected value was incorrect - %d", val)
	}

	// Quick
	var (
		f = func(currency string, price, numOfTickets uint64, itemFee, totalFee float64) uint64 {
			return CalculateTotalFee(selectors.Cost{
				Currency: currency,
				Amount:   price,
				Fees: map[selectors.FeeType]selectors.Fee{
					FeeTypeFull: selectors.Fee{
						Type:    PercentFeeAmountType,
						Percent: totalFee,
					},
					FeeTypeItem: selectors.Fee{
						Type:    PercentFeeAmountType,
						Percent: itemFee,
					},
				},
			}, numOfTickets)
		}
		g = func(currency string, price, numOfTickets uint64, itemFee, totalFee float64) uint64 {
			var (
				totalAmount = price * numOfTickets
				itemFees    = uint64(math.Ceil(float64(price)*itemFee)) * numOfTickets
			)
			return uint64(math.Ceil(float64(totalAmount+itemFees) * totalFee))
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		typex.Fatal(err)
	}
}
