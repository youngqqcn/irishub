package rand

import (
	sdk "github.com/irisnet/irishub/types"
)

// handle all "rand" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgRequestRand:
			return handleMsgRequestRand(ctx, k, msg)
		default:
			return sdk.ErrTxDecode("invalid message parsed in rand module").Result()
		}

		return sdk.ErrTxDecode("invalid message parsed in rand module").Result()
	}
}

// handleMsgRequestRand handles MsgRequestRand
func handleMsgRequestRand(ctx sdk.Context, k Keeper, msg MsgRequestRand) sdk.Result {
	tags, err := k.RequestRand(ctx, msg.Consumer)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// BeginBlocker handles block beginning logic for rand
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) (tags sdk.Tags) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "beginBlock").With("module", "iris/rand"))

	// get data of the last block
	lastBlockHeight := ctx.BlockHeight() - 1
	lastBlockTimestamp := ctx.BlockHeader().GetTime().Unix()
	lastBlockHash := []byte(ctx.BlockHeader().GetLastBlockId().Hash)
	
	// get pending random number requests for lastBloskHeight
	iterator := k.IterateRandRequestQueueByHeight(ctx, lastBlockHeight)
	defer iterator.Close()
	
	handledRandReqNum := 0
	for ; iterator.Valid(); iterator.Next() {
		var reqID string
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &reqID)

		request, err := k.GetRandRequest(ctx, reqID)
		if err != nil {
			continue
		}
		
		// generate a random number
		randNum := MakePRNG(lastBlockHash, lastBlockTimestamp, request.Consumer).GetRand()
		k.SetRand(ctx, NewRand(lastBlockHeight, sdk.NewDec(randNum))
		
		// remove the request
		K.DequeueRandRequest(ctx, lastBlockHeight, reqID)

		handledRandReqNum += 1
	}

	ctx.Logger().Info(fmt.Sprintf("the count of handled rand requests is %d", handledRandReqNum))
}