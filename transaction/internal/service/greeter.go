package service

import (
	"context"

	v1 "transaction/api/helloworld/v1"
	"transaction/internal/biz"
)

type TransactionService struct {
	v1.UnimplementedCardServer
	uc *biz.TransactionUsecase
}

func NewTransactionService(uc *biz.TransactionUsecase) *TransactionService {
	return &TransactionService{uc: uc}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, req *v1.CreateTransactionRequest) (*v1.CreateTransactionReply, error) {
	transaction := &biz.Transaction{
		TransactionId:        req.Transactionid,
		AccountId:           req.AccountId,
		RelatedTransactionId: req.RelatedTransactionId,
		TransactionType:     req.TransactionType,
		Amount:              float64(req.Amount),
		Currency:            req.Currency,
		Status:              req.Status,
		Description:         req.Description,
		ReferenceNumber:     req.ReferenceNumber,
		PostingDate:         req.PostingDate,
	}

	t, err := s.uc.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return &v1.CreateTransactionReply{
		TransactionId:        t.TransactionId,
		AccountId:           t.AccountId,
		RelatedTransactionId: t.RelatedTransactionId,
		TransactionType:     t.TransactionType,
		Amount:              float64(t.Amount),
		Currency:            t.Currency,
		Status:              t.Status,
		Description:         t.Description,
		ReferenceNumber:     t.ReferenceNumber,
		PostingDate:         t.PostingDate,
		Message:             "Transaction Details are stored",
	}, nil
}

func (s *TransactionService) UpdateTransaction(ctx context.Context, req *v1.UpdateTransactionRequest) (*v1.UpdateTransactionReply, error) {
	transaction := &biz.Transaction{
		TransactionId:        req.TransactionId,
		AccountId:           req.AccountId,
		RelatedTransactionId: req.RelatedTransactionId,
		TransactionType:     req.TransactionType,
		Amount:              float64(req.Amount),
		Currency:            req.Currency,
		Status:              req.Status,
		Description:         req.Description,
		ReferenceNumber:     req.ReferenceNumber,
		PostingDate:         req.PostingDate,
	}

	updatedTransaction, err := s.uc.UpdateTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateTransactionReply{
		Message: "Transaction Details are updated",
		Display: []*v1.CreateTransactionReply{
			{
				TransactionId:        updatedTransaction.TransactionId,
				AccountId:           updatedTransaction.AccountId,
				RelatedTransactionId: updatedTransaction.RelatedTransactionId,
				TransactionType:     updatedTransaction.TransactionType,
				Amount:              float64(updatedTransaction.Amount),
				Currency:            updatedTransaction.Currency,
				Status:              updatedTransaction.Status,
				Description:         updatedTransaction.Description,
				ReferenceNumber:     updatedTransaction.ReferenceNumber,
				PostingDate:         updatedTransaction.PostingDate,
			},
		},
	}, nil
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, req *v1.DeleteTransactionRequest) (*v1.DeleteTransactionReply, error) {
	err := s.uc.DeleteTransaction(ctx, req.TransactionId)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteTransactionReply{
		Message: "Transaction Details are deleted",
	}, nil
}

func (s *TransactionService) DisplayTransaction(ctx context.Context, req *v1.FindTransactionByIdRequest) (*v1.FindTransactionByIdReply, error) {
	transaction, err := s.uc.GetTransactionById(ctx, req.TransactionId)
	if err != nil {
		return nil, err
	}

	reply := &v1.CreateTransactionReply{
		TransactionId:        transaction.TransactionId,
		AccountId:           transaction.AccountId,
		RelatedTransactionId: transaction.RelatedTransactionId,
		TransactionType:     transaction.TransactionType,
		Amount:              float64(transaction.Amount),
		Currency:            transaction.Currency,
		Status:              transaction.Status,
		Description:         transaction.Description,
		ReferenceNumber:     transaction.ReferenceNumber,
		PostingDate:         transaction.PostingDate,
	}

	return &v1.FindTransactionByIdReply{
		Display: []*v1.CreateTransactionReply{reply},
	}, nil
}