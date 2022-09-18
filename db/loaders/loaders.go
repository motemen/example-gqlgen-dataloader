package loaders

import (
	"context"
	"fmt"
	"net/http"

	"github.com/graph-gophers/dataloader"
	"github.com/motemen/example-gqlgen-dataloader/db"
)

var contextKey = &struct{ name string }{"loaders"}

type loaders struct {
	userLoader *dataloader.Loader
}

func Middleware(dbConn *db.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		br := &batchReader{db: dbConn}
		ctx := context.WithValue(r.Context(), contextKey, &loaders{
			userLoader: dataloader.NewBatchedLoader(br.GetUsers),
		})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func GetUser(ctx context.Context, userID string) (*db.User, error) {
	loaders := ctx.Value(contextKey).(*loaders)
	thunk := loaders.userLoader.Load(ctx, dataloader.StringKey(userID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.(*db.User), nil
}

type batchReader struct {
	db *db.DB
}

var _ dataloader.BatchFunc = ((*batchReader)(nil)).GetUsers

func (br *batchReader) GetUsers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var users []db.User
	err := br.db.Find(&users, keys).Error
	if err != nil {
		panic(err)
	}

	userById := map[string]db.User{}
	for _, user := range users {
		userById[user.ID] = user
	}

	result := make([]*dataloader.Result, len(keys))
	for i, key := range keys {
		if user, ok := userById[key.String()]; ok {
			result[i] = &dataloader.Result{Data: &user}
		} else {
			result[i] = &dataloader.Result{Error: fmt.Errorf("user not found: %s", key)}
		}
	}

	return result
}
