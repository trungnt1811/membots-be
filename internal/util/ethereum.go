package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/params"
)

func WeiToEther(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
}

func EtherToWei(eth *big.Float) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.Ether))
	fracStr := strings.Split(fmt.Sprintf("%.18f", eth), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}

func ConvertWeiToUSD(ethPrice uint64, amountWei string) (string, error) {
	wei := new(big.Int)
	_, ok := wei.SetString(amountWei, 10)
	if !ok {
		return "", fmt.Errorf("invalid wei amount: %s", amountWei)
	}
	ether := WeiToEther(wei)
	ethPriceFloat := new(big.Float).SetUint64(ethPrice)
	usdValue := new(big.Float).Mul(ether, ethPriceFloat)
	return usdValue.Text('f', 4), nil
}

func GetETHPrice() (uint64, error) {
	url := "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2"

	// Define the request payload
	requestBody, err := json.Marshal(map[string]string{
		"query": `{
            bundle(id: "1") {
                ethPrice
            }
        }`,
	})
	if err != nil {
		return 0, fmt.Errorf("error marshaling JSON: %w", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response: %w", err)
	}

	// Define a structure to parse the response
	var response struct {
		Data struct {
			Bundle struct {
				EthPrice string `json:"ethPrice"`
			} `json:"bundle"`
		} `json:"data"`
	}

	// Parse the response
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %w", err)
	}

	// Convert ethPrice string to float64
	ethPriceFloat, err := strconv.ParseFloat(response.Data.Bundle.EthPrice, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing ethPrice to float64: %w", err)
	}

	// Convert float64 to uint64
	ethPriceUint64 := uint64(math.Round(ethPriceFloat))

	return ethPriceUint64, nil
}
