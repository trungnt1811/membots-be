package util

import (
	"math"

	"golang.org/x/text/message"
)

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func FormatNotiAmt(number float64) string {
	roundNum := RoundFloat(number, 2)
	p := message.NewPrinter(message.MatchLanguage("en"))
	return p.Sprintf("%.2f", roundNum)
}
