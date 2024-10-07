package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/hyle-team/bridgeless-signer/internal/core/api/ctx"
	"github.com/hyle-team/bridgeless-signer/internal/core/api/requests"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"slices"
	"time"
)

const paramOriginTxId = "origin_tx_id"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func CheckWithdrawalWs(w http.ResponseWriter, r *http.Request) {
	var (
		ctxt    = r.Context()
		proxies = ctx.Proxies(ctxt)
	)

	req, err := requests.CheckWithdrawalRequest(
		&types.CheckWithdrawalRequest{
			OriginTxId: chi.URLParam(r, paramOriginTxId),
		}, proxies,
	)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	depositIdentifier := data.DepositIdentifier{
		TxHash:    req.Deposit.TxHash,
		TxEventId: int(req.Deposit.TxEventIndex),
		ChainId:   req.Deposit.ChainId,
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		ctx.Logger(ctxt).WithError(err).Debug("websocket upgrade error")
		return
	}

	gracefulClose := make(chan struct{})
	go watchConnectionClosing(ws, gracefulClose)
	watchWithdrawalStatus(ctxt, ws, gracefulClose, depositIdentifier)
}

func watchConnectionClosing(ws *websocket.Conn, done chan struct{}) {
	defer close(done)

	for {
		// collecting only errors and close message to signalize writer
		mt, _, err := ws.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}
	}
}

func watchWithdrawalStatus(ctxt context.Context, ws *websocket.Conn, connClosed chan struct{}, id data.DepositIdentifier) {
	defer ws.Close()

	var (
		db     = ctx.DB(ctxt)
		logger = ctx.Logger(ctxt)

		prevStatus types.WithdrawalStatus = -1

		cancelled, graceful bool
		// TODO: make ticker configurable
		ticker     = time.NewTicker(1 * time.Second)
		tillCancel = func() {
			select {
			case <-connClosed:
				cancelled = true
			case <-ctxt.Done():
				cancelled, graceful = true, true
			case <-ticker.C:
				// doing nothing, just waiting some period
			}
		}
	)

	defer ticker.Stop()

	// fast-start without waiting for initial tick
	for ; ; tillCancel() {
		if cancelled {
			if graceful {
				_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, "Server shutting down"))
			}
			return
		}

		deposit, err := db.Get(id)
		if err != nil {
			logger.WithError(err).Error("failed to get deposit")
			_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Internal server error"))
			return
		}
		if deposit == nil {
			_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(4004, "deposit not found"))
			return
		}
		if deposit.Status == prevStatus {
			continue
		}

		raw, err := protojson.Marshal(deposit.ToStatusResponse())
		if err != nil {
			logger.WithError(err).Error("failed to marshal deposit status")
			_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Internal server error"))
			return
		}
		if err = ws.WriteMessage(websocket.TextMessage, raw); err != nil {
			logger.WithError(err).Error("failed to write message to websocket")
			_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Internal server error"))
			return
		}

		// is it a time for websocket closing
		if slices.Contains(
			[]types.WithdrawalStatus{
				// transaction is sent
				types.WithdrawalStatus_TX_PENDING,
				types.WithdrawalStatus_TX_SUCCESSFUL,
				types.WithdrawalStatus_TX_FAILED,
				// ready to be sent
				types.WithdrawalStatus_WITHDRAWAL_SIGNED,
				// data invalid or something goes wrong
				types.WithdrawalStatus_INVALID,
				types.WithdrawalStatus_FAILED,
			}, deposit.Status,
		) {
			err = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.WithError(err).Error("failed to send close msg after finish")
			}
			return
		}

		prevStatus = deposit.Status
	}
}
