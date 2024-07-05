package processor

import (
	"fmt"
	"testing"

	"github.com/hyle-team/bridgeless-signer/internal/bridge/evm"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/config"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/internal/data/pg"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/kit/kv"
)

const (
	depositTxHash    = "0x28246db1eccf7d7d431d6b3f05a5e1a34d3155ebbd421b3b0efbeeb339c8af77"
	depositTxEventId = 1
	amoyChainId      = "80002"
)

func TestProcessor_HappyPath(t *testing.T) {
	cfg := config.New(kv.MustFromEnv())

	proxies, err := evm.NewProxiesRepository(cfg.Chains(), cfg.Signer().Address())
	if err != nil {
		t.Fatal(errors.Wrap(err, "failed to create proxies repository"))
	}

	db := pg.NewDepositsQ(cfg.DB())

	processor := New(proxies, db, cfg.Signer())

	deposit := data.Deposit{
		DepositIdentifier: data.DepositIdentifier{
			TxHash:    depositTxHash,
			TxEventId: depositTxEventId,
			ChainId:   amoyChainId,
		},
		Status: types.WithdrawStatus_PROCESSING,
	}

	deposit.Id, err = db.Insert(deposit)
	if err != nil {
		if !errors.Is(err, data.ErrAlreadySubmitted) {
			t.Fatal(errors.Wrap(err, "failed to insert deposit"))
		}
	}

	req := bridgeTypes.GetDepositRequest{
		DepositDbId:       deposit.Id,
		DepositIdentifier: deposit.DepositIdentifier,
	}

	formWithdrawRequest, _, err := processor.ProcessGetDepositRequest(req)
	if err != nil {
		t.Fatal(errors.Wrap(err, "failed to process get deposit request"))
	}

	assert.Equal(t, "100000000000000000000", formWithdrawRequest.Data.Amount.String())

	// modify the amount to send a different amount
	formWithdrawRequest.Data.Amount.SetString("123456", 10)

	formedWithRequest, _, err := processor.ProcessFormWithdrawRequest(*formWithdrawRequest)
	if err != nil {
		t.Fatal(errors.Wrap(err, "failed to process form withdraw request"))
	}

	signedRequest, _, err := processor.ProcessSignWithdrawRequest(*formedWithRequest)
	if err != nil {
		t.Fatal(errors.Wrap(err, "failed to process sign withdraw request"))
	}

	_, err = processor.ProcessSendWithdrawRequest(*signedRequest)
	if err != nil {
		t.Fatal(errors.Wrap(err, "failed to process send withdraw request"))
	}

	fmt.Println(signedRequest.Transaction.Hash().String())
}
