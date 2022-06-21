package chain

import (
	"context"
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/spf13/cobra"
)

type Client struct {
	ckbUrl     string
	indexerUrl string
	client     rpc.Client
	ctx        context.Context
}

func NewClient(ctx context.Context, ckbUrl, indexerUrl string) *Client {
	rpcClient, err := rpc.DialWithIndexer(ckbUrl, indexerUrl)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("DialWithIndexer err: %v", err.Error()))
	}
	return &Client{
		ckbUrl:     ckbUrl,
		indexerUrl: indexerUrl,
		client:     rpcClient,
		ctx:        ctx,
	}
}

func (c *Client) Client() rpc.Client {
	return c.client
}

func (c *Client) GetTransactionByHash(hash types.Hash) *types.TransactionWithStatus {
	res, err := c.client.GetTransaction(c.ctx, hash)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("GetTransaction err: %v", err.Error()))
	}
	return res
}
