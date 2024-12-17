package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	"net/http"
	"time"
)

func (server *Server) KaspiSaleScoutCronEvent(c *gin.Context) {
	unpaidOrders, err := server.paymentManager.GetUnpaidPaymentsByPaymentSystem(c.Request.Context(), 30, "kaspi_salescout")
	if err != nil {
		server.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	server.Logger.Infof("unpaid payments for kaspi salescout orders: %s", unpaidOrders)

	seconds := 0
	for seconds < 55 {
		server.processingKaspiSaleScoutUnpaidOrders(c.Request.Context(), unpaidOrders)
		time.Sleep(3 * time.Second)
		seconds = seconds + 3
	}

	c.Status(http.StatusOK)
}

func (server *Server) processingKaspiSaleScoutUnpaidOrders(ctx context.Context, unpaidOrders []models.PaymentOrder) {
	for i := range unpaidOrders {
		unpaidOrder := unpaidOrders[i]

		status, err := server.paymentManager.GetKaspiSaleScoutPaymentStatus(ctx, unpaidOrder.PaymentOrderID)
		if err != nil {
			server.Logger.Error(err)
			continue
		}
		if status != "Processed" {
			continue
		}

		_, err = server.paymentManager.SavePaymentDetailsEvent(ctx, models.PaymentOrder{
			PaymentOrderID:       unpaidOrder.PaymentOrderID,
			PaymentOrderStatus:   models.PAID,
			PaymentStatusHistory: []models.StatusHistory{{Status: models.PAID}},
			CartID:               unpaidOrder.CartID,
		})

		if err != nil {
			server.Logger.Error(err)
			continue
		}
	}
}
