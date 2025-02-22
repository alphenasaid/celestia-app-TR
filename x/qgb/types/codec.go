package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgDataCommitmentConfirm{}, "qgb/DataCommitmentConfirm", nil)
	cdc.RegisterConcrete(&MsgValsetConfirm{}, "qgb/ValSetConfirm", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDataCommitmentConfirm{},
	)

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgValsetConfirm{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
