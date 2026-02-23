package main

import (
	"errors"
	"fmt"

	"git.ghink.net/ghink/refundTool/internal/config"
	"git.ghink.net/ghink/refundTool/internal/logger"
	"git.ghink.net/ghink/refundTool/internal/method/payment"
	"github.com/ghinknet/payutils/v2/model"
	"go.uber.org/zap"
)

func main() {
	// Load public config
	config.LoadPublic()

	// Init logger
	logger.InitLogger()

	// Read space name
	spaceName := ""
	fmt.Println("Please input space name:")
	if _, err := fmt.Scanln(&spaceName); err != nil {
		logger.L.Error("failed to scan input", zap.Error(err))
		return
	}

	// Init payutils
	payment.InitPayutils(spaceName)

	// Main process
	for {
		// Read order ID
		fmt.Println("Please input order id:")
		orderID := ""
		if _, err := fmt.Scanln(&orderID); err != nil {
			logger.L.Error("failed to scan input", zap.Error(err))
			continue
		}

		// Check order ID content
		if orderID == "" {
			logger.L.Error("order id is empty, please retry")
			continue
		}

		// Read payment method
		fmt.Println("Please input payment method:")
		payMethod := ""
		if _, err := fmt.Scanln(&payMethod); err != nil {
			logger.L.Error("failed to scan input", zap.Error(err))
			continue
		}

		payMethodType := model.TradeMethod(payMethod)

		// Read currency
		fmt.Println("Please input currency (default CNY):")
		currency := ""
		_, _ = fmt.Scanln(&currency)
		if currency == "" {
			currency = "CNY"
		}

		// Read totalAmount
		fmt.Println("Please input total amount (in cent):")
		var totalAmount int64 = 0
		if _, err := fmt.Scanln(&totalAmount); err != nil {
			logger.L.Error("failed to scan input", zap.Error(err))
			continue
		}
		if totalAmount <= 0 {
			logger.L.Error("total amount must be greater than 0")
			continue
		}

		// Read refund ID
		fmt.Println("Please input refund ID:")
		refundID := ""
		if _, err := fmt.Scanln(&refundID); err != nil {
			logger.L.Error("failed to scan input", zap.Error(err))
			continue
		}
		if refundID == "" {
			logger.L.Error("refund id is empty, please retry")
			continue
		}

		// Read refund amount
		fmt.Println("Please input refund amount (in cent):")
		var refundAmount int64 = 0
		if _, err := fmt.Scanln(&refundAmount); err != nil {
			logger.L.Error("failed to scan input", zap.Error(err))
			continue
		}
		if refundAmount <= 0 {
			logger.L.Error("refund amount must be greater than 0")
			continue
		}

		// Reason
		fmt.Println("Please input refund reason:")
		reason := ""
		_, _ = fmt.Scanln(&reason)
		if reason == "" {
			reason = "主动退款"
		}

		// Refund
		logger.L.Info("refunding...")
		if err := payment.C.Refund(
			orderID, payMethodType, currency, totalAmount, refundID, refundAmount, reason,
		); err != nil {
			var payutilsError *model.PayutilsError
			if errors.As(err, &payutilsError) {
				upstreamCode, upstreamResponse, upstreamMessage := payutilsError.UpstreamDetail()
				logger.L.Error(
					"failed to refund", zap.Error(err),
					zap.Int("upstreamCode", upstreamCode),
					zap.String("upstreamResponse", upstreamResponse),
					zap.String("upstreamMessage", upstreamMessage),
				)
			} else {
				logger.L.Error("failed to refund", zap.Error(err))
			}
			continue
		}
	}
}
