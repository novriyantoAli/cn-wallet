package dto

import "time"

// CreatePaylaterRepaymentRequest represents a request to create a paylater repayment
type CreatePaylaterRepaymentRequest struct {
	PaylaterLoanID uint   `json:"paylater_loan_id" binding:"required"`
	UserID         uint   `json:"user_id" binding:"required"`
	Amount         int64  `json:"amount" binding:"required,gt=0"`
	PaymentSource  string `json:"payment_source" binding:"required,oneof=wallet external_payment"`
}

// GetPaylaterRepaymentResponse represents paylater repayment information response
type GetPaylaterRepaymentResponse struct {
	ID             uint      `json:"id"`
	PaylaterLoanID uint      `json:"paylater_loan_id"`
	UserID         uint      `json:"user_id"`
	Amount         int64     `json:"amount"`
	PaymentSource  string    `json:"payment_source"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// ListPaylaterRepaymentsRequest represents a request to list repayments with filters
type ListPaylaterRepaymentsRequest struct {
	LoanID   *uint   `form:"loan_id"`
	UserID   *uint   `form:"user_id"`
	Status   *string `form:"status"`
	FromDate *string `form:"from_date"`
	ToDate   *string `form:"to_date"`
	Page     int     `form:"page" binding:"min=1"`
	PageSize int     `form:"page_size" binding:"min=1,max=100"`
}

// ListPaylaterRepaymentsResponse represents paginated repayment list response
type ListPaylaterRepaymentsResponse struct {
	Data       []GetPaylaterRepaymentResponse `json:"data"`
	TotalCount int64                          `json:"total_count"`
	Page       int                            `json:"page"`
	PageSize   int                            `json:"page_size"`
	TotalPages int                            `json:"total_pages"`
}
