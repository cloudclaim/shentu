// TODO remove

package crisis

import (
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/crisis/client/cli"
	"github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	"github.com/cosmos/cosmos-sdk/x/crisis/types"

	"github.com/shentufoundation/shentu/v2/common"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// Module init related flags
const (
	FlagSkipGenesisInvariants = "x-crisis-skip-assert-invariants"
)

// AppModuleBasic defines the basic application module used by the crisis module.
type AppModuleBasic struct{}

// Name returns the crisis module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the crisis module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the crisis
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	defaultGenesisState := types.DefaultGenesisState()
	defaultGenesisState.ConstantFee.Denom = common.MicroCTKDenom
	return cdc.MustMarshalJSON(defaultGenesisState)
}

// ValidateGenesis performs genesis state validation for the crisis module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return types.ValidateGenesis(&data)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the capability module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}

// GetTxCmd returns the root tx command for the crisis module.
func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd returns no root query command for the crisis module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command { return nil }

// RegisterInterfaces registers interfaces and implementations of the crisis
// module.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

//____________________________________________________________________________

// AppModule implements an application module for the crisis module.
type AppModule struct {
	AppModuleBasic

	// NOTE: We store a reference to the keeper here so that after a module
	// manager is created, the invariants can be properly registered and
	// executed.
	keeper *keeper.Keeper

	skipGenesisInvariants bool
}

// NewAppModule creates a new AppModule object. If initChainAssertInvariants is set,
// we will call keeper.AssertInvariants during InitGenesis (it may take a significant time)
// - which doesn't impact the chain security unless 66+% of validators have a wrongly
// modified genesis file.
func NewAppModule(keeper *keeper.Keeper, skipGenesisInvariants bool) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,

		skipGenesisInvariants: skipGenesisInvariants,
	}
}

// AddModuleInitFlags implements servertypes.ModuleInitFlags interface.
func AddModuleInitFlags(startCmd *cobra.Command) {
	startCmd.Flags().Bool(FlagSkipGenesisInvariants, false, "Skip x/crisis invariants check on startup")
}

// Name returns the crisis module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants performs a no-op.
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Deprecated: Route returns the message routing key for the crisis module.
func (am AppModule) Route() sdk.Route {
	return sdk.Route{}
}

// QuerierRoute returns no querier route.
func (AppModule) QuerierRoute() string { return "" }

// LegacyQuerierHandler returns no sdk.Querier.
func (AppModule) LegacyQuerierHandler(*codec.LegacyAmino) sdk.Querier { return nil }

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), am.keeper)
}

// InitGenesis performs genesis initialization for the crisis module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	am.keeper.InitGenesis(ctx, &genesisState)
	if !am.skipGenesisInvariants {
		am.keeper.AssertInvariants(ctx)
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the crisis
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (am AppModule) ConsensusVersion() uint64 { return 1 }

// EndBlock returns the end blocker for the crisis module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	crisis.EndBlocker(ctx, *am.keeper)
	return []abci.ValidatorUpdate{}
}
