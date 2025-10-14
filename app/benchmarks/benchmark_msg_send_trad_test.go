//go:build Tradbenchmarks

package benchmarks_test

import (
	"testing"
	"time"

	"github.com/celestiaorg/celestia-app/v6/app"
	"github.com/celestiaorg/celestia-app/v6/app/encoding"
	"github.com/celestiaorg/celestia-app/v6/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/v6/pkg/user"
	testutil "github.com/celestiaorg/celestia-app/v6/test/util"
	"github.com/celestiaorg/celestia-app/v6/test/util/testfactory"
	"github.com/celestiaorg/celestia-app/v6/test/util/testnode"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cometbft/cometbft/abci/types"

	"github.com/stretchr/testify/require"
)

const blockTime = time.Duration(6 * time.Second)

func BenchmarkTradCheckTx_MsgSend_1(b *testing.B) {
	testApp, rawTxs := generateMsgSendTransactions(b, 1)

	finalizeBlockResp, err := testApp.FinalizeBlock(&types.RequestFinalizeBlock{
		Time:   testutil.GenesisTime.Add(blockTime),
		Height: testApp.LastBlockHeight() + 1,
		Hash:   testApp.LastCommitID().Hash,
	})
	require.NotNil(b, finalizeBlockResp)
	require.NoError(b, err)

	commitResp, err := testApp.Commit()
	require.NotNil(b, commitResp)
	require.NoError(b, err)

	checkTxReq := types.RequestCheckTx{
		Tx:   rawTxs[0],
		Type: types.CheckTxType_New,
	}

	b.ReportAllocs()
	b.ResetTimer()
    for i := 0; i < b.N; i++ {
		_, err := testApp.CheckTx(&checkTxReq)
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}

}

func BenchmarkTradFinalizeBlock_MsgSend_1(b *testing.B) {
	testApp, rawTxs := generateMsgSendTransactions(b, 1)

	finalizeBlockReq := types.RequestFinalizeBlock{
		Time:   testutil.GenesisTime.Add(blockTime),
		Height: testApp.LastBlockHeight() + 1,
		Hash:   testApp.LastCommitID().Hash,
		Txs:    rawTxs,
	}
	b.ReportAllocs()
	b.ResetTimer()
    for i := 0; i < b.N; i++ {
		_, err := testApp.FinalizeBlock(&finalizeBlockReq)
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}
}

func BenchmarkTradFinalizeBlock_MsgSend_8MB(b *testing.B) {
	testApp, rawTxs := generateMsgSendTransactions(b, 31645)

	finalizeBlockReq := types.RequestFinalizeBlock{
		Time:   testutil.GenesisTime.Add(blockTime),
		Height: testApp.LastBlockHeight() + 1,
		Hash:   testApp.LastCommitID().Hash,
		Txs:    rawTxs,
	}
	b.ReportAllocs()
	b.ResetTimer()
    for i := 0; i < b.N; i++ {
		_, err := testApp.FinalizeBlock(&finalizeBlockReq)
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}

}

func BenchmarkTradPrepareProposal_MsgSend_1(b *testing.B) {
	testApp, rawTxs := generateMsgSendTransactions(b, 1)

	prepareProposalReq := types.RequestPrepareProposal{
		Txs:    rawTxs,
		Height: testApp.LastBlockHeight() + 1,
	}
	b.ReportAllocs()
	b.ResetTimer()
    for i := 0; i < b.N; i++ {
		_, err := testApp.PrepareProposal(&prepareProposalReq)
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}
}

func BenchmarkTradPrepareProposal_MsgSend_8MB(b *testing.B) {
	// a full 8mb block equals to around 31645 msg send transactions.
	// using 31645 to let prepare proposal choose the maximum
	testApp, rawTxs := generateMsgSendTransactions(b, 31645)

	prepareProposalReq := types.RequestPrepareProposal{
		Txs:    rawTxs,
		Height: testApp.LastBlockHeight() + 1,
	}
	b.ReportAllocs()
	b.ResetTimer()
    for i := 0; i < b.N; i++ {
		_, err := testApp.PrepareProposal(&prepareProposalReq)
		b.StartTimer()
		require.NoError(b, err)
		b.StopTimer()
	}
}

