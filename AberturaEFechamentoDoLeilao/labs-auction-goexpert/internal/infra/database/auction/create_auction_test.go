package auction

import (
	"context"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestAuctionRepository_closeExpiredAuctions(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should close expired auctions", func(mt *mtest.T) {
		ctx := context.TODO()
		repo := &AuctionRepository{Collection: mt.Coll}

		expiredAuction := AuctionEntityMongo{
			Id:          "12345",
			ProductName: "Product Test",
			Category:    "Category Test",
			Description: "Description Test",
			Condition:   auction_entity.New,
			Status:      auction_entity.Open,
			Timestamp:   time.Now().Unix(),
			ExpiresAt:   time.Now().Add(-1 * time.Hour).Unix(),
		}

		activeAuction := AuctionEntityMongo{
			Id:          "67890",
			ProductName: "Product Active",
			Category:    "Category Active",
			Description: "Description Active",
			Condition:   auction_entity.New,
			Status:      auction_entity.Open,
			Timestamp:   time.Now().Unix(),
			ExpiresAt:   time.Now().Add(1 * time.Hour).Unix(),
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		_, err := repo.Collection.InsertOne(ctx, expiredAuction)
		assert.NoError(t, err)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		_, err = repo.Collection.InsertOne(ctx, activeAuction)
		assert.NoError(t, err)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		repo.closeExpiredAuctions(ctx)

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "auctions.test", mtest.FirstBatch, bson.D{
			{"_id", expiredAuction.Id},
			{"status", auction_entity.Closed},
		}))

		var result AuctionEntityMongo
		err = repo.Collection.FindOne(ctx, bson.M{"_id": expiredAuction.Id}).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, auction_entity.Closed, result.Status)

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "auctions.test", mtest.FirstBatch, bson.D{
			{"_id", activeAuction.Id},
			{"status", auction_entity.Open},
		}))

		err = repo.Collection.FindOne(ctx, bson.M{"_id": activeAuction.Id}).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, auction_entity.Open, result.Status)
	})
}
