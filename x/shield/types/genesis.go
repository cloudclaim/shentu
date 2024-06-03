package types

import (
	"encoding/json"
	"fmt"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new genesis state.
func NewGenesisState(shieldAdmin sdk.AccAddress, nextPoolID, nextPurchaseID uint64, poolParams PoolParams,
	claimProposalParams ClaimProposalParams, distrParams DistributionParams, totalCollateral, totalWithdrawing, totalShield, totalClaimed math.Int, serviceFees,
	remainingServiceFees sdk.DecCoins, pools []Pool, providers []Provider, purchase []PurchaseList, withdraws []Withdraw,
	lastUpdateTime time.Time, sSRate sdk.Dec, globalStakingPool math.Int, stakingPurchases []ShieldStaking, originalStaking []OriginalStaking,
	proposalIDReimbursementPairs []ProposalIDReimbursementPair) GenesisState {
	return GenesisState{
		ShieldAdmin:                  shieldAdmin.String(),
		NextPoolId:                   nextPoolID,
		NextPurchaseId:               nextPurchaseID,
		PoolParams:                   poolParams,
		ClaimProposalParams:          claimProposalParams,
		DistributionParams:           distrParams,
		TotalCollateral:              totalCollateral,
		TotalWithdrawing:             totalWithdrawing,
		TotalShield:                  totalShield,
		TotalClaimed:                 totalClaimed,
		ServiceFees:                  serviceFees,
		RemainingServiceFees:         remainingServiceFees,
		Pools:                        pools,
		Providers:                    providers,
		PurchaseLists:                purchase,
		Withdraws:                    withdraws,
		LastUpdateTime:               lastUpdateTime,
		ShieldStakingRate:            sSRate,
		GlobalStakingPool:            globalStakingPool,
		StakeForShields:              stakingPurchases,
		OriginalStakings:             originalStaking,
		ProposalIDReimbursementPairs: proposalIDReimbursementPairs,
	}
}

// DefaultGenesisState returns a default genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		NextPoolId:           uint64(1),
		NextPurchaseId:       uint64(1),
		PoolParams:           DefaultPoolParams(),
		ClaimProposalParams:  DefaultClaimProposalParams(),
		DistributionParams:   DefaultDistributionParams(),
		TotalCollateral:      sdk.ZeroInt(),
		TotalWithdrawing:     sdk.ZeroInt(),
		TotalShield:          sdk.ZeroInt(),
		TotalClaimed:         sdk.ZeroInt(),
		ServiceFees:          sdk.DecCoins{},
		RemainingServiceFees: sdk.DecCoins{},
		ShieldStakingRate:    sdk.NewDec(2),
		LastUpdateTime:       time.Now(),
	}
}

// ValidateGenesis validates shield genesis data.
func ValidateGenesis(data GenesisState) error {
	if data.NextPoolId < 1 {
		return fmt.Errorf("failed to validate %s genesis state: NextPoolID must be positive ", ModuleName)
	}
	if err := validatePoolParams(data.PoolParams); err != nil {
		return fmt.Errorf("failed to validate %s pool params: %w", ModuleName, err)
	}
	if err := validateClaimProposalParams(data.ClaimProposalParams); err != nil {
		return fmt.Errorf("failed to validate %s claim proposal params: %w", ModuleName, err)
	}

	return nil
}

// GetGenesisStateFromAppState returns GenesisState given raw application genesis state.
func GetGenesisStateFromAppState(cdc codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}

func NewOriginalStaking(purchaseID uint64, amount math.Int) OriginalStaking {
	return OriginalStaking{
		PurchaseId: purchaseID,
		Amount:     amount,
	}
}
