package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/dto"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/logger"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/model"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/repository"
)

type TransferService struct {
	repo *repository.TransferRepository
}

func NewTransferService(repo *repository.TransferRepository) *TransferService {
	return &TransferService{
		repo: repo,
	}
}

func (s *TransferService) Transfer(ctx context.Context, req dto.TransferRequest) (string, error) {
	transactionID := fmt.Sprintf("TRX-%s", uuid.New().String()[:8])
	ctx, span := otel.Tracer("banking-api/service").Start(ctx, "banking.transfer",
		trace.WithAttributes(
			attribute.String("banking.operation", "transfer"),
			attribute.String("banking.transaction.id", transactionID),
			attribute.String("banking.account.from.id", req.FromAccount),
			attribute.String("banking.account.to.id", req.ToAccount),
			attribute.Float64("banking.transaction.amount", req.Amount),
		),
	)
	defer span.End()

	logger.Log.Info(
		"transfer started",
		zap.String("fromAccount", req.FromAccount),
		zap.String("toAccount", req.ToAccount),
		zap.Float64("amount", req.Amount),
	)

	from, err := s.repo.FindAccount(ctx, req.FromAccount)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "source account not found")
		logger.Log.Error(
			"source account not found",
			zap.String("account", req.FromAccount),
			zap.Error(err),
		)

		return "", err
	}

	logger.Log.Info(
		"source account loaded",
		zap.String("account", from.AccountNumber),
		zap.Float64("balance", from.Balance),
	)

	to, err := s.repo.FindAccount(ctx, req.ToAccount)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "destination account not found")
		logger.Log.Error(
			"destination account not found",
			zap.String("account", req.ToAccount),
			zap.Error(err),
		)

		return "", err
	}

	logger.Log.Info(
		"destination account loaded",
		zap.String("account", to.AccountNumber),
		zap.Float64("balance", to.Balance),
	)

	if from.Balance < req.Amount {
		err := errors.New("insufficient balance")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Log.Warn(
			"insufficient balance",
			zap.String("account", from.AccountNumber),
			zap.Float64("balance", from.Balance),
			zap.Float64("requested", req.Amount),
		)

		return "", err
	}

	from.Balance -= req.Amount
	to.Balance += req.Amount

	trx := &model.Transaction{
		TransactionID: transactionID,
		FromAccount:   from.AccountNumber,
		ToAccount:     to.AccountNumber,
		Amount:        req.Amount,
		Status:        "SUCCESS",
	}

	err = s.repo.ExecuteTransfer(ctx, from, to, trx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "transfer failed")
		logger.Log.Error(
			"transfer failed",
			zap.String("transactionId", trx.TransactionID),
			zap.Error(err),
		)

		return "", err
	}

	logger.Log.Info(
		"transfer completed",
		zap.String("transactionId", trx.TransactionID),
		zap.String("fromAccount", trx.FromAccount),
		zap.String("toAccount", trx.ToAccount),
		zap.Float64("amount", trx.Amount),
		zap.String("status", trx.Status),
	)
	return trx.TransactionID, nil
}

func (s *TransferService) Inquiry(ctx context.Context, account string) (*model.Account, error) {
	ctx, span := otel.Tracer("banking-api/service").Start(ctx, "banking.account.inquiry",
		trace.WithAttributes(
			attribute.String("banking.operation", "inquiry"),
			attribute.String("banking.account.id", account),
		),
	)
	defer span.End()

	acc, err := s.repo.GetAccount(ctx, account)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "account not found")
	}
	return acc, err
}

func (s *TransferService) History(ctx context.Context) ([]model.Transaction, error) {
	ctx, span := otel.Tracer("banking-api/service").Start(ctx, "banking.transaction.history",
		trace.WithAttributes(attribute.String("banking.operation", "transaction_history")),
	)
	defer span.End()

	transactions, err := s.repo.GetTransactions(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "transaction history failed")
	}
	return transactions, err
}
