package auction

import (
	"context"
	"os"
	"strconv"
	"time"

	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
	ExpiresAt   int64                           `bson:"expires_at"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(ctx context.Context, auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionMongo := AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
		ExpiresAt:   auctionEntity.Timestamp.Add(24 * time.Hour).Unix(),
	}

	_, err := ar.Collection.InsertOne(ctx, auctionMongo)
	if err != nil {
		logger.Error("Failed to create auction", err)
	}
	return nil
}

func (ar *AuctionRepository) MonitorExpiredAuctions(ctx context.Context) {
	interval := getMonitorInterval()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ar.closeExpiredAuctions(ctx)
		case <-ctx.Done():
			logger.Info("Stopping auction expiration monitor.")
			return
		}
	}
}

func (ar *AuctionRepository) closeExpiredAuctions(ctx context.Context) {
	now := time.Now().Unix()

	filter := bson.M{
		"expires_at": bson.M{"$lt": now},
		"status":     auction_entity.Open,
	}

	update := bson.M{
		"$set": bson.M{"status": auction_entity.Closed},
	}
	result, err := ar.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Error updating expired auctions", err)
		return
	}

	logger.Info("Closed expired auctions", zap.Int64("modified_count", result.ModifiedCount))
}

func getMonitorInterval() time.Duration {
	intervalStr := os.Getenv("AUCTION_MONITOR_INTERVAL")
	if intervalStr == "" {
		return 1 * time.Minute
	}
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		logger.Error("Invalid AUCTION_MONITOR_INTERVAL, using default 1m", err)
		return 1 * time.Minute
	}
	return time.Duration(interval) * time.Second
}
