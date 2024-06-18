package constant

import "time"

type memeStatus uint

const (
	PROCESSING memeStatus = 0
	SUCCEED    memeStatus = 1
)

type memeceptionStatus uint

const (
	LIVE           memeceptionStatus = 1
	ENDED_SOLD_OUT memeceptionStatus = 2
)

const CacheTimeETHPrice = 10 * time.Second

