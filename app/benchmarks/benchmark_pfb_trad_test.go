//go:build benchmarks

package benchmarks_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/log"
	testutil "github.com/celestiaorg/celestia-app/v6/test/util"
	blobtx "github.com/celestiaorg/go-square/v3/tx"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/require"
)

func init() {
	testutil.TestAppLogger = log.NewNopLogger()
}

func TradBenchmarkCheckTx_PFB_Multi(b *testing.B) {
	testCases := []struct {
		blobSize int
	}{
		{blobSize: 300},
		{blobSize: 500},
		{blobSize: 1000},
		{blobSize: 5000},
		{blobSize: 10_000},
		{blobSize: 50_000},
		{blobSize: 100_000},
		{blobSize: 200_000},
		{blobSize: 300_000},
		{blobSize: 400_000},
		{blobSize: 500_000},
		{blobSize: 1_000_000},
		{blobSize: 2_000_000}, // maxTxSize is capped at 2MiB in checkTx
	}
	for _, testCase := range testCases {
		b.Run(fmt.Sprintf("%d bytes", testCase.blobSize), func(b *testing.B) {
			TradbenchmarkCheckTxPFB(b, testCase.blobSize)
		})
	}
}

func TradbenchmarkCheckTxPFB(b *testing.B, size int) {
	testApp, rawTxs := generatePayForBlobTransactions(b, 1, size)

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

	checkTxRequest := types.RequestCheckTx{
		Tx:   rawTxs[0],
		Type: types.CheckTxType_New,
	}
	b.ReportAllocs()
	for b.Loop() {
		resp, err := testApp.CheckTx(&checkTxRequest)
		require.NoError(b, err)
	}
}

func TradBenchmarkFinalizeBlock_PFB_Multi(b *testing.B) {
	testCases := []struct {
		blobSize int
	}{
		{blobSize: 300},
		{blobSize: 500},
		{blobSize: 1000},
		{blobSize: 5000},
		{blobSize: 10_000},
		{blobSize: 50_000},
		{blobSize: 100_000},
		{blobSize: 200_000},
		{blobSize: 300_000},
		{blobSize: 400_000},
		{blobSize: 500_000},
		{blobSize: 1_000_000},
		{blobSize: 2_000_000},
		{blobSize: 3_000_000},
		{blobSize: 4_000_000},
		{blobSize: 5_000_000},
		{blobSize: 6_000_000},
	}
	for _, testCase := range testCases {
		b.Run(fmt.Sprintf("%d bytes", testCase.blobSize), func(b *testing.B) {
			TradbenchmarkFinalizeBlockPFB(b, testCase.blobSize)
		})
	}
}

func TradbenchmarkFinalizeBlockPFB(b *testing.B, size int) {
	testApp, rawTxs := generatePayForBlobTransactions(b, 1, size)

	blobTx, ok, err := blobtx.UnmarshalBlobTx(rawTxs[0])
	require.NoError(b, err)
	require.True(b, ok)

	finalizeBlockReq := types.RequestFinalizeBlock{
		Time:   testutil.GenesisTime.Add(blockTime),
		Height: testApp.LastBlockHeight() + 1,
		Hash:   testApp.LastCommitID().Hash,
		Txs:    [][]byte{blobTx.Tx},
	}
	b.ReportAllocs()
	for b.Loop() {
		resp, err := testApp.FinalizeBlock(&finalizeBlockReq)
		require.NoError(b, err)
	}
}

func TradBenchmarkPrepareProposal_PFB_Multi(b *testing.B) {
	testCases := []struct {
		numberOfTransactions, blobSize int
	}{
		{numberOfTransactions: 15_000, blobSize: 300},
		{numberOfTransactions: 10_000, blobSize: 500},
		{numberOfTransactions: 6_000, blobSize: 1000},
		{numberOfTransactions: 3_000, blobSize: 5000},
		{numberOfTransactions: 1_000, blobSize: 10_000},
		{numberOfTransactions: 500, blobSize: 50_000},
		{numberOfTransactions: 100, blobSize: 100_000},
		{numberOfTransactions: 100, blobSize: 200_000},
		{numberOfTransactions: 50, blobSize: 300_000},
		{numberOfTransactions: 50, blobSize: 400_000},
		{numberOfTransactions: 30, blobSize: 500_000},
		{numberOfTransactions: 10, blobSize: 1_000_000},
		{numberOfTransactions: 5, blobSize: 2_000_000},
		{numberOfTransactions: 3, blobSize: 3_000_000},
		{numberOfTransactions: 3, blobSize: 4_000_000},
		{numberOfTransactions: 2, blobSize: 5_000_000},
		{numberOfTransactions: 2, blobSize: 6_000_000},
	}
	for _, testCase := range testCases {
		b.Run(fmt.Sprintf("%d transactions of %d bytes", testCase.numberOfTransactions, testCase.blobSize), func(b *testing.B) {
			TradbenchmarkPrepareProposalPFB(b, testCase.numberOfTransactions, testCase.blobSize)
		})
	}
}

func TradbenchmarkPrepareProposalPFB(b *testing.B, count, size int) {
	testApp, rawTxs := generatePayForBlobTransactions(b, count, size)

	prepareProposalReq := types.RequestPrepareProposal{
		Txs:    rawTxs,
		Height: testApp.LastBlockHeight() + 1,
	}
	b.ReportAllocs()
	for b.Loop() {
		prepareProposalResp, err := testApp.PrepareProposal(&prepareProposalReq)
		require.NoError(b, err)
	}
}

func TradBenchmarkProcessProposal_PFB_Multi(b *testing.B) {
	testCases := []struct {
		numberOfTransactions, blobSize int
	}{
		{numberOfTransactions: 15_000, blobSize: 300},
		{numberOfTransactions: 10_000, blobSize: 500},
		{numberOfTransactions: 6_000, blobSize: 1000},
		{numberOfTransactions: 3_000, blobSize: 5000},
		{numberOfTransactions: 1_000, blobSize: 10_000},
		{numberOfTransactions: 500, blobSize: 50_000},
		{numberOfTransactions: 100, blobSize: 100_000},
		{numberOfTransactions: 100, blobSize: 200_000},
		{numberOfTransactions: 50, blobSize: 300_000},
		{numberOfTransactions: 50, blobSize: 400_000},
		{numberOfTransactions: 30, blobSize: 500_000},
		{numberOfTransactions: 10, blobSize: 1_000_000},
		{numberOfTransactions: 5, blobSize: 2_000_000},
		{numberOfTransactions: 3, blobSize: 3_000_000},
		{numberOfTransactions: 3, blobSize: 4_000_000},
		{numberOfTransactions: 2, blobSize: 5_000_000},
		{numberOfTransactions: 2, blobSize: 6_000_000},
	}
	for _, testCase := range testCases {
		b.Run(fmt.Sprintf("%d transactions of %d bytes", testCase.numberOfTransactions, testCase.blobSize), func(b *testing.B) {
			TradbenchmarkProcessProposalPFB(b, testCase.numberOfTransactions, testCase.blobSize)
		})
	}
}

func TradbenchmarkProcessProposalPFB(b *testing.B, count, size int) {
	testApp, rawTxs := generatePayForBlobTransactions(b, count, size)

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
	for b.Loop() {
		resp, err := testApp.ProcessProposal(&processProposalReq)
		require.NoError(b, err)
	}
}
