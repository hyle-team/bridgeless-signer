package gosdk

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"sync/atomic"

	"fmt"
	"net/http"
)

const rpcVersion = "2.0"

type RPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      uint32      `json:"id"`
}

type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error"`
	ID      int             `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Client struct {
	walletRPC string
	nodeRPC   string
	idCounter atomic.Uint32
}

func NewClient(walletRPC, nodeRPC string) *Client {
	return &Client{
		walletRPC: walletRPC,
		nodeRPC:   nodeRPC,
	}
}

func (c *Client) Call(method string, res interface{}, params interface{}, isWalletMethod bool) error {
	req, err := c.prepareMessage(method, params)
	if err != nil {
		return errors.Wrap(err, "failed to prepare request")
	}

	rpc := c.nodeRPC
	if isWalletMethod {
		rpc = c.walletRPC
	}
	resp, err := http.Post(rpc, "application/json", req)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	var rpcResponse RPCResponse
	if err = json.Unmarshal(responseBody, &rpcResponse); err != nil {
		return errors.Wrap(err, "failed to unmarshal response body")
	}

	if rpcResponse.Error != nil {
		return errors.New(fmt.Sprintf("RPC Error: Code=%d, Message=%s\n", rpcResponse.Error.Code, rpcResponse.Error.Message))
	}
	fmt.Printf("RPC Result: %s\n", string(rpcResponse.Result))

	if err = json.Unmarshal(rpcResponse.Result, res); err != nil {
		return errors.Wrap(err, "failed to unmarshal result")
	}
	return nil
}

func (c *Client) prepareMessage(method string, params interface{}) (*bytes.Buffer, error) {
	request := RPCRequest{
		JSONRPC: rpcVersion,
		Method:  method,
		Params:  params,
		ID:      c.nextID(),
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request")

	}
	return bytes.NewBuffer(requestBody), err
}

func (c *Client) nextID() uint32 {
	return c.idCounter.Add(1)
}
