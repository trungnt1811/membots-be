import { newMockEvent } from "matchstick-as"
import { ethereum, Address, BigInt } from "@graphprotocol/graph-ts"
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
} from "../generated/TestnetTruglyMemeception/TestnetTruglyMemeception"

export function createCollectFeesEvent(
  memeToken: Address,
  recipient: Address,
  amount0: BigInt,
  amount1: BigInt,
  fee0: BigInt,
  fee1: BigInt
): CollectFees {
  let collectFeesEvent = changetype<CollectFees>(newMockEvent())

  collectFeesEvent.parameters = new Array()

  collectFeesEvent.parameters.push(
    new ethereum.EventParam("memeToken", ethereum.Value.fromAddress(memeToken))
  )
  collectFeesEvent.parameters.push(
    new ethereum.EventParam("recipient", ethereum.Value.fromAddress(recipient))
  )
  collectFeesEvent.parameters.push(
    new ethereum.EventParam(
      "amount0",
      ethereum.Value.fromUnsignedBigInt(amount0)
    )
  )
  collectFeesEvent.parameters.push(
    new ethereum.EventParam(
      "amount1",
      ethereum.Value.fromUnsignedBigInt(amount1)
    )
  )
  collectFeesEvent.parameters.push(
    new ethereum.EventParam("fee0", ethereum.Value.fromUnsignedBigInt(fee0))
  )
  collectFeesEvent.parameters.push(
    new ethereum.EventParam("fee1", ethereum.Value.fromUnsignedBigInt(fee1))
  )

  return collectFeesEvent
}

export function createMeme404CreatedEvent(
  memeToken: Address,
  pool: Address,
  params: ethereum.Tuple,
  tiers: Array<ethereum.Tuple>
): Meme404Created {
  let meme404CreatedEvent = changetype<Meme404Created>(newMockEvent())

  meme404CreatedEvent.parameters = new Array()

  meme404CreatedEvent.parameters.push(
    new ethereum.EventParam("memeToken", ethereum.Value.fromAddress(memeToken))
  )
  meme404CreatedEvent.parameters.push(
    new ethereum.EventParam("pool", ethereum.Value.fromAddress(pool))
  )
  meme404CreatedEvent.parameters.push(
    new ethereum.EventParam("params", ethereum.Value.fromTuple(params))
  )
  meme404CreatedEvent.parameters.push(
    new ethereum.EventParam("tiers", ethereum.Value.fromTupleArray(tiers))
  )

  return meme404CreatedEvent
}

export function createMemeCreatedEvent(
  memeToken: Address,
  pool: Address,
  params: ethereum.Tuple
): MemeCreated {
  let memeCreatedEvent = changetype<MemeCreated>(newMockEvent())

  memeCreatedEvent.parameters = new Array()

  memeCreatedEvent.parameters.push(
    new ethereum.EventParam("memeToken", ethereum.Value.fromAddress(memeToken))
  )
  memeCreatedEvent.parameters.push(
    new ethereum.EventParam("pool", ethereum.Value.fromAddress(pool))
  )
  memeCreatedEvent.parameters.push(
    new ethereum.EventParam("params", ethereum.Value.fromTuple(params))
  )

  return memeCreatedEvent
}

export function createMemeKOLCreatedEvent(
  memeToken: Address,
  pool: Address,
  params: ethereum.Tuple
): MemeKOLCreated {
  let memeKolCreatedEvent = changetype<MemeKOLCreated>(newMockEvent())

  memeKolCreatedEvent.parameters = new Array()

  memeKolCreatedEvent.parameters.push(
    new ethereum.EventParam("memeToken", ethereum.Value.fromAddress(memeToken))
  )
  memeKolCreatedEvent.parameters.push(
    new ethereum.EventParam("pool", ethereum.Value.fromAddress(pool))
  )
  memeKolCreatedEvent.parameters.push(
    new ethereum.EventParam("params", ethereum.Value.fromTuple(params))
  )

  return memeKolCreatedEvent
}

