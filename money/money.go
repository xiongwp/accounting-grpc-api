// Package money 跨服务共享的币种精度 + Money 换算工具。
//
// 约定：
//   - Money.MinorUnits = ISO 最小货币单位（PHP cents / JPY yen / KWD fils）
//   - storage = minor_units × StorageScale（= 100，全币种统一）
//
// 所有需要在 "minor units ↔ storage" 之间换算、或查询币种精度的地方（包括
// order-core 对账 job / e2e 工具 / accounting-system server-side 落账）都从
// 这里读，不要再各自硬编码 100 / 10000 / 1000。
package money

import (
	"fmt"

	accountingv1 "github.com/xiongwp/accounting-grpc-api/gen/accounting/v1"
)

// StorageScale accounting-system 内部 storage 相对 ISO minor units 的放大倍数。
// 跨币种恒定，目的是给未来高精度（加密货币小额手续费等）场景留余量。
const StorageScale int64 = 100

// precisionMinor ISO 4217 最小货币单位相对主单位的小数位。
//   PHP=2 → 1 PHP = 100 cents
//   JPY=0 → 1 JPY = 1 yen（无小数）
//   KWD=3 → 1 KWD = 1000 fils
var precisionMinor = map[string]int{
	"PHP": 2,
	"USD": 2,
	"EUR": 2,
	"GBP": 2,
	"HKD": 2,
	"SGD": 2,
	"AUD": 2,
	"CAD": 2,
	"CHF": 2,
	"MYR": 2,
	"THB": 2,
	"INR": 2,
	"TWD": 2,
	"CNY": 2,
	"VND": 0,
	"IDR": 0,
	"JPY": 0,
	"KRW": 0,
	"KWD": 3,
	"BHD": 3,
	"OMR": 3,
}

// Precision 返回币种的 ISO 小数位。未登记币种返回 err。
func Precision(code string) (int, error) {
	p, ok := precisionMinor[code]
	if !ok {
		return 0, fmt.Errorf("money: unknown currency %q", code)
	}
	return p, nil
}

// MinorUnitsPerMajor 1 个主单位对应的最小单位数。PHP=100, JPY=1, KWD=1000。
func MinorUnitsPerMajor(code string) (int64, error) {
	p, err := Precision(code)
	if err != nil {
		return 0, err
	}
	n := int64(1)
	for i := 0; i < p; i++ {
		n *= 10
	}
	return n, nil
}

// ToStorage minor_units → storage。跨币种都是 ×StorageScale。
func ToStorage(m *accountingv1.Money) (int64, error) {
	if m == nil {
		return 0, fmt.Errorf("money: nil")
	}
	if _, err := Precision(m.GetCurrency()); err != nil {
		return 0, err
	}
	return m.GetMinorUnits() * StorageScale, nil
}

// FromStorage storage → Money{minor_units, currency}。要求 storage % StorageScale == 0，
// 否则视为内部 bug（不应出现非整数 minor）。
func FromStorage(storage int64, code string) (*accountingv1.Money, error) {
	if _, err := Precision(code); err != nil {
		return nil, err
	}
	if storage%StorageScale != 0 {
		return nil, fmt.Errorf("money: storage %d not divisible by scale %d (internal corruption)", storage, StorageScale)
	}
	return &accountingv1.Money{MinorUnits: storage / StorageScale, Currency: code}, nil
}
