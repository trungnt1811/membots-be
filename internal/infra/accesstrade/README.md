# AccessTrade module

This module sync campaigns from AccessTrade API and manage affiliate links

## Main features

- Read AccessTrade campaigns from API
- Update and add new records
- Mark as removed old records

## Core components

- Repository: functions which call APIs and save data to DB
- Usecase: logic functions for sync campaigns and create links
- Worker: provide cronjob configure for sync job
- Test: provide unit tests for this module

## Dependencies

- AccessTrade API Key
- Redis connection
- DB connection

## How to use

Import module then initialize and run worker

```go
import "github.com/astraprotocol/affiliate-system/internal/app/affiliate/accesstrade"

worker := accesstrade.NewAccessTradeWorker(rdc, atUsecase)
worker.RunJob()
```

## Error handling

TODO: write error handling
