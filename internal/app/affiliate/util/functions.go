package util

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

const IMPORT_ID_PREFIX string = "import_id:"

const (
	UnicodeAccent = `ÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚÝàáâãèéêìíòóôõùúýĂăĐđĨĩŨũƠơƯưẠạẢảẤấẦầẨẩẪẫẬậẮắẰằẲẳẴẵẶặẸẹẺẻẼẽẾếỀềỂểỄễỆệỈỉỊịỌọỎỏỐốỒồỔổỖỗỘộỚớỜờỞởỠỡỢợỤụỦủỨứỪừỬửỮữỰự`
	AsciiAccent   = `AAAAEEEIIOOOOUUYaaaaeeeiioooouuyAaDdIiUuOoUuAaAaAaAaAaAaAaAaAaAaAaAaEeEeEeEeEeEeEeEeIiIiOoOoOoOoOoOoOoOoOoOoOoOoUuUuUuUuUuUuUu`
)

var (
	accentMap map[string]string
)

func init() {
	accentMap = make(map[string]string)

	unicodeList := stringToRune(UnicodeAccent)
	asciiList := stringToRune(AsciiAccent)

	for idx, char := range unicodeList {
		accentMap[char] = asciiList[idx]
	}
}

func IsLowBalance(evmClient *ethclient.Client, address common.Address, threshold float64) (bool, float64, error) {
	var low = false

	walletBallance, err := evmClient.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return low, 0, errors.Errorf("Failed to find balance of account %v", address)
	}
	balance, _ := AastraToASA(walletBallance).Float64()

	if balance < threshold {
		low = true
	}

	return low, balance, nil
}

func RequestIdToImportId(requestId string) int {
	parsed, err := strconv.Atoi(strings.TrimPrefix(requestId, IMPORT_ID_PREFIX))
	if err != nil {
		return 0
	}
	return parsed
}

func ImportIdToRequestId(importId int) string {
	return fmt.Sprintf("%s%d", IMPORT_ID_PREFIX, importId)
}

func ToAsciiString(in string) string {
	out := in

	for _, runeValue := range in {
		char := string(runeValue)

		asciiChar, ok := accentMap[char]
		if ok {
			out = strings.ReplaceAll(out, char, asciiChar)
		}
	}
	return out
}

func stringToRune(s string) []string {

	ll := utf8.RuneCountInString(s)

	var texts = make([]string, ll+1)

	var index = 0

	for _, runeValue := range s {

		texts[index] = string(runeValue)

		index++

	}

	return texts

}
