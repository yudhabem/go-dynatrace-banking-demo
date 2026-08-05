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

type PaymentService struct {
	repo *repository.PaymentRepository
}

func NewPaymentService(repo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{
		repo: repo,
	}
}

func (s *PaymentService) Payment(ctx context.Context, req dto.PaymentRequest) (string, error) {

	transactionID := fmt.Sprintf("PAY-%s", uuid.New().String()[:8])

	ctx, span := otel.Tracer("banking-api/service").Start(
		ctx,
		"banking.payment",
		trace.WithAttributes(
			attribute.String("banking.operation", "payment"),

			attribute.String("payment.id", transactionID),
			attribute.String("payment.account", req.AccountNumber),
			attribute.String("payment.merchant", req.Merchant),
			attribute.Float64("payment.amount", req.Amount),
		),
	)
	defer span.End()

	logger.Log.Info(
		"payment started",
		zap.String("transactionId", transactionID),
		zap.String("account", req.AccountNumber),
		zap.String("merchant", req.Merchant),
		zap.Float64("amount", req.Amount),
	)

	account, err := s.repo.GetAccount(ctx, req.AccountNumber)
	if err != nil {

		span.RecordError(err)
		span.SetStatus(codes.Error, "payment account not found")
		span.SetAttributes(
			attribute.String("payment.status", "FAILED"),
		)

		logger.Log.Error(
			"payment account not found",
			zap.String("transactionId", transactionID),
			zap.String("account", req.AccountNumber),
			zap.Error(err),
		)

		return "", err
	}

	logger.Log.Info(
		"payment account loaded",
		zap.String("transactionId", transactionID),
		zap.String("account", account.AccountNumber),
		zap.Float64("balance", account.Balance),
	)

	if account.Balance < req.Amount {

		err := errors.New("insufficient balance")

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(
			attribute.String("payment.status", "FAILED"),
		)

		logger.Log.Warn(
			"insufficient balance",
			zap.String("transactionId", transactionID),
			zap.String("account", account.AccountNumber),
			zap.Float64("balance", account.Balance),
			zap.Float64("requested", req.Amount),
		)

		return "", err
	}

	account.Balance -= req.Amount

	logger.Log.Info(
		"payment balance deducted",
		zap.String("transactionId", transactionID),
		zap.String("account", account.AccountNumber),
		zap.Float64("amount", req.Amount),
		zap.Float64("balance", account.Balance),
	)

	trx := &model.Transaction{
		TransactionID: transactionID,
		FromAccount:   account.AccountNumber,
		Amount:        req.Amount,
		Type:          "PAYMENT",
		Merchant:      req.Merchant,
		Status:        "SUCCESS",
	}

	if err := s.repo.ExecutePayment(ctx, account, trx); err != nil {

		span.RecordError(err)
		span.SetStatus(codes.Error, "payment failed")
		span.SetAttributes(
			attribute.String("payment.status", "FAILED"),
		)

		logger.Log.Error(
			"payment failed",
			zap.String("transactionId", trx.TransactionID),
			zap.Error(err),
		)

		return "", err
	}

	span.SetStatus(codes.Ok, "payment success")
	span.SetAttributes(
		attribute.String("payment.status", "SUCCESS"),
	)

	logger.Log.Info(
		"payment completed",
		zap.String("transactionId", trx.TransactionID),
		zap.String("account", trx.FromAccount),
		zap.String("merchant", trx.Merchant),
		zap.Float64("amount", trx.Amount),
		zap.String("status", trx.Status),
	)

	return trx.TransactionID, nil
}
