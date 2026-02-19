package payment

import (
	"time"

	"git.ghink.net/ghink/refundTool/internal/config"
	"git.ghink.net/ghink/refundTool/internal/logger"
	payutilsFiber "github.com/ghinknet/payutils/v2/framework/fiber"
	payutilsModel "github.com/ghinknet/payutils/v2/model"
	payutilsAlipay "github.com/ghinknet/payutils/v2/payment/alipay"
	payutilsWeChat "github.com/ghinknet/payutils/v2/payment/wechat"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

var C *payutilsFiber.Client

// InitPayutils inits global payment utils client
func InitPayutils(spaceName string) {
	// Load config in need
	var AlipayConfig *payutilsAlipay.Config = nil
	var WeChatPayConfig *payutilsWeChat.Config = nil
	if config.C.GetString(spaceName+".alipay.appID") != "" {
		AlipayConfig = &payutilsAlipay.Config{
			AppID:             config.C.GetString(spaceName + ".alipay.appID"),
			AppCertPrivateKey: config.C.GetString(spaceName + ".alipay.appCertPrivateKey"),
			AppCert:           config.C.GetString(spaceName + ".alipay.appCert"),
			RootCert:          config.C.GetString(spaceName + ".alipay.rootCert"),
			PublicCert:        config.C.GetString(spaceName + ".alipay.publicCert"),
			IsProd:            config.C.GetBool(spaceName + ".alipay.isProd"),
		}
	}
	if config.C.GetString(spaceName+".wechatPay.appID") != "" {
		WeChatPayConfig = &payutilsWeChat.Config{
			AppID:                    config.C.GetString(spaceName + ".wechatPay.appID"),
			AppSecret:                config.C.GetString(spaceName + ".wechatPay.appSecret"),
			MerchantID:               config.C.GetString(spaceName + ".wechatPay.merchantID"),
			MerchantAPIv3Key:         config.C.GetString(spaceName + ".wechatPay.merchantAPIv3Key"),
			MerchantCertSerialNumber: config.C.GetString(spaceName + ".wechatPay.merchantCertSerialNumber"),
			MerchantPrivateKey:       config.C.GetString(spaceName + ".wechatPay.merchantPrivateKey"),
			PublicKey:                config.C.GetString(spaceName + ".wechatPay.publicKey"),
			PublicKeyID:              config.C.GetString(spaceName + ".wechatPay.publicKeyID"),
		}
	}

	// Construct config of payutils fiber
	Config := payutilsFiber.Config{
		Basic: payutilsModel.Config{
			Debug:        config.Debug,
			AllowOrigins: make([]string, 0),
			Endpoint:     "https://www.ghink.net",
			Prefix:       config.C.GetString(spaceName + ".prefix"),
			Suffix:       config.C.GetString(spaceName + ".suffix"),
		},
		Alipay:    AlipayConfig,
		WeChatPay: WeChatPayConfig,
		Fiber:     fiber.New().Group("/payutils"),
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return nil
		},
		DetailProvider: func(c fiber.Ctx, orderID string, method payutilsModel.TradeMethod) (payutilsModel.OrderDetail, error) {
			return payutilsModel.OrderDetail{}, nil
		},
		StatusUpdater: func(c fiber.Ctx, orderID string, status payutilsModel.TradeState, method payutilsModel.TradeMethod, tm time.Time) error {
			return nil
		},
	}

	var err error
	C, err = payutilsFiber.CreateClient(Config)
	if err != nil {
		logger.L.Fatal("failed to create client payutils with Fiber", zap.Error(err))
	}

	logger.L.Debug("Payutils initialized")
}