func BenchmarkTradProcessProposal_MsgSend_1(b *testing.B) {
	testApp, rawTxs := generateMsgSendTransactions(b, 1)

	prepareProposalReq := types.RequestPrepareProposal{
		Txs:    rawTxs,
		Height: testApp.LastBlockHeight() + 1,
	}

	prepareProposalResp, err := testApp.PrepareProposal(&prepareProposalReq)
	require.NoError(b, err)
	require.Len(b, prepareProposalResp.Txs, 1)

	processProposalReq := types.RequestProcessProposal{
		Txs:          prepareProposalResp.Txs,
		Height:       testApp.LastBlockHeight() + 1,
		DataRootHash: prepareProposalResp.DataRootHash,
		SquareSize:   1,
	}
	b.ReportAllocs()
	b.ResetTimer()
    for i := 0; i < b.N; i++ {
		_, err := testApp.ProcessProposal(&processProposalReq)
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}
}

func BenchmarkTradProcessProposal_MsgSend_8MB(b *testing.B) {
	// a full 8mb block equals to around 31645 msg send transactions.
	// using 31645 to let prepare proposal choose the maximum
	testApp, rawTxs := generateMsgSendTransactions(b, 31645)

	prepareProposalReq := types.RequestPrepareProposal{
		Txs:    rawTxs,
		Height: testApp.LastBlockHeight() + 1,
	}

	prepareProposalResp, err := testApp.PrepareProposal(&prepareProposalReq)
	require.NoError(b, err)
	require.GreaterOrEqual(b, len(prepareProposalResp.Txs), 1)

	processProposalReq := types.RequestProcessProposal{
		Txs:          prepareProposalResp.Txs,
		Height:       testApp.LastBlockHeight() + 1,
		DataRootHash: prepareProposalResp.DataRootHash,
		SquareSize:   128,
	}
	b.ReportAllocs()
	b.ResetTimer()
    for i := 0; i < b.N; i++ {
		_, err := testApp.ProcessProposal(&processProposalReq)
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}
}

// generateMsgSendTransactions creates a test app then generates a number
// of valid msg send transactions.
func generateMsgSendTransactions(b *testing.B, count int) (*app.App, [][]byte) {
	account := "test"
	testApp, kr := testutil.SetupTestAppWithGenesisValSetAndMaxSquareSize(app.DefaultConsensusParams(), 128, account)
	addr := testfactory.GetAddress(kr, account)
	enc := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	acc := testutil.DirectQueryAccount(testApp, addr)
	signer, err := user.NewSigner(kr, enc.TxConfig, testutil.ChainID, user.NewAccount(account, acc.GetAccountNumber(), acc.GetSequence()))
	require.NoError(b, err)
	rawTxs := make([][]byte, 0, count)
	for i := 0; i < count; i++ {
		msg := banktypes.NewMsgSend(
			addr,
			testnode.RandomAddress().(sdk.AccAddress),
			sdk.NewCoins(sdk.NewInt64Coin(appconsts.BondDenom, 10)),
		)
		rawTx, _, err := signer.CreateTx([]sdk.Msg{msg}, user.SetGasLimit(1000000), user.SetFee(10))
		require.NoError(b, err)
		rawTxs = append(rawTxs, rawTx)
		err = signer.IncrementSequence(account)
		require.NoError(b, err)
	}
	return testApp, rawTxs
}

// mebibyte the number of bytes in a mebibyte
const mebibyte = 1048576

// calculateBlockSizeInMb returns the block size in mb given a set
// of raw transactions.
func calculateBlockSizeInMb(txs [][]byte) float64 {
	numberOfBytes := 0
	for _, tx := range txs {
		numberOfBytes += len(tx)
	}
	mb := float64(numberOfBytes) / mebibyte
	return mb
}

// calculateTotalGasUsed simulates the provided transactions and returns the
// total gas used by all of them
func calculateTotalGasUsed(testApp *app.App, txs [][]byte) uint64 {
	var totalGas uint64
	for _, tx := range txs {
		gasInfo, _, _ := testApp.Simulate(tx)
		totalGas += gasInfo.GasUsed
	}
	return totalGas
}
