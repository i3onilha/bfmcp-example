// Package usecase defines the input/output DTOs for use cases.
package usecase

// ProcessOrderInput defines the input for processing an order.
type ProcessOrderInput struct {
	OrderID  string `json:"orderId" validate:"notblank"`
	UserID   string `json:"userId" validate:"notblank"`
	Priority string `json:"priority" validate:"required,oneof=high normal low"`
}

// ProcessOrderOutput defines the enriched output after processing an order.
type ProcessOrderOutput struct {
	OrderID     string  `json:"orderId"`
	Status      string  `json:"status"`
	EstimatedAt string  `json:"estimatedAt"`
	UserEmail   string  `json:"userEmail"`
	UserName    string  `json:"userName"`
	BFF         BFFMeta `json:"bff"`
}

// BFFMeta contains BFF-level processing metadata.
type BFFMeta struct {
	ProcessedBy string `json:"processedBy"`
	ProcessedAt string `json:"processedAt"`
}
