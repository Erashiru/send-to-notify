package offline_orders

import (
	"context"
	"database/sql"
	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
	"github.com/rs/zerolog/log"
)

type Service struct {
	db *sql.DB
}

func NewOfflineOrdersService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) SaveOrders(ctx context.Context, offlineOrders []models.ExtendedOrderEvent) error {

	err := s.insertOfflineOrders(s.db, offlineOrders)
	if err != nil {
		log.Err(err).Msgf("error inserting message %s", err.Error())
		return err
	}
	log.Info().Msgf("Saving offline order finished")
	return nil
}

func (s *Service) insertOfflineOrders(postgres *sql.DB, orders []models.ExtendedOrderEvent) error {
	tx, err := postgres.Begin()
	if err != nil {
		return err
	}
	for _, order := range orders {
		var count int
		query := "SELECT COUNT(*) FROM orderevents WHERE eventId = $1"
		err := s.db.QueryRow(query, order.EventId).Scan(&count)
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		orderEventID, err := insertOrderEvent(tx, order)
		if err != nil {
			log.Err(err).Msgf("order %v", order)
			_ = tx.Rollback()
			return err
		}
		var customerCount int

		if order.RegularCustomerInfo.Id != "" {
			query = "SELECT COUNT(*) FROM customers WHERE customerId = $1"
			err = s.db.QueryRow(query, order.RegularCustomerInfo.Id).Scan(&customerCount)
			if err != nil {
				return err
			}
		}

		if customerCount == 0 {
			err = insertCustomer(tx, order.Order.Customer, orderEventID, order.RegularCustomerInfo)
			if err != nil {
				log.Err(err).Msgf("customer %v", order.Order.Customer)
				_ = tx.Rollback()
				return err
			}
		}

		for _, item := range order.Order.Items {
			err := insertItem(tx, item, orderEventID)
			if err != nil {
				log.Err(err).Msgf("item %v", item)
				_ = tx.Rollback()
				return err
			}
		}
		if len(order.Order.OfflineOrderPayment) > 0 {
			err = insertPaymentType(tx, order.Order.OfflineOrderPayment[0], orderEventID)
			if err != nil {
				log.Err(err).Msgf("payment type %v", order.Order.OfflineOrderPayment[0].PaymentTypes.Id)
				_ = tx.Rollback()
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func insertItem(tx *sql.Tx, item models.Item, orderEventID int) error {
	stmt := `
  INSERT INTO Items (
   ProductId, Price, PositionID, Type, Amount,
   ProductSizeID, Comment, OrderEventID, TableProductID, Name
  ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := tx.Exec(
		stmt,
		item.ProductId, item.Price, item.PositionId, item.Type, item.Amount,
		item.ProductSizeID, item.Comment, orderEventID, item.TableProduct.Id, item.TableProduct.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

func insertCustomer(tx *sql.Tx, customer models.Customer, orderEventId int, regularCustomerInfo models.GetCustomerInfoResponse) error {

	stmt := `
  INSERT INTO Customers (
   Name, Surname, Email,
   Gender, InBlacklist,
    BlacklistReason, Type, OrderEventId, 
    CustomerId, Comment, Birthday, 
    Anonymized, ReferrerId, ConsentStatus,
    UserData, ShouldReceiveLoyaltyInfo, IsDeleted
  ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	res, err := tx.Exec(
		stmt, customer.Name,
		customer.Surname,
		customer.Email,
		customer.Gender,
		customer.InBlacklist,
		customer.BlacklistReason,
		customer.Type,
		orderEventId,
		regularCustomerInfo.Id,
		regularCustomerInfo.Comment,
		regularCustomerInfo.Birthday,
		regularCustomerInfo.Anonymized,
		regularCustomerInfo.ReferrerId,
		regularCustomerInfo.ConsentStatus,
		regularCustomerInfo.UserData,
		regularCustomerInfo.ShouldReceiveLoyaltyInfo,
		regularCustomerInfo.IsDeleted,
	)
	if err != nil {
		return err
	}

	customerId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	if regularCustomerInfo.Id != "" {
		for _, card := range regularCustomerInfo.Cards {
			err := insertCustomerCard(tx, card, int(customerId))
			if err != nil {
				return err
			}
		}
		for _, walletBalance := range regularCustomerInfo.WalletBalances {
			err := insertWalletBalance(tx, walletBalance, int(customerId))
			if err != nil {
				return err
			}
		}
		for _, category := range regularCustomerInfo.Categories {
			err := insertCustomerCategory(tx, category, int(customerId))
			if err != nil {
				return err
			}
		}
	}
	return err
}

func insertOrderEvent(tx *sql.Tx, order models.ExtendedOrderEvent) (int, error) {
	stmt := `
      INSERT INTO OrderEvents (
         TerminalID, Phone, Status, CompleteBefore, WhenCreated, CookingStartTime, IsDeleted, Sum, Number, StoreId, StoreName, EventId
      ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
      RETURNING ID
   `

	sum, err := order.Order.Sum.Float64()
	if err != nil {
		return 0, err
	}
	res, err := tx.Exec(
		stmt,
		order.Order.TerminalID, order.Order.Phone,
		order.Order.Status, order.Order.CompleteBefore.String(),
		order.Order.WhenCreated.String(), order.Order.CookingStartTime.String(),
		order.Order.IsDeleted, sum, order.Order.Number,
		order.StoreId, order.StoreName, order.EventId,
	)
	if err != nil {
		return 0, err
	}

	orderEventID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(orderEventID), nil
}

func insertPaymentType(tx *sql.Tx, payment models.OfflineOrderPayment, orderEventId int) error {
	stmt := `
  INSERT INTO Payment (
  PaymentTypeId, Kind, Sum, IsProcessedExternally,orderEventID
  ) VALUES ($1, $2, $3, $4, $5)
 `

	_, err := tx.Exec(
		stmt,
		payment.PaymentTypes.Id, payment.PaymentTypes.Kind, payment.Sum,
		payment.IsProcessedExternally, orderEventId,
	)

	if err != nil {
		return err
	}
	return nil
}

func insertCustomerCard(tx *sql.Tx, card models.CustomerCard, customerId int) error {
	stmt := `INSERT INTO Cards (CardId, Track, Number, ValidToDate, CustomerId) VALUES ($1, $2, $3, $4, $5)`

	_, err := tx.Exec(stmt, card.Id, card.Track, card.Number, card.ValidToDate, customerId)
	if err != nil {
		return err
	}
	return nil
}

func insertWalletBalance(tx *sql.Tx, walletBalance models.WalletBalance, customerId int) error {
	stmt := `INSERT INTO WalletBalances (WalletId, Name, Type, Balance, CustomerId) VALUES ($1, $2, $3, $4, $5)`

	_, err := tx.Exec(stmt, walletBalance.Id, walletBalance.Name, walletBalance.Type, walletBalance.Balance, customerId)
	if err != nil {
		return err
	}
	return nil
}

func insertCustomerCategory(tx *sql.Tx, category models.CustomerCategory, customerId int) error {
	stmt := `INSERT INTO CustomerCategories (CategoryId, Name, IsActive, IsDefaultForNewGuests, CustomerId) VALUES ($1, $2, $3, $4, $5)`

	_, err := tx.Exec(stmt, category.Id, category.Name, category.IsActive, category.IsDefaultForNewGuests, customerId)
	if err != nil {
		return err
	}
	return nil
}
