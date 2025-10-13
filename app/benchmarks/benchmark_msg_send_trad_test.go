//go:build benchmarks

package benchmarks_test

import (
	"testing"
	"time"

	testutil "github.com/celestiaorg/celestia-app/v6/test/util"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/require"
)

const blockTime = time.Duration(6 * time.Second)

func TradBenchmarkCheckTx_MsgSend_1(b *testing.B) {
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
	for b.Loop() {
		resp, err := testApp.CheckTx(&checkTxReq)
		require.NoError(b, err)
	}

}

func TradBenchmarkCheckTx_MsgSend_8MB(b *testing.B) {
	testApp, rawTxs := generateMsgSendTransactions(b, 31645)

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

	var totalGas int64
	for _, tx := range rawTxs {
		checkTxRequest := types.RequestCheckTx{
			Tx:   tx,
			Type: types.CheckTxType_New,
		}
		b.ReportAllocs()
		for b.Loop() {
			resp, err := testApp.CheckTx(&checkTxRequest)
			require.NoError(b, err)
		}
	}
}

func TradBenchmarkFinalizeBlock_MsgSend_1(b *testing.B) {
	testApp, rawTxs := generateMsgSendTransactions(b, 1)

	finalizeBlockReq := types.RequestFinalizeBlock{
		Time:   testutil.GenesisTime.Add(blockTime),
		Height: testApp.LastBlockHeight() + 1,
		Hash:   testApp.LastCommitID().Hash,
		Txs:    rawTxs,
	}
	b.ReportAllocs()
	for b.Loop() {
		resp, err := testApp.FinalizeBlock(&finalizeBlockReq)
		require.NoError(b, err)
	}
}

func TradBenchmarkFinalizeBlock_MsgSend_8MB(b *testing.B) {
	testApp, rawTxs := generateMsgSendTransactions(b, 31645)

	finalizeBlockReq := types.RequestFinalizeBlock{
		Time:   testutil.GenesisTime.Add(blockTime),
		Height: testApp.LastBlockHeight() + 1,
		Hash:   testApp.LastCommitID().Hash,
		Txs:    rawTxs,
	}
	b.ReportAllocs()
	for b.Loop() {
		resp, err := testApp.FinalizeBlock(&finalizeBlockReq)
		require.NoError(b, err)
	}

	var totalGas int64
	for i := range rawTxs {
		require.Equal(b, uint32(0), resp.TxResults[i].Code)
		require.Equal(b, "", resp.TxResults[i].Codespace)
		totalGas += resp.TxResults[i].GasUsed
	}
}

func TradBenchmarkPrepareProposal_MsgSend_1(b *testing.B) {
	testApp, rawTxs := generateMsgSendTransactions(b, 1)

	prepareProposalReq := types.RequestPrepareProposal{
		Txs:    rawTxs,
		Height: testApp.LastBlockHeight() + 1,
	}
	b.ReportAllocs()
	for b.Loop() {
		resp, err := testApp.PrepareProposal(&prepareProposalReq)
		require.NoError(b, err)
	}
}

func TradBenchmarkPrepareProposal_MsgSend_8MB(b *testing.B) {
	// a full 8mb block equals to around 31645 msg send transactions.
	// using 31645 to let prepare proposal choose the maximum
	testApp, rawTxs := generateMsgSendTransactions(b, 31645)

	prepareProposalReq := types.RequestPrepareProposal{
		Txs:    rawTxs,
		Height: testApp.LastBlockHeight() + 1,
	}
	b.ReportAllocs()
	for b.Loop() {
		resp, err := testApp.PrepareProposal(&prepareProposalReq)
		require.NoError(b, err)
	}
}

func TradBenchmarkProcessProposal_MsgSend_1(b *testing.B) {
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
	for b.Loop() {
		resp, err := testApp.ProcessProposal(&processProposalReq)
		require.NoError(b, err)
	}
}

func TradBenchmarkProcessProposal_MsgSend_8MB(b *testing.B) {
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
	for b.Loop() {
		resp, err := testApp.ProcessProposal(&processProposalReq)
		require.NoError(b, err)
	}
}
