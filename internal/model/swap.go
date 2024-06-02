package model

type Transaction struct {
	ID          string `json:"id"`
	BlockNumber string `json:"blockNumber"`
	Timestamp   int64  `json:"timestamp"`
	GasUsed     string `json:"gasUsed"`
	GasPrice    string `json:"gasPrice"`
}

type Token struct {
	ID                           string `json:"id"`
	Symbol                       string `json:"symbol"`
	Name                         string `json:"name"`
	Decimals                     string `json:"decimals"`
	TotalSupply                  string `json:"totalSupply"`
	Volume                       string `json:"volume"`
	VolumeUSD                    string `json:"volumeUSD"`
	UntrackedVolumeUSD           string `json:"untrackedVolumeUSD"`
	FeesUSD                      string `json:"feesUSD"`
	TxCount                      string `json:"txCount"`
	PoolCount                    string `json:"poolCount"`
	TotalValueLocked             string `json:"totalValueLocked"`
	TotalValueLockedUSD          string `json:"totalValueLockedUSD"`
	TotalValueLockedUSDUntracked string `json:"totalValueLockedUSDUntracked"`
	DerivedETH                   string `json:"derivedETH"`
}

type Pool struct {
	ID                           string `json:"id"`
	CreatedAtTimestamp           string `json:"createdAtTimestamp"`
	CreatedAtBlockNumber         string `json:"createdAtBlockNumber"`
	Token0                       Token  `json:"token0"`
	Token1                       Token  `json:"token1"`
	FeeTier                      string `json:"feeTier"`
	Liquidity                    string `json:"liquidity"`
	SqrtPrice                    string `json:"sqrtPrice"`
	FeeGrowthGlobal0X128         string `json:"feeGrowthGlobal0X128"`
	FeeGrowthGlobal1X128         string `json:"feeGrowthGlobal1X128"`
	Token0Price                  string `json:"token0Price"`
	Token1Price                  string `json:"token1Price"`
	Tick                         string `json:"tick"`
	ObservationIndex             string `json:"observationIndex"`
	VolumeToken0                 string `json:"volumeToken0"`
	VolumeToken1                 string `json:"volumeToken1"`
	VolumeUSD                    string `json:"volumeUSD"`
	UntrackedVolumeUSD           string `json:"untrackedVolumeUSD"`
	FeesUSD                      string `json:"feesUSD"`
	TxCount                      string `json:"txCount"`
	CollectedFeesToken0          string `json:"collectedFeesToken0"`
	CollectedFeesToken1          string `json:"collectedFeesToken1"`
	CollectedFeesUSD             string `json:"collectedFeesUSD"`
	TotalValueLockedToken0       string `json:"totalValueLockedToken0"`
	TotalValueLockedToken1       string `json:"totalValueLockedToken1"`
	TotalValueLockedETH          string `json:"totalValueLockedETH"`
	TotalValueLockedUSD          string `json:"totalValueLockedUSD"`
	TotalValueLockedUSDUntracked string `json:"totalValueLockedUSDUntracked"`
	LiquidityProviderCount       string `json:"liquidityProviderCount"`
}

type Swap struct {
	ID           string      `json:"id"`
	Transaction  Transaction `json:"transaction"`
	Timestamp    int64       `json:"timestamp"`
	Pool         Pool        `json:"pool"`
	Token0       Token       `json:"token0"`
	Token1       Token       `json:"token1"`
	Sender       string      `json:"sender"`
	Recipient    string      `json:"recipient"`
	Origin       string      `json:"origin"`
	Amount0      string      `json:"amount0"`
	Amount1      string      `json:"amount1"`
	AmountUSD    string      `json:"amountUSD"`
	SqrtPriceX96 string      `json:"sqrtPriceX96"`
	Tick         string      `json:"tick"`
	LogIndex     string      `json:"logIndex"`
}
