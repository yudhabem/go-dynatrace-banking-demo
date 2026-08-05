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

	ctx, span := otel.Tracer("banking-api/service").Start(
		ctx,
		"banking.transfer",
		trace.WithAttributes(
			attribute.String("banking.operation", "transfer"),

			attribute.String("transfer.id", transactionID),
			attribute.String("transfer.from_account", req.FromAccount),
			attribute.String("transfer.to_account", req.ToAccount),
			attribute.Float64("transfer.amount", req.Amount),
		),
	)
	defer span.End()

	logger.Log.Info(
		"transfer started",
		zap.String("transactionId", transactionID),
		zap.String("fromAccount", req.FromAccount),
		zap.String("toAccount", req.ToAccount),
		zap.Float64("amount", req.Amount),
	)

	// ------------------------------------------------------------------
	// Source Account
	// ------------------------------------------------------------------

	from, err := s.repo.FindAccount(ctx, req.FromAccount)
	if err != nil {

		span.RecordError(err)
		span.SetStatus(codes.Error, "source account not found")
		span.SetAttributes(
			attribute.String("transfer.status", "FAILED"),
		)

		logger.Log.Error(
			"source account not found",
			zap.String("transactionId", transactionID),
			zap.String("account", req.FromAccount),
			zap.Error(err),
		)

		return "", err
	}

	logger.Log.Info(
		"source account loaded",
		zap.String("transactionId", transactionID),
		zap.String("account", from.AccountNumber),
		zap.Float64("balance", from.Balance),
	)

	// ------------------------------------------------------------------
	// Destination Account
	// ------------------------------------------------------------------

	to, err := s.repo.FindAccount(ctx, req.ToAccount)
	if err != nil {

		span.RecordError(err)
		span.SetStatus(codes.Error, "destination account not found")
		span.SetAttributes(
			attribute.String("transfer.status", "FAILED"),
		)

		logger.Log.Error(
			"destination account not found",
			zap.String("transactionId", transactionID),
			zap.String("account", req.ToAccount),
			zap.Error(err),
		)

		return "", err
	}

	logger.Log.Info(
		"destination account loaded",
		zap.String("transactionId", transactionID),
		zap.String("account", to.AccountNumber),
		zap.Float64("balance", to.Balance),
	)

	// ------------------------------------------------------------------
	// Validation
	// ------------------------------------------------------------------

	if from.Balance < req.Amount {

		err := errors.New("insufficient balance")

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(
			attribute.String("transfer.status", "FAILED"),
		)

		logger.Log.Warn(
			"insufficient balance",
			zap.String("transactionId", transactionID),
			zap.String("account", from.AccountNumber),
			zap.Float64("balance", from.Balance),
			zap.Float64("requested", req.Amount),
		)

		return "", err
	}

	// ------------------------------------------------------------------
	// Update Balance
	// ------------------------------------------------------------------

	from.Balance -= req.Amount
	to.Balance += req.Amount

	trx := &model.Transaction{
		TransactionID: transactionID,
		FromAccount:   from.AccountNumber,
		ToAccount:     to.AccountNumber,
		Amount:        req.Amount,
		Type:          "TRANSFER",
		Status:        "SUCCESS",
	}

	err = s.repo.ExecuteTransfer(ctx, from, to, trx)
	if err != nil {

		span.RecordError(err)
		span.SetStatus(codes.Error, "transfer failed")
		span.SetAttributes(
			attribute.String("transfer.status", "FAILED"),
		)

		logger.Log.Error(
			"transfer failed",
			zap.String("transactionId", trx.TransactionID),
			zap.Error(err),
		)

		return "", err
	}

	// ------------------------------------------------------------------
	// Success
	// ------------------------------------------------------------------

	span.SetStatus(codes.Ok, "transfer success")

	span.SetAttributes(
		attribute.String("transfer.status", "SUCCESS"),
	)

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

	ctx, span := otel.Tracer("banking-api/service").Start(
		ctx,
		"banking.account.inquiry",
		trace.WithAttributes(
			attribute.String("banking.operation", "inquiry"),
			attribute.String("account.id", account),
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

	ctx, span := otel.Tracer("banking-api/service").Start(
		ctx,
		"banking.transaction.history",
		trace.WithAttributes(
			attribute.String("banking.operation", "transaction_history"),
		),
	)
	defer span.End()

	transactions, err := s.repo.GetTransactions(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "transaction history failed")
	}

	return transactions, err
}
