package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateTrade{}

func NewMsgCreateTrade(creator string, coinInputs []CoinInput, itemInputs []ItemInput, coinOutputs sdk.Coins, itemOutputs []ItemRef, extraInfo string) *MsgCreateTrade {
	return &MsgCreateTrade{
		Creator:     creator,
		CoinInputs:  coinInputs,
		ItemInputs:  itemInputs,
		CoinOutputs: coinOutputs,
		ItemOutputs: itemOutputs,
		ExtraInfo:   extraInfo,
	}
}

func (msg *MsgCreateTrade) Route() string {
	return RouterKey
}

func (msg *MsgCreateTrade) Type() string {
	return "CreateTrade"
}

func (msg *MsgCreateTrade) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateTrade) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateTrade) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	for i, coinInput := range msg.CoinInputs {
		if !coinInput.Coins.IsValid() {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid coinInputs at index %d", i)
		}
	}

	if !msg.CoinOutputs.Empty() && !msg.CoinOutputs.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid coinOutputs")
	}

	// ensure that there is only one payment token
	for _, coinInput := range msg.CoinInputs {

		paymentDenom := ""
		for _, coin := range coinInput.Coins {
			switch {
			case !IsCookbookDenom(coin.Denom) && !IsIBCDenomRepresentation(coin.Denom):
				if paymentDenom != "" {
					return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "multiple paymentDenoms in CoinInputs")
				}
				paymentDenom = coin.Denom
			case IsIBCDenomRepresentation(coin.Denom):
				if paymentDenom != "" {
					return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "multiple paymentDenoms in CoinInputs")
				}
				paymentDenom = coin.Denom
			}
		}

		for _, coin := range msg.CoinOutputs {
			switch {
			case !IsCookbookDenom(coin.Denom) && !IsIBCDenomRepresentation(coin.Denom):
				if coin.Denom != paymentDenom && paymentDenom != "" {
					return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "multiple paymentDenoms in CoinOutputs")
				} else if paymentDenom == "" {
					paymentDenom = coin.Denom
				}
			case IsIBCDenomRepresentation(coin.Denom):
				if coin.Denom != paymentDenom && paymentDenom != "" {
					return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "multiple paymentDenoms in CoinOutputs")
				} else if paymentDenom == "" {
					paymentDenom = coin.Denom
				}
			}
		}
	}

	for _, item := range msg.ItemOutputs {
		err := ValidateItemID(item.ItemId)
		if err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
		err = ValidateID(item.CookbookId)
		if err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
	}

	for _, ii := range msg.ItemInputs {
		if err = ValidateItemInput(ii); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
	}

	return nil
}

var _ sdk.Msg = &MsgCancelTrade{}

func NewMsgCancelTrade(creator string, id uint64) *MsgCancelTrade {
	return &MsgCancelTrade{
		Id:      id,
		Creator: creator,
	}
}

func (msg *MsgCancelTrade) Route() string {
	return RouterKey
}

func (msg *MsgCancelTrade) Type() string {
	return "CancelTrade"
}

func (msg *MsgCancelTrade) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCancelTrade) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCancelTrade) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
