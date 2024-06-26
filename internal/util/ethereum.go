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
	"time"

	"github.com/avast/retry-go/v4"
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

func WeiStrToEtherStr(wei string) string {
	tmp, _ := big.NewInt(0).SetString(wei, 10)
	result := WeiToEther(tmp)
	etherStr := result.Text('f', 8)
	// Remove trailing zeros and the decimal point if necessary
	etherStr = strings.TrimRight(etherStr, "0")
	etherStr = strings.TrimRight(etherStr, ".")
	return etherStr
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

	var ethPriceUint64 uint64

	err := retry.Do(
		func() error {
			// Define the request payload
			requestBody, err := json.Marshal(map[string]string{
				"query": `{
                    bundle(id: "1") {
                        ethPrice
                    }
                }`,
			})
			if err != nil {
				return fmt.Errorf("error marshaling JSON: %w", err)
			}

			// Create a new HTTP POST request
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
			if err != nil {
				return fmt.Errorf("error creating request: %w", err)
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			// Send the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("error sending request: %w", err)
			}
			defer resp.Body.Close()

			// Check if response status is not OK
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
			}

			// Read the response
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("error reading response: %w", err)
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
				return fmt.Errorf("error unmarshalling response: %w", err)
			}

			// Convert ethPrice string to float64
			ethPriceFloat, err := strconv.ParseFloat(response.Data.Bundle.EthPrice, 64)
			if err != nil {
				return fmt.Errorf("error parsing ethPrice to float64: %w", err)
			}

			// Convert float64 to uint64
			ethPriceUint64 = uint64(math.Round(ethPriceFloat))
			return nil
		},
		retry.Attempts(5),          // Number of retry attempts
		retry.Delay(2*time.Second), // Delay between retry attempts
		retry.OnRetry(func(n uint, err error) {
			fmt.Printf("Retry GetETHPrice attempt %d: %v\n", n+1, err)
		}),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get ETH price after retries: %w", err)
	}

	return ethPriceUint64, nil
}
