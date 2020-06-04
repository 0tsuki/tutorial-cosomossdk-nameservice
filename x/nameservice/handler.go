package nameservice

import (
	"fmt"

	"github.com/0tsuki/tutorial-cosomossdk-nameservice/x/nameservice/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the nameservice type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgSetName:
			return handleMsgSetName(ctx, k, msg)
		case types.MsgBuyName:
			return handleMsgBuyName(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgSetName(ctx sdk.Context, k Keeper, msg types.MsgSetName) (*sdk.Result, error) {
	if !msg.Owner.Equals(k.GetOwner(ctx, msg.Name)) { // Checks if the msg sender is the same as the current owner
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}
	k.SetName(ctx, msg.Name, msg.Value) // If so, set the name to the value specified in the msg.
	return &sdk.Result{}, nil
}

// Handle a message to buy name
func handleMsgBuyName(ctx sdk.Context, k Keeper, msg types.MsgBuyName) (*sdk.Result, error) {
	// Checks if the bid price is greater than the price paid by the current owner
	if k.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid not high enough") // If not, throw an error
	}
	if k.HasOwner(ctx, msg.Name) {
		err := k.CoinKeeper.SendCoins(ctx, msg.Buyer, k.GetOwner(ctx, msg.Name), msg.Bid)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := k.CoinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
		if err != nil {
			return nil, err
		}
	}
	k.SetOwner(ctx, msg.Name, msg.Buyer)
	k.SetPrice(ctx, msg.Name, msg.Bid)
	return &sdk.Result{}, nil
}
