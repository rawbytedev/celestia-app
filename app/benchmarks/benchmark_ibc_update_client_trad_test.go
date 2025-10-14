//go:build benchmarks

package benchmarks_test

import (
	"fmt"
	"testing"

	testutil "github.com/celestiaorg/celestia-app/v6/test/util"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/require"
)

func BenchmarkTradIBC_CheckTx_Update_Client_Multi(b *testing.B) {
	testCases := []struct {
		numberOfValidators int
	}{
		{numberOfValidators: 2},
		{numberOfValidators: 10},
		{numberOfValidators: 25},
		{numberOfValidators: 50},
		{numberOfValidators: 75},
		{numberOfValidators: 100},
		{numberOfValidators: 125},
		{numberOfValidators: 150},
		{numberOfValidators: 175},
		{numberOfValidators: 200},
		{numberOfValidators: 225},
		{numberOfValidators: 250},
		{numberOfValidators: 300},
		{numberOfValidators: 400},
		{numberOfValidators: 500},
	}
	for _, testCase := range testCases {
		b.Run(fmt.Sprintf("number of validators: %d", testCase.numberOfValidators), func(b *testing.B) {
			benchmarkTradIBCCheckTxUpdateClient(b, testCase.numberOfValidators)
		})
	}
}

func benchmarkTradIBCCheckTxUpdateClient(b *testing.B, numberOfValidators int) {
	testApp, rawTxs := generateIBCUpdateClientTransaction(b, numberOfValidators, 1, 1)
	testApp.Commit()

	checkTxReq := types.RequestCheckTx{
		Type: types.CheckTxType_New,
		Tx:   rawTxs[0],
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

func BenchmarkTradIBC_FinalizeBlock_Update_Client_Multi(b *testing.B) {
	testCases := []struct {
		numberOfValidators int
	}{
		{numberOfValidators: 2},
		{numberOfValidators: 10},
		{numberOfValidators: 25},
		{numberOfValidators: 50},
		{numberOfValidators: 75},
		{numberOfValidators: 100},
		{numberOfValidators: 125},
		{numberOfValidators: 150},
		{numberOfValidators: 175},
		{numberOfValidators: 200},
		{numberOfValidators: 225},
		{numberOfValidators: 250},
		{numberOfValidators: 300},
		{numberOfValidators: 400},
		{numberOfValidators: 500},
	}
	for _, testCase := range testCases {
		b.Run(fmt.Sprintf("number of validators: %d", testCase.numberOfValidators), func(b *testing.B) {
			benchmarkTradIBCFinalizeBlockUpdateClient(b, testCase.numberOfValidators)
		})
	}
}

func benchmarkTradIBCFinalizeBlockUpdateClient(b *testing.B, numberOfValidators int) {
	testApp, rawTxs := generateIBCUpdateClientTransaction(b, numberOfValidators, 1, 1)

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

func BenchmarkTradIBC_PrepareProposal_Update_Client_Multi(b *testing.B) {
	testCases := []struct {
		numberOfTransactions, numberOfValidators int
	}{
		{numberOfTransactions: 6_000, numberOfValidators: 2},
		{numberOfTransactions: 3_000, numberOfValidators: 10},
		{numberOfTransactions: 2_000, numberOfValidators: 25},
		{numberOfTransactions: 1_000, numberOfValidators: 50},
		{numberOfTransactions: 500, numberOfValidators: 75},
		{numberOfTransactions: 500, numberOfValidators: 100},
		{numberOfTransactions: 500, numberOfValidators: 125},
		{numberOfTransactions: 500, numberOfValidators: 150},
		{numberOfTransactions: 500, numberOfValidators: 175},
		{numberOfTransactions: 500, numberOfValidators: 200},
		{numberOfTransactions: 500, numberOfValidators: 225},
		{numberOfTransactions: 500, numberOfValidators: 250},
		{numberOfTransactions: 500, numberOfValidators: 300},
		{numberOfTransactions: 500, numberOfValidators: 400},
		{numberOfTransactions: 500, numberOfValidators: 500},
	}
	for _, testCase := range testCases {
		b.Run(fmt.Sprintf("number of validators: %d", testCase.numberOfValidators), func(b *testing.B) {
			benchmarkTradIBCPrepareProposalUpdateClient(b, testCase.numberOfValidators, testCase.numberOfTransactions)
		})
	}
}

func benchmarkTradIBCPrepareProposalUpdateClient(b *testing.B, numberOfValidators, count int) {
	testApp, rawTxs := generateIBCUpdateClientTransaction(b, numberOfValidators, count, count)

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

func BenchmarkTradIBC_ProcessProposal_Update_Client_Multi(b *testing.B) {
	testCases := []struct {
		numberOfTransactions, numberOfValidators int
	}{
		{numberOfTransactions: 6_000, numberOfValidators: 2},
		{numberOfTransactions: 3_000, numberOfValidators: 10},
		{numberOfTransactions: 2_000, numberOfValidators: 25},
		{numberOfTransactions: 1_000, numberOfValidators: 50},
		{numberOfTransactions: 500, numberOfValidators: 75},
		{numberOfTransactions: 500, numberOfValidators: 100},
		{numberOfTransactions: 500, numberOfValidators: 125},
		{numberOfTransactions: 500, numberOfValidators: 150},
		{numberOfTransactions: 500, numberOfValidators: 175},
		{numberOfTransactions: 500, numberOfValidators: 200},
		{numberOfTransactions: 500, numberOfValidators: 225},
		{numberOfTransactions: 500, numberOfValidators: 250},
		{numberOfTransactions: 500, numberOfValidators: 300},
		{numberOfTransactions: 500, numberOfValidators: 400},
		{numberOfTransactions: 500, numberOfValidators: 500},
	}
	for _, testCase := range testCases {
		b.Run(fmt.Sprintf("number of validators: %d", testCase.numberOfValidators), func(b *testing.B) {
			benchmarkTradIBCProcessProposalUpdateClient(b, testCase.numberOfValidators, testCase.numberOfTransactions)
		})
	}
}

func benchmarkTradIBCProcessProposalUpdateClient(b *testing.B, numberOfValidators, count int) {
	testApp, rawTxs := generateIBCUpdateClientTransaction(b, numberOfValidators, count, count)

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
		SquareSize:   prepareProposalResp.SquareSize,
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
