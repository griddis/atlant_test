package repository

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/griddis/atlant_test/tools/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Test_mongodbStore_UpdatePrice(t *testing.T) {
	logger := logging.NewLogger("debug", "2006-01-02T15:04:05.999999999Z07:00")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockClient := NewMockIMongoClient(ctrl)
	mockCollection := NewMockIMongoCollection(ctrl)
	mockSingleResult := NewMockIMongoSingleResult(ctrl)

	repo := &mongodbStore{
		logger,
		mockClient,
		mockCollection,
	}

	p := ProductPrice{
		Name:  "test",
		Price: 3.14,
	}

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"name", p.Name}}
	mockCollection.EXPECT().
		FindOneAndUpdate(ctx, filter, gomock.Any(), opts).
		Return(mockSingleResult)

	mockSingleResult.EXPECT().
		Decode(gomock.AssignableToTypeOf(&p)).
		SetArg(0, p).
		Return(nil)

	err := repo.UpdatePrice(ctx, p)
	if err != nil {
		t.Errorf("bad result, unexpected error, got %v", err)
	}
}
