package bounty_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	db := dbm.NewMemDB()
	encCdc := shentuapp.MakeEncodingConfig()

	app := shentuapp.NewShentuApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, shentuapp.DefaultNodeHome, 5, encCdc, shentuapp.EmptyAppOptions{})

	genesisState := shentuapp.GenesisStateWithSingleValidator(t, app)
	stateBytes, err := tmjson.Marshal(genesisState)
	require.NoError(t, err)

	app.InitChain(
		abcitypes.RequestInitChain{
			AppStateBytes: stateBytes,
			ChainId:       "test-chain-id",
		},
	)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	acc := app.AccountKeeper.GetModuleAccount(ctx, types.ModuleName)
	require.NotNil(t, acc)
	//acc1 := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.ModuleName))
	//require.NotNil(t, acc1)
}
