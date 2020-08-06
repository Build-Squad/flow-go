package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/onflow/flow/protobuf/go/flow/legacy/access"
	"github.com/onflow/flow/protobuf/go/flow/legacy/entities"

	"github.com/dapperlabs/flow-go/engine/access/rpc/handler"
	"github.com/dapperlabs/flow-go/engine/access/rpc/handler/legacy/convert"
	"github.com/dapperlabs/flow-go/model/flow"
)

type Handler struct {
	api   handler.AccessAPI
	chain flow.Chain
}

func New(api handler.AccessAPI, chain flow.Chain) *Handler {
	return &Handler{
		api:   api,
		chain: chain,
	}
}

// Ping the Access API server for a response.
func (h *Handler) Ping(ctx context.Context, req *access.PingRequest) (*access.PingResponse, error) {
	return &access.PingResponse{}, nil
}

func (h *Handler) GetNetworkParameters(
	context.Context,
	*access.GetNetworkParametersRequest,
) (*access.GetNetworkParametersResponse, error) {
	panic("implement me")
}

// SendTransaction submits a transaction to the network.
func (h *Handler) SendTransaction(
	ctx context.Context,
	req *access.SendTransactionRequest,
) (*access.SendTransactionResponse, error) {
	txMsg := req.GetTransaction()

	tx, err := convert.MessageToTransaction(txMsg, h.chain)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = h.api.SendTransaction(ctx, &tx)
	if err != nil {
		return nil, err
	}

	txID := tx.ID()

	return &access.SendTransactionResponse{
		Id: txID[:],
	}, nil
}

// GetLatestBlockHeader gets the latest sealed block header.
func (h *Handler) GetLatestBlockHeader(
	ctx context.Context,
	req *access.GetLatestBlockHeaderRequest,
) (*access.BlockHeaderResponse, error) {
	header, err := h.api.GetLatestBlockHeader(ctx, req.GetIsSealed())
	if err != nil {
		return nil, err
	}

	return blockHeaderResponse(header)
}

// GetBlockHeaderByHeight gets a block header by height.
func (h *Handler) GetBlockHeaderByHeight(
	ctx context.Context,
	req *access.GetBlockHeaderByHeightRequest,
) (*access.BlockHeaderResponse, error) {
	header, err := h.api.GetBlockHeaderByHeight(ctx, req.GetHeight())
	if err != nil {
		return nil, err
	}

	return blockHeaderResponse(header)
}

// GetBlockHeaderByID gets a block header by ID.
func (h *Handler) GetBlockHeaderByID(
	ctx context.Context,
	req *access.GetBlockHeaderByIDRequest,
) (*access.BlockHeaderResponse, error) {
	blockID := convert.MessageToIdentifier(req.GetId())

	header, err := h.api.GetBlockHeaderByID(ctx, blockID)
	if err != nil {
		return nil, err
	}

	return blockHeaderResponse(header)
}

// GetLatestBlock gets the latest sealed block.
func (h *Handler) GetLatestBlock(
	ctx context.Context,
	req *access.GetLatestBlockRequest,
) (*access.BlockResponse, error) {
	block, err := h.api.GetLatestBlock(ctx, req.GetIsSealed())
	if err != nil {
		return nil, err
	}

	return blockResponse(block)
}

// GetBlockByHeight gets a block by height.
func (h *Handler) GetBlockByHeight(
	ctx context.Context,
	req *access.GetBlockByHeightRequest,
) (*access.BlockResponse, error) {
	block, err := h.api.GetBlockByHeight(ctx, req.GetHeight())
	if err != nil {
		return nil, err
	}

	return blockResponse(block)
}

// GetBlockByHeight gets a block by ID.
func (h *Handler) GetBlockByID(
	ctx context.Context,
	req *access.GetBlockByIDRequest,
) (*access.BlockResponse, error) {
	blockID := convert.MessageToIdentifier(req.GetId())

	block, err := h.api.GetBlockByID(ctx, blockID)
	if err != nil {
		return nil, err
	}

	return blockResponse(block)
}

// GetCollectionByID gets a collection by ID.
func (h *Handler) GetCollectionByID(
	ctx context.Context,
	req *access.GetCollectionByIDRequest,
) (*access.CollectionResponse, error) {
	id := convert.MessageToIdentifier(req.GetId())

	col, err := h.api.GetCollectionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	colMsg, err := convert.LightCollectionToMessage(col)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &access.CollectionResponse{
		Collection: colMsg,
	}, nil
}

// GetTransaction gets a transaction by ID.
func (h *Handler) GetTransaction(
	ctx context.Context,
	req *access.GetTransactionRequest,
) (*access.TransactionResponse, error) {
	id := convert.MessageToIdentifier(req.GetId())

	tx, err := h.api.GetTransaction(ctx, id)
	if err != nil {
		return nil, err
	}

	return &access.TransactionResponse{
		Transaction: convert.TransactionToMessage(*tx),
	}, nil
}

// GetTransactionResult gets a transaction by ID.
func (h *Handler) GetTransactionResult(
	ctx context.Context,
	req *access.GetTransactionRequest,
) (*access.TransactionResultResponse, error) {
	id := convert.MessageToIdentifier(req.GetId())

	result, err := h.api.GetTransactionResult(ctx, id)
	if err != nil {
		return nil, err
	}

	return convert.TransactionResultToMessage(*result), nil
}

