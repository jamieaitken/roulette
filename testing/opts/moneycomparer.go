package opts

import (
	"github.com/Rhymond/go-money"
	"github.com/google/go-cmp/cmp"
)

var MoneyComparer = cmp.Comparer(func(x, y money.Money) bool {
	return x.Amount() == y.Amount() && x.Currency() == y.Currency()
})
