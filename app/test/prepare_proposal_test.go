package app_test

import (
	"bytes"
	"testing"

	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/app/encoding"
	"github.com/celestiaorg/celestia-app/testutil"
	"github.com/celestiaorg/celestia-app/x/payment/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	core "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestPrepareProposal(t *testing.T) {
	signer := testutil.GenerateKeyringSigner(t, testAccName)

	encCfg := encoding.MakeEncodingConfig(app.ModuleEncodingRegisters...)

	testApp := testutil.SetupTestAppWithGenesisValSet(t)

	type test struct {
		input            abci.RequestPrepareProposal
		expectedMessages []*core.Message
		expectedTxs      int
	}

	firstNS := []byte{2, 2, 2, 2, 2, 2, 2, 2}
	firstMessage := bytes.Repeat([]byte{4}, 512)
	firstRawTx := generateRawTx(t, encCfg.TxConfig, firstNS, firstMessage, signer, 2, 4, 8, 16)

	secondNS := []byte{1, 1, 1, 1, 1, 1, 1, 1}
	secondMessage := []byte{2}
	secondRawTx := generateRawTx(t, encCfg.TxConfig, secondNS, secondMessage, signer, 2, 4, 8, 16)

	thirdNS := []byte{3, 3, 3, 3, 3, 3, 3, 3}
	thirdMessage := []byte{1}
	thirdRawTx := generateRawTx(t, encCfg.TxConfig, thirdNS, thirdMessage, signer, 2, 4, 8, 16)

	tests := []test{
		{
			input: abci.RequestPrepareProposal{
				BlockData: &core.Data{
					Txs: [][]byte{firstRawTx, secondRawTx, thirdRawTx},
				},
			},
			expectedMessages: []*core.Message{
				{
					NamespaceId: secondNS, // the second message should be first
					Data:        []byte{2},
				},
				{
					NamespaceId: firstNS,
					Data:        firstMessage,
				},
				{
					NamespaceId: thirdNS,
					Data:        []byte{1},
				},
			},
			expectedTxs: 3,
		},
	}

	for _, tt := range tests {
		res := testApp.PrepareProposal(tt.input)
		assert.Equal(t, tt.expectedMessages, res.BlockData.Messages.MessagesList)
		assert.Equal(t, tt.expectedTxs, len(res.BlockData.Txs))

		// verify the signatures of the prepared txs
		sdata, err := signer.GetSignerData()
		if err != nil {
			require.NoError(t, err)
		}
		dec := app.MalleatedTxDecoder(encCfg.TxConfig.TxDecoder())
		for _, tx := range res.BlockData.Txs {
			sTx, err := dec(tx)
			require.NoError(t, err)

			sigTx, ok := sTx.(authsigning.SigVerifiableTx)
			require.True(t, ok)

			sigs, err := sigTx.GetSignaturesV2()
			require.NoError(t, err)
			require.Equal(t, 1, len(sigs))
			sig := sigs[0]

			err = authsigning.VerifySignature(
				sdata.PubKey,
				sdata,
				sig.Data,
				encCfg.TxConfig.SignModeHandler(),
				sTx,
			)
			assert.NoError(t, err)
		}
	}
}

func generateRawTx(t *testing.T, txConfig client.TxConfig, ns, message []byte, signer *types.KeyringSigner, ks ...uint64) (rawTx []byte) {
	coin := sdk.Coin{
		Denom:  app.BondDenom,
		Amount: sdk.NewInt(10),
	}

	opts := []types.TxBuilderOption{
		types.SetFeeAmount(sdk.NewCoins(coin)),
		types.SetGasLimit(10000000),
	}

	// create a msg
	msg := generateSignedWirePayForData(t, ns, message, signer, opts, ks...)

	builder := signer.NewTxBuilder(opts...)

	tx, err := signer.BuildSignedTx(builder, msg)
	require.NoError(t, err)

	// encode the tx
	rawTx, err = txConfig.TxEncoder()(tx)
	require.NoError(t, err)

	return rawTx
}

func generateSignedWirePayForData(t *testing.T, ns, message []byte, signer *types.KeyringSigner, options []types.TxBuilderOption, ks ...uint64) *types.MsgWirePayForData {
	msg, err := types.NewWirePayForData(ns, message, ks...)
	if err != nil {
		t.Error(err)
	}

	err = msg.SignShareCommitments(signer, options...)
	if err != nil {
		t.Error(err)
	}

	return msg
}

const (
	testAccName = "test-account"
)
