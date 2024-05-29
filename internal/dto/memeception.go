package dto

type Meme struct {
	ID                     	uint         			`json:"id"`
	Name                   	string       			`json:"name"`
	Symbol                 	string    	  			`json:"symbol"`
	Description            	string       			`json:"description"`
	TotalSupply            	string       			`json:"totalSupply"`
	Decimals               	uint         			`json:"decimals"`
	LogoUrl                	string      			`json:"logoUrl"`
	BannerUrl              	string      			`json:"bannerUrl"`
	CreatorAddress         	string      			`json:"creatorAddress"`
	ContractAddress        	string       			`json:"contractAddress"`
	SwapFeeBps             	uint         			`json:"swapFeeBps"`
	VestingAllocBps			uint         			`json:"vestingAllocBps"`
	Memerc20              	bool         			`json:"memerc20"`
	Live                   	bool         			`json:"live"`
	NetworkID              	uint         			`json:"networkId"`
	Website                	string         			`json:"website"`
	Memeception            	Memeception    			`json:"memeception"`
	Socials					map[string]Social      	`json:"socials"`
}

type Memeception struct {
	StartAt					uint			`json:"startAt"`
	Status                  uint    		`json:"status"`
	Ama                    	bool    		`json:"ama"`
	ContractAddress			string       	`json:"contractAddress"`
}

type MemeceptionResp struct {
	Meme                  	Meme         	`json:"meme"`
	Price                	uint64    		`json:"price"`
}