// GetAccount returns an account by address at the latest sealed block.
func (h *Handler) GetAccount(
	ctx context.Context,
	req *access.GetAccountRequest,
) (*access.GetAccountResponse, error) {
	address := flow.BytesToAddress(req.GetAddress())

	account, err := h.api.GetAccount(ctx, address)
	if err != nil {
		return nil, err
	}

	accountMsg, err := convert.AccountToMessage(account)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &access.GetAccountResponse{
		Account: accountMsg,
	}, nil
}

// GetAccountAtLatestBlock returns an account by address at the latest sealed block.
func (h *Handler) GetAccountAtLatestBlock(
	ctx context.Context,
	req *access.GetAccountAtLatestBlockRequest,
) (*access.AccountResponse, error) {
	address := flow.BytesToAddress(req.GetAddress())

	account, err := h.api.GetAccountAtLatestBlock(ctx, address)
	if err != nil {
		return nil, err
	}

	accountMsg, err := convert.AccountToMessage(account)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &access.AccountResponse{
		Account: accountMsg,
	}, nil
}

func (h *Handler) GetAccountAtBlockHeight(
	ctx context.Context,
	request *access.GetAccountAtBlockHeightRequest,
) (*access.AccountResponse, error) {
	panic("implement me")
}

// ExecuteScriptAtLatestBlock executes a script at a the latest block
func (h *Handler) ExecuteScriptAtLatestBlock(
	ctx context.Context,
	req *access.ExecuteScriptAtLatestBlockRequest,
) (*access.ExecuteScriptResponse, error) {
	script := req.GetScript()
	arguments := req.GetArguments()

	value, err := h.api.ExecuteScriptAtLatestBlock(ctx, script, arguments)
	if err != nil {
		return nil, err
	}

	return &access.ExecuteScriptResponse{
		Value: value,
	}, nil
}

// ExecuteScriptAtBlockHeight executes a script at a specific block height
func (h *Handler) ExecuteScriptAtBlockHeight(
	ctx context.Context,
	req *access.ExecuteScriptAtBlockHeightRequest,
) (*access.ExecuteScriptResponse, error) {
	script := req.GetScript()
	arguments := req.GetArguments()
	blockHeight := req.GetBlockHeight()

	value, err := h.api.ExecuteScriptAtBlockHeight(ctx, blockHeight, script, arguments)
	if err != nil {
		return nil, err
	}

	return &access.ExecuteScriptResponse{
		Value: value,
	}, nil
}

// ExecuteScriptAtBlockID executes a script at a specific block ID
func (h *Handler) ExecuteScriptAtBlockID(
	ctx context.Context,
	req *access.ExecuteScriptAtBlockIDRequest,
) (*access.ExecuteScriptResponse, error) {
	script := req.GetScript()
	arguments := req.GetArguments()
	blockID := convert.MessageToIdentifier(req.GetBlockId())

	value, err := h.api.ExecuteScriptAtBlockID(ctx, blockID, script, arguments)
	if err != nil {
		return nil, err
	}

	return &access.ExecuteScriptResponse{
		Value: value,
	}, nil
}

// GetEventsForHeightRange returns events matching a query.
func (h *Handler) GetEventsForHeightRange(
	ctx context.Context,
	req *access.GetEventsForHeightRangeRequest,
) (*access.EventsResponse, error) {
	eventType := req.GetType()
	startHeight := req.GetStartHeight()
	endHeight := req.GetEndHeight()

	results, err := h.api.GetEventsForHeightRange(ctx, eventType, startHeight, endHeight)
	if err != nil {
		return nil, err
	}

	return &access.EventsResponse{
		Results: blockEventsToMessages(results),
	}, nil
}

// GetEventsForBlockIDs returns events matching a set of block IDs.
func (h *Handler) GetEventsForBlockIDs(
	ctx context.Context,
	req *access.GetEventsForBlockIDsRequest,
) (*access.EventsResponse, error) {
	eventType := req.GetType()
	blockIDs := convert.MessagesToIdentifiers(req.GetBlockIds())

	results, err := h.api.GetEventsForBlockIDs(ctx, eventType, blockIDs)
	if err != nil {
		return nil, err
	}

	return &access.EventsResponse{
		Results: blockEventsToMessages(results),
	}, nil
}

// blockToHeaderResponse constructs a block header response from a block.
func blockToHeaderResponse(block *flow.Block) (*access.BlockHeaderResponse, error) {
	msg, err := convert.BlockHeaderToMessage(block.Header)
	if err != nil {
		return nil, err
	}

	return &access.BlockHeaderResponse{
		Block: msg,
	}, nil
}

func blockResponse(block *flow.Block) (*access.BlockResponse, error) {
	msg, err := convert.BlockToMessage(block)
	if err != nil {
		return nil, err
	}

	return &access.BlockResponse{
		Block: msg,
	}, nil
}

func blockHeaderResponse(header *flow.Header) (*access.BlockHeaderResponse, error) {
	msg, err := convert.BlockHeaderToMessage(header)
	if err != nil {
		return nil, err
	}

	return &access.BlockHeaderResponse{
		Block: msg,
	}, nil
}

func blockEventsToMessages(blocks []flow.BlockEvents) []*access.EventsResponse_Result {
	results := make([]*access.EventsResponse_Result, len(blocks))

	for i, block := range blocks {
		results[i] = blockEventsToMessage(block)
	}

	return results
}

func blockEventsToMessage(block flow.BlockEvents) *access.EventsResponse_Result {
	eventMessages := make([]*entities.Event, len(block.Events))
	for i, event := range block.Events {
		eventMessages[i] = convert.EventToMessage(event)
	}

	return &access.EventsResponse_Result{
		BlockId:     block.BlockID[:],
		BlockHeight: block.BlockHeight,
		Events:      eventMessages,
	}
}
