package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/evmos/ethermint/x/committer/types"
)

type (
	Keeper struct {
		cdc      	codec.BinaryCodec
		storeKey 	storetypes.StoreKey
		memKey   	storetypes.StoreKey
		paramstore	paramtypes.Subspace
        
		
	}
)

func NewKeeper(
    cdc codec.BinaryCodec,
    storeKey,
    memKey storetypes.StoreKey,
    ps paramtypes.Subspace,
    
    
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
        cdc:      	cdc,
        storeKey: 	storeKey,
        memKey:   	memKey,
        paramstore:	ps,
        
		
	}
}



func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetLastProposal(ctx sdk.Context, proposal types.Proposal) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&proposal)
	store.Set(types.KeyPrefix(types.LastProposalId), b)
}

func (k Keeper) GetLastProposal(ctx sdk.Context) types.Proposal {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.KeyPrefix(types.LastProposalId))
	if b == nil {
		return types.Proposal{Id: 0, EndIndex: 0}
	}
	var lastProposal types.Proposal
	k.cdc.MustUnmarshal(b, &lastProposal)
	return lastProposal
}

func (k Keeper) SetProposal(ctx sdk.Context, proposal types.Proposal) {
	p := types.KeyPrefix(types.ProposalKeyPrefix)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	b := k.cdc.MustMarshal(&proposal)
	store.Set(types.KeyPrefix(fmt.Sprintf("%d", proposal.Id)), b)
}

func (k Keeper) GetProposal(ctx sdk.Context, id uint64) (types.Proposal, bool) {
	p := types.KeyPrefix(types.ProposalKeyPrefix)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	b := store.Get(types.KeyPrefix(fmt.Sprintf("%d", id)))
	if b == nil {
		return types.Proposal{}, false
	}

	var proposal types.Proposal
	k.cdc.MustUnmarshal(b, &proposal)
	return proposal, true
}

// func (k Keeper) AddCommitter(ctx sdk.Context, committer string) {
// 	committers := k.GetAllCommitters(ctx)
// 	committers.CommitterList = append(committers.CommitterList, committer)
// 	k.SetCommitter(ctx, committers)
// }

func (k Keeper) SetCommitter(ctx sdk.Context, committer types.Committer) {
	p := types.KeyPrefix(types.CommitterKeyPrefix)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	b := k.cdc.MustMarshal(&committer)
	store.Set(types.KeyPrefix(types.CommitterKeyPrefix), b)
}

func (k Keeper) GetAllCommitters(ctx sdk.Context) types.Committer {
	p := types.KeyPrefix(types.CommitterKeyPrefix)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	b := store.Get(types.KeyPrefix(types.CommitterKeyPrefix))
	if b == nil {
		return types.Committer{}
	}
	var committers types.Committer
	k.cdc.MustUnmarshal(b, &committers)
	return committers
}

func (k Keeper) IsExistCommitter(ctx sdk.Context, address string) bool {
	committers := k.GetAllCommitters(ctx)
	for _, committer := range committers.CommitterList {
		if committer == address {
			return true
		}
	}
	return false
}

func (k Keeper) NewProposal(ctx sdk.Context) types.Proposal {
	lastProposalId := k.GetLastProposal(ctx).Id;
	proposal := types.Proposal{
		Id: lastProposalId + 1,
	}

	return proposal
}
