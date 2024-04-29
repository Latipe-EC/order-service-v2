package order

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

const (
	OI_PENDING  = 0
	OI_PREPARED = 1
	OI_CANCEL   = -1
)

type OrderItem struct {
	OrderType   string
	Id          string    `gorm:"not null;type:char(16);primary_key" json:"item_id"`
	OrderID     string    `gorm:"not null;type:char(16);" json:"order_id"`
	Order       *Order    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductID   string    `gorm:"not null;type:varchar(255)" json:"product_id"`
	ProductName string    `gorm:"not null;type:varchar(255)" json:"product_name"`
	ProdImg     string    `gorm:"not null;type:TEXT" json:"image"`
	RatingID    string    `gorm:"not null;type:varchar(255)" json:"rating_id"`
	OptionID    string    `gorm:"not null;type:varchar(250)" json:"option_id"`
	NameOption  string    `gorm:"not null;type:varchar(250)" json:"name_option"`
	Quantity    int       `gorm:"not null;type:int" json:"quantity"`
	Price       int       `gorm:"not null;type:bigint" json:"price"`
	NetPrice    int       `gorm:"not null;type:bigint" json:"net_price"`
	SubTotal    int       `gorm:"not null;type:bigint" json:"sub_total"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;type:datetime(6)" json:"updated_at"`
	CreatedAt   time.Time `gorm:"autoCreateTime;type:datetime(6)" json:"created_at"`
}

func (o *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	keyGen := strings.ReplaceAll(uuid.NewString(), "-", "")
	o.Id = fmt.Sprintf("%v%v", o.OrderID[:5], keyGen[21:])
	return nil
}

func (OrderItem) TableName() string {
	return "order_items"
}

type OrderStatusLog struct {
	OrderType    string
	Id           int       `gorm:"not null;autoIncrement;primaryKey;type:bigint" json:"id"`
	OrderID      string    `gorm:"not nulltype:char(16)" json:"order_id"`
	Order        *Order    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Message      string    `gorm:"type:longtext" json:"message"`
	StatusChange int       `gorm:"type:int" json:"status_change"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;type:datetime(6)" json:"updated_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime;type:datetime(6)" json:"created_at"`
}

func (OrderStatusLog) TableName() string {
	return "order_status_logs"
}

type Order struct {
	OrderID          string            `gorm:"not null;type:char(16);primaryKey;" json:"order_id"`
	UserId           string            `gorm:"not null;type:varchar(250)" json:"user_id"`
	Username         string            `gorm:"not null;type:varchar(250)" json:"email"`
	Amount           int               `gorm:"not null;type:bigint" json:"amount"`
	ShippingCost     int               `gorm:"not null;type:int" json:"shipping_cost"`
	Vouchers         string            `gorm:"not null;type:varchar(250)" json:"vouchers"`
	PaymentDiscount  int               `gorm:"not null;type:int" json:"payment_discount"`
	StoreDiscount    int               `gorm:"not null;type:int" json:"store_discount"`
	ShippingDiscount int               `gorm:"not null;type:int" json:"shipping_discount"`
	SubTotal         int               `gorm:"not null;type:int" json:"sub_total"`
	Status           int               `gorm:"not null;type:int" json:"status"`
	PaymentMethod    int               `json:"payment_method" gorm:"not null;type:int"`
	StoreId          string            `gorm:"not null;type:varchar(255)" json:"store_id"`
	Thumbnail        string            `gorm:"type:varchar(255)" json:"thumbnail"`
	UpdatedAt        time.Time         `gorm:"autoUpdateTime;type:datetime(6)" json:"updated_at"`
	CreatedAt        time.Time         `gorm:"autoCreateTime;type:datetime(6)" json:"created_at"`
	OrderItem        []*OrderItem      `gorm:"constraint:OnUpdate:CASCADE;polymorphic:Order;" json:"order_items"`
	OrderStatusLog   []*OrderStatusLog `gorm:"constraint:OnUpdate:CASCADE;polymorphic:Order;" json:"order_status"`
	OrderCommission  *OrderCommission  `gorm:"constraint:OnUpdate:CASCADE;polymorphic:Order;" json:"order_commissions"`
	Delivery         *DeliveryOrder    `gorm:"constraint:OnUpdate:CASCADE;polymorphic:Order;" json:"delivery"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderCommission struct {
	OrderType         string
	Id                int       `gorm:"not null;autoIncrement;primaryKey;type:bigint" json:"id"`
	OrderID           string    `gorm:"not null;type:char(16)" json:"order_id"`
	StoreID           string    `gorm:"not null;type:varchar(250)" json:"store_id"`
	DiscountFromStore int       `gorm:"not null;type:bigint" json:"discount_from_store"`
	SystemFee         int       `gorm:"not null;type:bigint" json:"system_fee"`
	Status            int       `gorm:"not null;int" json:"status"`
	AmountReceived    int       `gorm:"not null;type:bigint" json:"amount_received"`
	Order             *Order    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt         time.Time `gorm:"autoCreateTime;type:datetime(6)" json:"created_at"`
}

func (OrderCommission) TableName() string {
	return "orders_commission"
}
