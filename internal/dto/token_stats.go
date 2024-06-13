package dto

type Stats struct {
	MemeID              string `json:"memeId"`
	StartDate           string `json:"startDate"`
	VolumeUSD           string `json:"volumeUSD"`
	SwapCount           uint64 `json:"swapCount"`
	CloseUSD            string `json:"closeUSD"`
	TotalValueLockedUSD string `json:"totalValueLockedUSD"`
	HoldersCount        uint64 `json:"holdersCount"`
}

type DailyDiffPct struct {
	VolumeUSD           uint64 `json:"volumeUSD"`
	SwapCount           uint64 `json:"swapCount"`
	CloseUSD            uint64 `json:"closeUSD"`
	TotalValueLockedUSD uint64 `json:"totalValueLockedUSD"`
	HoldersCount        uint64 `json:"holdersCount"`
}
