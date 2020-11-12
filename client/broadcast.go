package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/terra-project/terra.go/msg"
	"github.com/terra-project/terra.go/tx"
)

// BroadcastReq broadcast request body
type BroadcastReq struct {
	Tx   tx.StdTxData `json:"tx"`
	Mode string       `json:"mode"`
}

// TxResponse response
type TxResponse struct {
	Height msg.Int `json:"height"`
	TxHash string  `json:"txhash"`
	Code   uint32  `json:"code,omitempty"`
	RawLog string  `json:"raw_log,omitempty"`
}

// Broadcast - no-lint
func (LCDClient LCDClient) Broadcast(stdTx tx.StdTx) (TxResponse, error) {
	broadcastReq := BroadcastReq{
		Tx:   stdTx.Value,
		Mode: "sync",
	}

	reqBytes, err := json.Marshal(broadcastReq)
	fmt.Println(string(reqBytes))

	if err != nil {
		return TxResponse{}, sdkerrors.Wrap(err, "failed to marshal")
	}

	resp, err := http.Post(LCDClient.URL+"/txs", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return TxResponse{}, sdkerrors.Wrap(err, "failed to broadcast")
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TxResponse{}, sdkerrors.Wrap(err, "failed to read response")
	}

	var txResponse TxResponse
	err = json.Unmarshal(out, &txResponse)
	if err != nil {
		return TxResponse{}, sdkerrors.Wrap(err, "failed to unmarshal response")
	}

	return txResponse, nil
}
