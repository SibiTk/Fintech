package handler

import (
	v1 "BlockAccount/api/helloworld/v1"
	"BlockAccount/internal/data"
	accountpb "account/api/helloworld/v1"
	"context"
)

type AccountBlockHandler struct {
	repo          *data.GreeterRepo
	AccountClient accountpb.AccountClient
}

func (h *AccountBlockHandler) GetAccBlock(ctx context.Context, req *v1.GetAccBlockRequest) (*v1.GetAccBlockReply, error) {
	panic("unimplemented")
}

func NewAccountBlockHandler(
	repo *data.GreeterRepo,
	AccountClient accountpb.AccountClient,
) *AccountBlockHandler {
	return &AccountBlockHandler{repo: repo, AccountClient: AccountClient}
}
