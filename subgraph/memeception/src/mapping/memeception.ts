import {
  CollectFees as CollectFeesEvent,
  Meme404Created as Meme404CreatedEvent,
  MemeCreated as MemeCreatedEvent,
  MemeKOLCreated as MemeKOLCreatedEvent,
  MemeLiquidityAdded as MemeLiquidityAddedEvent,
  MemecoinBuy as MemecoinBuyEvent,
  MemecoinExit as MemecoinExitEvent,
  OwnershipTransferred as OwnershipTransferredEvent,
  TreasuryUpdated as TreasuryUpdatedEvent
} from "../types/Memeception/Memeception"
import {
  CollectFees,
  Meme404Created,
  MemeCreated,
  MemeKOLCreated,
  MemeLiquidityAdded,
  MemecoinBuy,
  MemecoinExit,
  OwnershipTransferred,
  TreasuryUpdated
} from "../types/schema"

export function handleCollectFees(event: CollectFeesEvent): void {
  let entity = new CollectFees(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.memeToken = event.params.memeToken
  entity.recipient = event.params.recipient
  entity.amount0 = event.params.amount0
  entity.amount1 = event.params.amount1
  entity.fee0 = event.params.fee0
  entity.fee1 = event.params.fee1

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleMeme404Created(event: Meme404CreatedEvent): void {
  let entity = new Meme404Created(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.memeToken = event.params.memeToken
  entity.pool = event.params.pool
  entity.params_name = event.params.params.name
  entity.params_symbol = event.params.params.symbol
  entity.params_startAt = event.params.params.startAt
  entity.params_swapFeeBps = event.params.params.swapFeeBps
  entity.params_vestingAllocBps = event.params.params.vestingAllocBps
  entity.params_salt = event.params.params.salt
  entity.params_creator = event.params.params.creator
  entity.params_targetETH = event.params.params.targetETH
  // entity.tiers = event.params.tiers
  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash
  entity.save()
}

export function handleMemeCreated(event: MemeCreatedEvent): void {
  let entity = new MemeCreated(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.memeToken = event.params.memeToken
  entity.pool = event.params.pool
  entity.params_name = event.params.params.name
  entity.params_symbol = event.params.params.symbol
  entity.params_startAt = event.params.params.startAt
  entity.params_swapFeeBps = event.params.params.swapFeeBps
  entity.params_vestingAllocBps = event.params.params.vestingAllocBps
  entity.params_salt = event.params.params.salt
  entity.params_creator = event.params.params.creator
  entity.params_targetETH = event.params.params.targetETH

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleMemeKOLCreated(event: MemeKOLCreatedEvent): void {
  let entity = new MemeKOLCreated(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.memeToken = event.params.memeToken
  entity.pool = event.params.pool
  entity.params_name = event.params.params.name
  entity.params_symbol = event.params.params.symbol
  entity.params_startAt = event.params.params.startAt
  entity.params_swapFeeBps = event.params.params.swapFeeBps
  entity.params_vestingAllocBps = event.params.params.vestingAllocBps
  entity.params_salt = event.params.params.salt
  entity.params_creator = event.params.params.creator
  entity.params_targetETH = event.params.params.targetETH

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleMemeLiquidityAdded(event: MemeLiquidityAddedEvent): void {
  let entity = new MemeLiquidityAdded(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.memeToken = event.params.memeToken
  entity.pool = event.params.pool
  entity.amountMeme = event.params.amountMeme
  entity.amountETH = event.params.amountETH

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleMemecoinBuy(event: MemecoinBuyEvent): void {
  let entity = new MemecoinBuy(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.memeToken = event.params.memeToken
  entity.user = event.params.user
  entity.amountETH = event.params.amountETH
  entity.amountMeme = event.params.amountMeme

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleMemecoinExit(event: MemecoinExitEvent): void {
  let entity = new MemecoinExit(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.memeToken = event.params.memeToken
  entity.user = event.params.user
  entity.amountETH = event.params.amountETH
  entity.amountMeme = event.params.amountMeme

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleOwnershipTransferred(
  event: OwnershipTransferredEvent
): void {
  let entity = new OwnershipTransferred(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.user = event.params.user
  entity.newOwner = event.params.newOwner

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleTreasuryUpdated(event: TreasuryUpdatedEvent): void {
  let entity = new TreasuryUpdated(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.oldTreasury = event.params.oldTreasury
  entity.newTreasury = event.params.newTreasury

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}
