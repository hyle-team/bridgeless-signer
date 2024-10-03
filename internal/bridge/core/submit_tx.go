package core

import (
	bridgetypes "github.com/hyle-team/bridgeless-core/x/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

func (c *Connector) SubmitDeposits(depositTxs ...bridgetypes.Transaction) error {
	if len(depositTxs) == 0 {
		return nil
	}

	msg := bridgetypes.NewMsgSubmitTransactions(c.settings.Account.CosmosAddress(), depositTxs...)
	if err := c.submitMsgs(msg); err != nil {
		if errors.Is(err, bridgetypes.ErrTranscationAlreadySubmitted.GRPCStatus().Err()) {
			return types.ErrTransactionAlreadySubmitted
		}
	}

	return nil
}
