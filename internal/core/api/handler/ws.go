package handler

import (
	"fmt"
	"github.com/gorilla/websocket"
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
	// TODO: add check origin
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *ServiceHandler) CheckWithdrawalWs(w http.ResponseWriter, r *http.Request) {
	req, err := h.CheckWithdrawalRequest(&types.CheckWithdrawalRequest{OriginTxId: r.PathValue(paramOriginTxId)})
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
		h.logger.WithError(err).Debug("websocket upgrade error")
		return
	}

	gracefulClose := make(chan struct{})
	// TODO: add app ctx for graceful shutdown
	go h.watchWithdrawalStatus(ws, gracefulClose, depositIdentifier)
	go h.watchConnectionClosing(ws, gracefulClose)

}

func (h *ServiceHandler) watchConnectionClosing(ws *websocket.Conn, done chan struct{}) {
	defer close(done)

	for {
		// collecting only errors and close message to signalize writer
		mt, raw, err := ws.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}
		h.logger.Debug(fmt.Sprintf("received message: %x", raw))
	}
}

func (h *ServiceHandler) watchWithdrawalStatus(ws *websocket.Conn, connClosed chan struct{}, id data.DepositIdentifier) {
	defer ws.Close()

	// TODO: make ticker configurable
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var prevStatus types.WithdrawalStatus = -1
	db := h.db.New()

	// fast-start
	for ; true; <-ticker.C {
		select {
		case <-connClosed:
			// websocket is closed by client
			// close msg should already be sent by default close handler
			return
		default:
		}

		// do smth
		deposit, err := db.Get(id)
		if err != nil {
			h.logger.WithError(err).Error("failed to get deposit")
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
			h.logger.WithError(err).Error("failed to marshal deposit status")
			_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Internal server error"))
			return
		}
		if err = ws.WriteMessage(websocket.TextMessage, raw); err != nil {
			h.logger.WithError(err).Error("failed to write message to websocket")
			_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Internal server error"))
			return
		}

		if slices.Contains(
			[]types.WithdrawalStatus{
				types.WithdrawalStatus_TX_PENDING,
				types.WithdrawalStatus_TX_SUCCESSFUL,
				types.WithdrawalStatus_TX_FAILED,
				types.WithdrawalStatus_INVALID,
				types.WithdrawalStatus_FAILED,
			}, deposit.Status,
		) {
			err = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				h.logger.WithError(err).Error("failed to send close msg after finish")
			}
			return
		}

		prevStatus = deposit.Status
	}
}
