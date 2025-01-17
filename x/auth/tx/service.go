package tx

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"

	gogogrpc "github.com/gogo/protobuf/grpc"
	"github.com/golang/protobuf/proto" //nolint: staticcheck
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/line/lbm-sdk/client"
	codectypes "github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	pagination "github.com/line/lbm-sdk/types/query"
	txtypes "github.com/line/lbm-sdk/types/tx"
)

// baseAppSimulateFn is the signature of the Baseapp#Simulate function.
type baseAppSimulateFn func(txBytes []byte) (sdk.GasInfo, *sdk.Result, error)

// txServer is the server for the protobuf Tx service.
type txServer struct {
	clientCtx         client.Context
	simulate          baseAppSimulateFn
	interfaceRegistry codectypes.InterfaceRegistry
}

// NewTxServer creates a new Tx service server.
func NewTxServer(clientCtx client.Context, simulate baseAppSimulateFn, interfaceRegistry codectypes.InterfaceRegistry) txtypes.ServiceServer {
	return txServer{
		clientCtx:         clientCtx,
		simulate:          simulate,
		interfaceRegistry: interfaceRegistry,
	}
}

var _ txtypes.ServiceServer = txServer{}

const (
	eventFormat = "{eventType}.{eventAttribute}={value}"
)

// TxsByEvents implements the ServiceServer.TxsByEvents RPC method.
func (s txServer) GetTxsEvent(ctx context.Context, req *txtypes.GetTxsEventRequest) (*txtypes.GetTxsEventResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	page, limit, err := pagination.ParsePagination(req.Pagination)
	if err != nil {
		return nil, err
	}
	orderBy := parseOrderBy(req.OrderBy)

	if len(req.Events) == 0 {
		return nil, status.Error(codes.InvalidArgument, "must declare at least one event to search")
	}

	for _, event := range req.Events {
		if !strings.Contains(event, "=") || strings.Count(event, "=") > 1 {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid event; event %s should be of the format: %s", event, eventFormat))
		}
	}

	result, err := QueryTxsByEvents(s.clientCtx, req.Events, page, limit, orderBy)
	if err != nil {
		return nil, err
	}

	// Create a proto codec, we need it to unmarshal the tx bytes.
	txsList := make([]*txtypes.Tx, len(result.Txs))

	for i, tx := range result.Txs {
		protoTx, ok := tx.Tx.GetCachedValue().(*txtypes.Tx)
		if !ok {
			return nil, status.Errorf(codes.Internal, "expected %T, got %T", txtypes.Tx{}, tx.Tx.GetCachedValue())
		}

		txsList[i] = protoTx
	}

	return &txtypes.GetTxsEventResponse{
		Txs:         txsList,
		TxResponses: result.Txs,
		Pagination: &pagination.PageResponse{
			Total: result.TotalCount,
		},
	}, nil
}

// Simulate implements the ServiceServer.Simulate RPC method.
func (s txServer) Simulate(ctx context.Context, req *txtypes.SimulateRequest) (*txtypes.SimulateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid empty tx")
	}

	txBytes := req.TxBytes
	if txBytes == nil && req.Tx != nil {
		// This block is for backwards-compatibility.
		// We used to support passing a `Tx` in req. But if we do that, sig
		// verification might not pass, because the .Marshal() below might not
		// be the same marshaling done by the client.
		var err error
		txBytes, err = proto.Marshal(req.Tx)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid tx; %v", err)
		}
	}

	if txBytes == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty txBytes is not allowed")
	}

	gasInfo, result, err := s.simulate(txBytes)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%v With gas wanted: '%d' and gas used: '%d' ", err, gasInfo.GasWanted, gasInfo.GasUsed)
	}

	return &txtypes.SimulateResponse{
		GasInfo: &gasInfo,
		Result:  result,
	}, nil
}

// GetTx implements the ServiceServer.GetTx RPC method.
func (s txServer) GetTx(ctx context.Context, req *txtypes.GetTxRequest) (*txtypes.GetTxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	if n := len(req.Hash); n != sha256.Size*2 {
		if n == 0 {
			return nil, status.Error(codes.InvalidArgument, "tx hash cannot be empty")
		}
		return nil, status.Error(codes.InvalidArgument, "The length of tx hash must be 64")
	}

	// TODO We should also check the proof flag in gRPC header.
	// https://github.com/cosmos/cosmos-sdk/issues/7036.
	result, err := QueryTx(s.clientCtx, req.Hash)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, status.Errorf(codes.NotFound, "tx not found: %s", req.Hash)
		}

		return nil, err
	}

	protoTx, ok := result.Tx.GetCachedValue().(*txtypes.Tx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "expected %T, got %T", txtypes.Tx{}, result.Tx.GetCachedValue())
	}

	return &txtypes.GetTxResponse{
		Tx:         protoTx,
		TxResponse: result,
	}, nil
}

// GetBlockWithTxs returns a block with decoded txs.
func (s txServer) GetBlockWithTxs(ctx context.Context, req *txtypes.GetBlockWithTxsRequest) (*txtypes.GetBlockWithTxsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "service not supported")
}

func (s txServer) BroadcastTx(ctx context.Context, req *txtypes.BroadcastTxRequest) (*txtypes.BroadcastTxResponse, error) {
	return client.TxServiceBroadcast(ctx, s.clientCtx, req)
}

// RegisterTxService registers the tx service on the gRPC router.
func RegisterTxService(
	qrt gogogrpc.Server,
	clientCtx client.Context,
	simulateFn baseAppSimulateFn,
	interfaceRegistry codectypes.InterfaceRegistry,
) {
	txtypes.RegisterServiceServer(
		qrt,
		NewTxServer(clientCtx, simulateFn, interfaceRegistry),
	)
}

// RegisterGRPCGatewayRoutes mounts the tx service's GRPC-gateway routes on the
// given Mux.
func RegisterGRPCGatewayRoutes(clientConn gogogrpc.ClientConn, mux *runtime.ServeMux) {
	if err := txtypes.RegisterServiceHandlerClient(context.Background(), mux, txtypes.NewServiceClient(clientConn)); err != nil {
		panic(err)
	}
}

func parseOrderBy(orderBy txtypes.OrderBy) string {
	switch orderBy {
	case txtypes.OrderBy_ORDER_BY_ASC:
		return "asc"
	case txtypes.OrderBy_ORDER_BY_DESC:
		return "desc"
	default:
		return "" // Defaults to Tendermint's default, which is `asc` now.
	}
}
