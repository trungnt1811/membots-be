package util

import (
	"crypto/ecdsa"
	"errors"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

// UserAmountToWei converts decimal user friendly representation of token amount to 'Wei' representation with provided amount of decimal places
// eg UserAmountToWei(1, 5) => 100000
func UserAmountToWei(amount string, decimal *big.Int) (*big.Int, error) {
	amountFloat, ok := big.NewFloat(0).SetString(amount)
	if !ok {
		return nil, errors.New("wrong amount format")
	}
	ethValueFloat := new(big.Float).Mul(amountFloat, big.NewFloat(math.Pow10(int(decimal.Int64()))))
	ethValueFloatString := strings.Split(ethValueFloat.Text('f', int(decimal.Int64())), ".")

	i, ok := big.NewInt(0).SetString(ethValueFloatString[0], 10)
	if !ok {
		return nil, errors.New(ethValueFloat.Text('f', int(decimal.Int64())))
	}

	return i, nil
}

// Astra amount to Wei
func AmountToWei(amount float64, decimal *big.Int) (*big.Int, error) {
	amountFloat := big.NewFloat(amount)
	ethValueFloat := new(big.Float).Mul(amountFloat, big.NewFloat(math.Pow10(int(decimal.Int64()))))
	ethValueFloatString := strings.Split(ethValueFloat.Text('f', int(decimal.Int64())), ".")

	i, ok := big.NewInt(0).SetString(ethValueFloatString[0], 10)
	if !ok {
		return nil, errors.New(ethValueFloat.Text('f', int(decimal.Int64())))
	}

	return i, nil
}

func RandomAddress() common.Address {
	now := time.Now().String()

	ranInt := rand.Intn(1000000)

	hash := crypto.Keccak256([]byte("LOCK_REWARD_HOLDER" + now + strconv.Itoa(ranInt)))

	address := hexutil.Encode(hash[12:])

	return common.HexToAddress(address)
}

func AastraToASA(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.Ether))
}

func PrivateKeyToAddr(privKey string) common.Address {
	pk, err := crypto.HexToECDSA(privKey)
	if err != nil {
		log.LG.Fatalf("failed convert string to ecdsa pk %v", err)
	}

	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.LG.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
}
