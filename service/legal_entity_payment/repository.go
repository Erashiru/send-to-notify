package legal_entity_payment

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kwaaka-team/orders-core/service/legal_entity_payment/models"
	"strconv"
	"strings"
	"time"
)

type Repository interface {
	Create(ctx context.Context, payment models.LegalEntityPayment) (string, error)
	GetByID(ctx context.Context, paymentID string) (models.LegalEntityPayment, error)
	GetList(ctx context.Context, req models.ListLegalEntityPaymentQuery) ([]models.LegalEntityPayment, error)
	Update(ctx context.Context, payment models.UpdateLegalEntityPayment) error
	Delete(ctx context.Context, paymentID string) error
	GetPaidPaymentsAnalytics(ctx context.Context, req models.LegalEntityPaymentAnalyticsRequest) (models.LegalEntityPaymentAnalytics, error)
	GetUnpaidPaymentsAnalytics(ctx context.Context, req models.LegalEntityPaymentAnalyticsRequest) (models.LegalEntityPaymentAnalytics, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) (*PostgresRepository, error) {
	return &PostgresRepository{
		db: db,
	}, nil
}

func (r *PostgresRepository) Create(ctx context.Context, req models.LegalEntityPayment) (string, error) {
	query := `INSERT INTO payment (name, legal_entity_id, legal_entity_name, 
                     amount, start_date, end_date, payment_type, status,
                     bill, bill_payment, billing_at, bill_payment_at, confirm_payment_at,
                     created_at, updated_at) 
			VALUES ($1, CAST($2 AS integer), $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) returning id`
	if _, err := r.db.Prepare(query); err != nil {
		return "", err
	}

	var id string
	req.CreatedAt = time.Now().UTC()

	row := r.db.QueryRowContext(ctx, query, req.Name, req.LegalEntityID, req.LegalEntityName,
		req.Amount, req.StartDate, req.EndDate, req.PaymentType, req.Status,
		req.Bill, req.BillPayment, req.BillingAt, req.BillPaymentAt, req.ConfirmPaymentAt,
		req.CreatedAt, req.UpdatedAt,
	)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (models.LegalEntityPayment, error) {
	if id == "" {
		return models.LegalEntityPayment{}, fmt.Errorf("legal entity payment id is empty")
	}

	query := `SELECT id, COALESCE(name, '') AS name, legal_entity_id, legal_entity_name, 
                     amount, start_date, end_date, payment_type, status,
                     bill, bill_payment, billing_at, bill_payment_at, confirm_payment_at,
                     created_at, updated_at 
					FROM payment WHERE id = $1`
	_, err := r.db.Prepare(query)
	if err != nil {
		return models.LegalEntityPayment{}, err
	}
	var payment models.LegalEntityPayment

	if err = r.db.QueryRowContext(ctx, query, id).Scan(&payment.ID, &payment.Name, &payment.LegalEntityID, &payment.LegalEntityName,
		&payment.Amount, &payment.StartDate, &payment.EndDate, &payment.PaymentType, &payment.Status,
		&payment.Bill, &payment.BillPayment, &payment.BillingAt, &payment.BillPaymentAt, &payment.ConfirmPaymentAt,
		&payment.CreatedAt, &payment.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return models.LegalEntityPayment{}, fmt.Errorf("legal entity payment with id %s not exist", id)
		}

		return models.LegalEntityPayment{}, err
	}

	return payment, nil
}

func (r *PostgresRepository) GetList(ctx context.Context, req models.ListLegalEntityPaymentQuery) ([]models.LegalEntityPayment, error) {
	fields, queryParams, _ := r.filterForList(req)

	query := `SELECT id, contract, COALESCE(name, '') AS name, legal_entity_id, legal_entity_name,
                     amount, paid_amount, start_date, end_date, payment_type, status,
                     bill, bill_payment, billing_at, bill_payment_at, confirm_payment_at,
                     comments, created_at, updated_at
				FROM legal_entity_payment WHERE 1=1 ` + fields

	if _, err := r.db.Prepare(query); err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.LegalEntityPayment
	for rows.Next() {
		var payment models.LegalEntityPayment

		if err := rows.Scan(&payment.ID, &payment.Contract, &payment.Name, &payment.LegalEntityID, &payment.LegalEntityName,
			&payment.Amount, &payment.PaidAmount, &payment.StartDate, &payment.EndDate, &payment.PaymentType, &payment.Status,
			&payment.Bill, &payment.BillPayment, &payment.BillingAt, &payment.BillPaymentAt, &payment.ConfirmPaymentAt,
			&payment.Comments, &payment.CreatedAt, &payment.UpdatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *PostgresRepository) filterForList(req models.ListLegalEntityPaymentQuery) (fields string, queryParams []interface{}, index int) {
	index = 1

	fields, queryParams, index = r.setTime(req.StartDate, req.EndDate, fields, queryParams, index)

	if len(req.LegalEntityIDs) > 0 {
		fields += " AND legal_entity_id IN ("
		for i, id := range req.LegalEntityIDs {
			fields += "$" + strconv.Itoa(index)
			index++
			queryParams = append(queryParams, id)
			if i < len(req.LegalEntityIDs)-1 {
				fields += ","
			}
		}
		fields += ")"
	}

	if len(req.PaymentTypes) > 0 {
		fields += " AND payment_type IN ("
		for i, t := range req.PaymentTypes {
			fields += "$" + strconv.Itoa(index)
			index++
			queryParams = append(queryParams, t)
			if i < len(req.PaymentTypes)-1 {
				fields += ","
			}
		}
		fields += ")"
	}

	if len(req.Statuses) > 0 {
		fields += " AND status IN ("
		for i, status := range req.Statuses {
			fields += "$" + strconv.Itoa(index)
			index++
			queryParams = append(queryParams, status)
			if i < len(req.Statuses)-1 {
				fields += ","
			}
		}
		fields += ")"
	}

	fields += " ORDER BY created_at DESC"

	if req.Limit > 0 {
		fields += fmt.Sprintf(" LIMIT $%d", index)
		index++
		queryParams = append(queryParams, req.Limit)

		fields += fmt.Sprintf(" OFFSET $%d", index)
		index++
		queryParams = append(queryParams, req.Pagination.AddOffset().Offset)
	}

	return fields, queryParams, index
}

func (r *PostgresRepository) Update(ctx context.Context, payment models.UpdateLegalEntityPayment) error {
	fields, queryParams, counter := r.setFields(payment)

	query := fmt.Sprintf(`UPDATE payment SET %s WHERE id = CAST($%d AS integer)`, strings.Join(fields, ", "), counter)

	if _, err := r.db.Prepare(query); err != nil {
		return err
	}

	_, err := r.db.ExecContext(ctx, query, queryParams...)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, paymentID string) error {
	query := `DELETE FROM payment WHERE id = CAST($1 AS integer)`

	if _, err := r.db.Prepare(query); err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx, query, paymentID); err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) setFields(payment models.UpdateLegalEntityPayment) (fields []string, queryParams []interface{}, index int) {
	index = 1

	if payment.Name != nil {
		fields = append(fields, fmt.Sprintf("name = $%d", index))
		queryParams = append(queryParams, *payment.Name)
		index++
	}

	if payment.LegalEntityID != nil {
		fields = append(fields, fmt.Sprintf("legal_entity_id = CAST($%d AS integer)", index))
		queryParams = append(queryParams, *payment.LegalEntityID)
		index++
	}

	if payment.LegalEntityName != nil {
		fields = append(fields, fmt.Sprintf("legal_entity_name = $%d", index))
		queryParams = append(queryParams, *payment.LegalEntityName)
		index++
	}

	if payment.Amount != nil {
		fields = append(fields, fmt.Sprintf("amount = $%d", index))
		queryParams = append(queryParams, *payment.Amount)
		index++
	}

	if payment.StartDate != nil {
		fields = append(fields, fmt.Sprintf("start_date = $%d", index))
		queryParams = append(queryParams, *payment.StartDate)
		index++
	}

	if payment.EndDate != nil {
		fields = append(fields, fmt.Sprintf("end_date = $%d", index))
		queryParams = append(queryParams, *payment.EndDate)
		index++
	}

	if payment.PaymentType != nil {
		fields = append(fields, fmt.Sprintf("payment_type = $%d", index))
		queryParams = append(queryParams, *payment.PaymentType)
		index++
	}

	if payment.Status != nil {
		fields = append(fields, fmt.Sprintf("status = $%d", index))
		queryParams = append(queryParams, *payment.Status)
		index++
	}

	if payment.Bill != nil {
		fields = append(fields, fmt.Sprintf("bill = $%d", index))
		queryParams = append(queryParams, *payment.Bill)
		index++
	}

	if payment.BillPayment != nil {
		fields = append(fields, fmt.Sprintf("bill_payment = $%d", index))
		queryParams = append(queryParams, *payment.BillPayment)
		index++
	}

	if payment.BillingAt != nil {
		fields = append(fields, fmt.Sprintf("billing_at = $%d", index))
		queryParams = append(queryParams, *payment.BillingAt)
		index++
	}

	if payment.BillPaymentAt != nil {
		fields = append(fields, fmt.Sprintf("bill_payment_at = $%d", index))
		queryParams = append(queryParams, *payment.BillPaymentAt)
		index++
	}

	if payment.ConfirmPaymentAt != nil {
		fields = append(fields, fmt.Sprintf("confirm_payment_at = $%d", index))
		queryParams = append(queryParams, *payment.ConfirmPaymentAt)
		index++
	}

	updatedAt := time.Now().UTC()
	fields = append(fields, fmt.Sprintf("updated_at = $%d", index))
	queryParams = append(queryParams, &updatedAt)
	index++

	queryParams = append(queryParams, *payment.ID)

	return fields, queryParams, index
}

func (r *PostgresRepository) GetPaidPaymentsAnalytics(ctx context.Context, req models.LegalEntityPaymentAnalyticsRequest) (models.LegalEntityPaymentAnalytics, error) {
	statuses := []string{models.PAID.String(), models.PAID_CONFIRMED.String()}
	req.Statuses = statuses

	fields, queryParams, _ := r.filterForAnalytics(req)

	query := `SELECT COALESCE(COUNT(id),0), COALESCE(SUM(amount),0) FROM payment WHERE 1=1 ` + fields

	if _, err := r.db.Prepare(query); err != nil {
		return models.LegalEntityPaymentAnalytics{}, err
	}

	var resp models.LegalEntityPaymentAnalytics

	if err := r.db.QueryRowContext(ctx, query, queryParams...).Scan(&resp.Quantity, &resp.Amount); err != nil {
		return models.LegalEntityPaymentAnalytics{}, err
	}

	return resp, nil
}

func (r *PostgresRepository) GetUnpaidPaymentsAnalytics(ctx context.Context, req models.LegalEntityPaymentAnalyticsRequest) (models.LegalEntityPaymentAnalytics, error) {
	statuses := []string{models.UNBILLED.String(), models.BILLED.String(), models.FAILED.String(), ""}
	req.Statuses = statuses

	fields, queryParams, _ := r.filterForAnalytics(req)

	query := `SELECT COALESCE(COUNT(id),0), COALESCE(SUM(amount),0) FROM payment WHERE 1=1 ` + fields

	if _, err := r.db.Prepare(query); err != nil {
		return models.LegalEntityPaymentAnalytics{}, err
	}

	var resp models.LegalEntityPaymentAnalytics

	if err := r.db.QueryRowContext(ctx, query, queryParams...).Scan(&resp.Quantity, &resp.Amount); err != nil {
		return models.LegalEntityPaymentAnalytics{}, err
	}

	return resp, nil
}

func (r *PostgresRepository) filterForAnalytics(req models.LegalEntityPaymentAnalyticsRequest) (fields string, queryParams []interface{}, index int) {
	index = 1

	fields, queryParams, index = r.setTime(req.StartDate, req.EndDate, fields, queryParams, index)

	if len(req.LegalEntityIDs) > 0 {
		fields += " AND legal_entity_id IN ("
		for i, id := range req.LegalEntityIDs {
			fields += "$" + strconv.Itoa(index)
			index++
			queryParams = append(queryParams, id)
			if i < len(req.LegalEntityIDs)-1 {
				fields += ","
			}
		}
		fields += ")"
	}

	if len(req.Statuses) > 0 {
		fields += " AND status IN ("
		for i, status := range req.Statuses {
			fields += "$" + strconv.Itoa(index)
			index++
			queryParams = append(queryParams, status)
			if i < len(req.Statuses)-1 {
				fields += ","
			}
		}
		fields += ")"
	}

	return fields, queryParams, index
}

func (r *PostgresRepository) setTime(startDate, endDate time.Time, fields string, queryParams []interface{}, index int) (string, []interface{}, int) {

	if !startDate.IsZero() {
		fields += fmt.Sprintf(" AND created_at >= $%d", index)
		index++
		queryParams = append(queryParams, startDate)
	}

	if !endDate.IsZero() {
		fields += fmt.Sprintf(" AND created_at <= $%d", index)
		index++
		queryParams = append(queryParams, endDate)
	}

	return fields, queryParams, index
}
