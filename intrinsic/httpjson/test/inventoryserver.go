// Copyright 2023 Intrinsic Innovation LLC

// Package inventoryserver implements InventoryService.
package inventoryserver

import (
	"context"

	log "github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ispb "intrinsic/httpjson/test/inventory_service_go_proto"
	pb "intrinsic/httpjson/test/inventory_service_go_proto"
)

// InventoryServer implements the pb.InventoryServiceServer interface.
type InventoryServer struct {
	ispb.UnimplementedInventoryServiceServer

	skus  map[string]*pb.Sku // sku_id -> Sku
	stock map[string]int64   // sku_id -> quantity
}

// NewInventoryServer creates a new server instance with initialized maps.
func NewInventoryServer() *InventoryServer {
	return &InventoryServer{
		skus:  make(map[string]*pb.Sku),
		stock: make(map[string]int64),
	}
}

// AddSku registers a new SKU in the inventory system.
func (s *InventoryServer) AddSku(ctx context.Context, req *pb.AddSkuRequest) (*pb.Sku, error) {
	log.InfoContextf(ctx, "Received AddSku request for ID: %s", req.GetSkuId())

	if req.GetSkuId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "sku_id cannot be empty")
	}

	// Check if SKU already exists
	if _, ok := s.skus[req.GetSkuId()]; ok {
		return nil, status.Errorf(codes.AlreadyExists, "SKU with ID %s already exists", req.GetSkuId())
	}

	// Create and store the new SKU
	sku := &pb.Sku{
		SkuId:       req.GetSkuId(),
		DisplayName: req.GetDisplayName(),
	}
	s.skus[sku.GetSkuId()] = sku

	// Initialize its stock level to 0
	s.stock[sku.GetSkuId()] = 0

	return sku, nil
}

// SetStockLevel sets how many items of SKU are in stock.
func (s *InventoryServer) SetStockLevel(ctx context.Context, req *pb.SetStockLevelRequest) (*pb.StockLevel, error) {
	log.InfoContextf(ctx, "Received SetStockLevel request for ID %s: quantity %d", req.GetSkuId(), req.GetQuantity())

	if req.GetSkuId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "sku_id cannot be empty")
	}
	if req.GetQuantity() < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "quantity cannot be negative")
	}

	// Check if the SKU exists before setting stock
	if _, ok := s.skus[req.GetSkuId()]; !ok {
		return nil, status.Errorf(codes.NotFound, "SKU with ID %s not found", req.GetSkuId())
	}

	// Set the stock level
	s.stock[req.GetSkuId()] = req.GetQuantity()

	return &pb.StockLevel{
		SkuId:    req.GetSkuId(),
		Quantity: req.GetQuantity(),
	}, nil
}

// GetStockLevel get the number of SKU in stock.
func (s *InventoryServer) GetStockLevel(ctx context.Context, req *pb.GetStockLevelRequest) (*pb.StockLevel, error) {
	log.InfoContextf(ctx, "Received GetStockLevel request for ID: %s", req.GetSkuId())

	if req.GetSkuId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "sku_id cannot be empty")
	}

	// Check if the SKU exists
	if _, ok := s.skus[req.GetSkuId()]; !ok {
		return nil, status.Errorf(codes.NotFound, "SKU with ID %s not found", req.GetSkuId())
	}

	// Get stock
	quantity := s.stock[req.GetSkuId()]

	return &pb.StockLevel{
		SkuId:    req.GetSkuId(),
		Quantity: quantity,
	}, nil
}

// ListSkus lists all SKUs and stock levels.
func (s *InventoryServer) ListSkus(ctx context.Context, req *pb.ListSkusRequest) (*pb.ListSkusResponse, error) {
	log.InfoContextf(ctx, "Received ListSkus request")
	skusWithStock := make([]*pb.SkuWithStockLevel, 0, len(s.skus))

	for skuID, sku := range s.skus {
		quantity := s.stock[skuID]
		skusWithStock = append(skusWithStock, &pb.SkuWithStockLevel{
			Sku: sku,
			StockLevel: &pb.StockLevel{
				SkuId:    skuID,
				Quantity: quantity,
			},
		})
	}

	return &pb.ListSkusResponse{
		Skus: skusWithStock,
	}, nil
}
