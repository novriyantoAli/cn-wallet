package dto

import "time"

// CreatePaylaterLoanRequest represents a request to create a paylater loan
type CreatePaylaterLoanRequest struct {
	UserID   uint      `json:"user_id" binding:"required"`
	Amount   int64     `json:"amount" binding:"required,gt=0"`
	Interest int64     `json:"interest" binding:"required,gte=0"`
	DueDate  time.Time `json:"due_date" binding:"required"`
	Source   string    `json:"source" binding:"required,oneof=transfer checkout"`
}

// UpdateLoanStatusRequest represents a request to update loan status
type UpdateLoanStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending active paid overdue defaulted"`
}

// GetPaylaterLoanResponse represents paylater loan information response
type GetPaylaterLoanResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Amount    int64     `json:"amount"`
	Interest  int64     `json:"interest"`
	Total     int64     `json:"total"`
	DueDate   time.Time `json:"due_date"`
	Source    string    `json:"source"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ListPaylaterLoansRequest represents a request to list loans with filters
type ListPaylaterLoansRequest struct {
	UserID   *uint   `form:"user_id"`
	Status   *string `form:"status"`
	Source   *string `form:"source"`
	FromDate *string `form:"from_date"`
	ToDate   *string `form:"to_date"`
	Page     int     `form:"page" binding:"min=1"`
	PageSize int     `form:"page_size" binding:"min=1,max=100"`
}

// ListPaylaterLoansResponse represents paginated loan list response
type ListPaylaterLoansResponse struct {
	Data       []GetPaylaterLoanResponse `json:"data"`
	TotalCount int64                     `json:"total_count"`
	Page       int                       `json:"page"`
	PageSize   int                       `json:"page_size"`
	TotalPages int                       `json:"total_pages"`
}

// PaylaterLoanStatsResponse represents loan statistics for a user
type PaylaterLoanStatsResponse struct {
	UserID           uint  `json:"user_id"`
	TotalLoans       int64 `json:"total_loans"`
	ActiveLoans      int64 `json:"active_loans"`
	TotalBorrowed    int64 `json:"total_borrowed"`
	TotalPaid        int64 `json:"total_paid"`
	TotalOutstanding int64 `json:"total_outstanding"`
	OverdueLoans     int64 `json:"overdue_loans"`
}