export function createMemeLiquidityAddedEvent(
  memeToken: Address,
  pool: Address,
  amountMeme: BigInt,
  amountETH: BigInt
): MemeLiquidityAdded {
  let memeLiquidityAddedEvent = changetype<MemeLiquidityAdded>(newMockEvent())

  memeLiquidityAddedEvent.parameters = new Array()

  memeLiquidityAddedEvent.parameters.push(
    new ethereum.EventParam("memeToken", ethereum.Value.fromAddress(memeToken))
  )
  memeLiquidityAddedEvent.parameters.push(
    new ethereum.EventParam("pool", ethereum.Value.fromAddress(pool))
  )
  memeLiquidityAddedEvent.parameters.push(
    new ethereum.EventParam(
      "amountMeme",
      ethereum.Value.fromUnsignedBigInt(amountMeme)
    )
  )
  memeLiquidityAddedEvent.parameters.push(
    new ethereum.EventParam(
      "amountETH",
      ethereum.Value.fromUnsignedBigInt(amountETH)
    )
  )

  return memeLiquidityAddedEvent
}

export function createMemecoinBuyEvent(
  memeToken: Address,
  user: Address,
  amountETH: BigInt,
  amountMeme: BigInt
): MemecoinBuy {
  let memecoinBuyEvent = changetype<MemecoinBuy>(newMockEvent())

  memecoinBuyEvent.parameters = new Array()

  memecoinBuyEvent.parameters.push(
    new ethereum.EventParam("memeToken", ethereum.Value.fromAddress(memeToken))
  )
  memecoinBuyEvent.parameters.push(
    new ethereum.EventParam("user", ethereum.Value.fromAddress(user))
  )
  memecoinBuyEvent.parameters.push(
    new ethereum.EventParam(
      "amountETH",
      ethereum.Value.fromUnsignedBigInt(amountETH)
    )
  )
  memecoinBuyEvent.parameters.push(
    new ethereum.EventParam(
      "amountMeme",
      ethereum.Value.fromUnsignedBigInt(amountMeme)
    )
  )

  return memecoinBuyEvent
}

export function createMemecoinExitEvent(
  memeToken: Address,
  user: Address,
  amountETH: BigInt,
  amountMeme: BigInt
): MemecoinExit {
  let memecoinExitEvent = changetype<MemecoinExit>(newMockEvent())

  memecoinExitEvent.parameters = new Array()

  memecoinExitEvent.parameters.push(
    new ethereum.EventParam("memeToken", ethereum.Value.fromAddress(memeToken))
  )
  memecoinExitEvent.parameters.push(
    new ethereum.EventParam("user", ethereum.Value.fromAddress(user))
  )
  memecoinExitEvent.parameters.push(
    new ethereum.EventParam(
      "amountETH",
      ethereum.Value.fromUnsignedBigInt(amountETH)
    )
  )
  memecoinExitEvent.parameters.push(
    new ethereum.EventParam(
      "amountMeme",
      ethereum.Value.fromUnsignedBigInt(amountMeme)
    )
  )

  return memecoinExitEvent
}

export function createOwnershipTransferredEvent(
  user: Address,
  newOwner: Address
): OwnershipTransferred {
  let ownershipTransferredEvent = changetype<OwnershipTransferred>(
    newMockEvent()
  )

  ownershipTransferredEvent.parameters = new Array()

  ownershipTransferredEvent.parameters.push(
    new ethereum.EventParam("user", ethereum.Value.fromAddress(user))
  )
  ownershipTransferredEvent.parameters.push(
    new ethereum.EventParam("newOwner", ethereum.Value.fromAddress(newOwner))
  )

  return ownershipTransferredEvent
}

export function createTreasuryUpdatedEvent(
  oldTreasury: Address,
  newTreasury: Address
): TreasuryUpdated {
  let treasuryUpdatedEvent = changetype<TreasuryUpdated>(newMockEvent())

  treasuryUpdatedEvent.parameters = new Array()

  treasuryUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "oldTreasury",
      ethereum.Value.fromAddress(oldTreasury)
    )
  )
  treasuryUpdatedEvent.parameters.push(
    new ethereum.EventParam(
      "newTreasury",
      ethereum.Value.fromAddress(newTreasury)
    )
  )

  return treasuryUpdatedEvent
}
