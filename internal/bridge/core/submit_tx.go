package core

import bridgetypes "github.com/hyle-team/bridgeless-core/x/bridge/types"

func (c *Connector) SubmitDeposits(depositTxs ...bridgetypes.Transaction) error {
	if len(depositTxs) == 0 {
		return nil
	}

	msg := bridgetypes.NewMsgSubmitTransactions(c.settings.Account.CosmosAddress(), depositTxs...)

	return c.submitMsgs(msg)
}
