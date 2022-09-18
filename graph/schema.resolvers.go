package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/google/uuid"
	"github.com/motemen/example-gqlgen-dataloader/db"
	"github.com/motemen/example-gqlgen-dataloader/db/loaders"
	"github.com/motemen/example-gqlgen-dataloader/graph/generated"
	"github.com/motemen/example-gqlgen-dataloader/graph/model"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*db.Todo, error) {
	todo := db.Todo{
		ID:     uuid.NewString(),
		Text:   input.Text,
		Done:   false,
		UserID: input.UserID,
	}
	err := r.DB.Create(&todo).Error
	if err != nil {
		return nil, err
	}
	return &todo, err
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*db.User, error) {
	user := db.User{
		Name: input.Name,
	}
	if input.ID != nil {
		user.ID = *input.ID
	} else {
		user.ID = uuid.NewString()
	}
	err := r.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*db.Todo, error) {
	var todos []*db.Todo
	err := r.DB.Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return todos, nil
}

// User is the resolver for the user field.
func (r *todoResolver) User(ctx context.Context, obj *db.Todo) (*db.User, error) {
	return loaders.GetUser(ctx, obj.UserID)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Todo returns generated.TodoResolver implementation.
func (r *Resolver) Todo() generated.TodoResolver { return &todoResolver{r} }

type (
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
	todoResolver     struct{ *Resolver }
)
