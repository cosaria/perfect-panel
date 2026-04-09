package billing

import "time"

type Order struct {
	ID               int64     `gorm:"primaryKey"`
	ParentOrderID    int64     `gorm:"index;default:0"`
	UserID           int64     `gorm:"index;not null;default:0"`
	OrderNo          string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	Type             uint8     `gorm:"type:tinyint(1);not null;default:1"`
	Quantity         int64     `gorm:"not null;default:1"`
	Price            int64     `gorm:"not null;default:0"`
	Amount           int64     `gorm:"not null;default:0"`
	Discount         int64     `gorm:"not null;default:0"`
	Coupon           string    `gorm:"type:varchar(255);default:''"`
	CouponDiscount   int64     `gorm:"not null;default:0"`
	PaymentGatewayID int64     `gorm:"index;not null;default:0"`
	Method           string    `gorm:"type:varchar(255);not null;default:''"`
	FeeAmount        int64     `gorm:"not null;default:0"`
	TradeNo          string    `gorm:"type:varchar(255);default:''"`
	GiftAmount       int64     `gorm:"not null;default:0"`
	Commission       int64     `gorm:"not null;default:0"`
	Status           uint8     `gorm:"type:tinyint(1);not null;default:1"`
	SubscribeID      int64     `gorm:"index;not null;default:0"`
	SubscribeToken   string    `gorm:"type:varchar(255);default:''"`
	IsNew            bool      `gorm:"not null;default:false"`
	CreatedAt        time.Time `gorm:"<-:create"`
	UpdatedAt        time.Time
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	ID          int64     `gorm:"primaryKey"`
	OrderID     int64     `gorm:"index;not null"`
	SubscribeID int64     `gorm:"index;not null;default:0"`
	Quantity    int64     `gorm:"not null;default:1"`
	UnitPrice   int64     `gorm:"not null;default:0"`
	Amount      int64     `gorm:"not null;default:0"`
	Snapshot    string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"<-:create"`
	UpdatedAt   time.Time
}

func (OrderItem) TableName() string {
	return "order_items"
}

type PaymentGateway struct {
	ID           int64     `gorm:"primaryKey"`
	Name         string    `gorm:"type:varchar(100);not null;default:''"`
	Platform     string    `gorm:"type:varchar(100);not null;default:''"`
	Icon         string    `gorm:"type:varchar(255);default:''"`
	Domain       string    `gorm:"type:varchar(255);default:''"`
	PublicConfig string    `gorm:"type:text;not null"`
	Description  string    `gorm:"type:text"`
	FeeMode      uint      `gorm:"type:tinyint(1);not null;default:0"`
	FeePercent   int64     `gorm:"not null;default:0"`
	FeeAmount    int64     `gorm:"not null;default:0"`
	Enable       *bool     `gorm:"not null;default:false"`
	Token        string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	CreatedAt    time.Time `gorm:"<-:create"`
	UpdatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"`
}

func (PaymentGateway) TableName() string {
	return "payment_gateways"
}

type PaymentGatewaySecret struct {
	ID               int64     `gorm:"primaryKey"`
	PaymentGatewayID int64     `gorm:"not null;uniqueIndex"`
	SecretConfig     string    `gorm:"type:text;not null"`
	CreatedAt        time.Time `gorm:"<-:create"`
	UpdatedAt        time.Time
}

func (PaymentGatewaySecret) TableName() string {
	return "payment_gateway_secrets"
}

type Payment struct {
	ID               int64     `gorm:"primaryKey"`
	OrderID          int64     `gorm:"index;not null;default:0"`
	PaymentGatewayID int64     `gorm:"index;not null;default:0"`
	Amount           int64     `gorm:"not null;default:0"`
	TradeNo          string    `gorm:"type:varchar(255);default:''"`
	Status           uint8     `gorm:"type:tinyint(1);not null;default:0"`
	RawPayload       string    `gorm:"type:text"`
	CreatedAt        time.Time `gorm:"<-:create"`
	UpdatedAt        time.Time
}

func (Payment) TableName() string {
	return "payments"
}

type PaymentCallback struct {
	ID           int64     `gorm:"primaryKey"`
	PaymentID    int64     `gorm:"index;not null;default:0"`
	CallbackType string    `gorm:"type:varchar(100);not null;default:''"`
	CallbackID   string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	RawPayload   string    `gorm:"type:text"`
	Status       string    `gorm:"type:varchar(50);not null;default:''"`
	CreatedAt    time.Time `gorm:"<-:create"`
	UpdatedAt    time.Time
}

func (PaymentCallback) TableName() string {
	return "payment_callbacks"
}

type Refund struct {
	ID        int64     `gorm:"primaryKey"`
	OrderID   int64     `gorm:"index;not null;default:0"`
	Amount    int64     `gorm:"not null;default:0"`
	Status    string    `gorm:"type:varchar(50);not null;default:''"`
	Reason    string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func (Refund) TableName() string {
	return "refunds"
}

type BillingLedger struct {
	ID        int64     `gorm:"primaryKey"`
	OrderID   int64     `gorm:"index;not null;default:0"`
	UserID    int64     `gorm:"index;not null;default:0"`
	Amount    int64     `gorm:"not null;default:0"`
	Kind      string    `gorm:"type:varchar(50);not null;default:''"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func (BillingLedger) TableName() string {
	return "billing_ledgers"
}

type GatewayFilter struct {
	Enable *bool
	Mark   string
	Search string
}

type GatewayRecord struct {
	ID          int64
	Name        string
	Platform    string
	Icon        string
	Domain      string
	Config      string
	Description string
	FeeMode     uint
	FeePercent  int64
	FeeAmount   int64
	Enable      *bool
	Token       string
}

type OrderRecord struct {
	ID             int64
	ParentID       int64
	UserID         int64
	OrderNo        string
	Type           uint8
	Quantity       int64
	Price          int64
	Amount         int64
	Discount       int64
	Coupon         string
	CouponDiscount int64
	PaymentID      int64
	Method         string
	FeeAmount      int64
	TradeNo        string
	GiftAmount     int64
	Commission     int64
	Status         uint8
	SubscribeID    int64
	SubscribeToken string
	IsNew          bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Payment        *GatewayRecord
}
