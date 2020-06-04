package keeper

import (
	"github.com/0tsuki/tutorial-cosomossdk-nameservice/x/nameservice/types"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the nameservice Querier
const (
	QueryResolve = "resolve"
	QueryWhois   = "whois"
	QueryNames   = "names"
)

// NewQuerier is the module level router for state queries
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryResolve:
			return queryResolve(ctx, path[1:], req, k)
		case QueryWhois:
			return queryWhois(ctx, path[1:], req, k)
		case QueryNames:
			return queryNames(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown nameservice query endpoint")
		}
	}
}

func queryResolve(ctx sdk.Context, path []string, _ abci.RequestQuery, k Keeper) ([]byte, error) {
	value := k.ResolveName(ctx, path[0])

	if value == "" {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "could not resolve name")
	}

	res, err := codec.MarshalJSONIndent(k.cdc, types.QueryResResolve{Value: value})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryWhois(ctx sdk.Context, path []string, _ abci.RequestQuery, k Keeper) ([]byte, error) {
	whois := k.GetWhois(ctx, path[0])
	res, err := codec.MarshalJSONIndent(k.cdc, whois)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryNames(ctx sdk.Context, _ abci.RequestQuery, k Keeper) ([]byte, error) {
	var nameList types.QueryResNames

	iterator := k.GetNamesIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		nameList = append(nameList, string(iterator.Key()))
	}

	res, err := codec.MarshalJSONIndent(k.cdc, nameList)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
