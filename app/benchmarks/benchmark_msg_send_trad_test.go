//go:build benchmarks

package benchmarks_test

import (
	"testing"
	

	
	
	
	testutil "github.com/celestiaorg/celestia-app/v6/test/util"

	"github.com/cometbft/cometbft/abci/types"
	
	"github.com/stretchr/testify/require"
)


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
		_, err := testApp.CheckTx(&checkTxReq)
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

	
	for _, tx := range rawTxs {
		checkTxRequest := types.RequestCheckTx{
			Tx:   tx,
			Type: types.CheckTxType_New,
		}
		b.ReportAllocs()
		for b.Loop() {
			_, err := testApp.CheckTx(&checkTxRequest)
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
		_, err := testApp.FinalizeBlock(&finalizeBlockReq)
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
		_, err := testApp.FinalizeBlock(&finalizeBlockReq)
		require.NoError(b, err)
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
		_, err := testApp.PrepareProposal(&prepareProposalReq)
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
		_, err := testApp.PrepareProposal(&prepareProposalReq)
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
		_, err := testApp.ProcessProposal(&processProposalReq)
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
		_, err := testApp.ProcessProposal(&processProposalReq)
		require.NoError(b, err)
	}
}

