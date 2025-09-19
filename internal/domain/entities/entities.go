package entities

import (
	"time"

	"github.com/google/uuid"
)

// AggregateRoot represents the base aggregate root with domain events
type AggregateRoot struct {
	domainEvents []DomainEvent
}

// DomainEvent represents a domain event interface
type DomainEvent interface {
	OccurredOn() time.Time
	EventType() string
}

// AddDomainEvent adds a domain event to the aggregate
func (ar *AggregateRoot) AddDomainEvent(event DomainEvent) {
	ar.domainEvents = append(ar.domainEvents, event)
}

// DomainEvents returns all domain events
func (ar *AggregateRoot) DomainEvents() []DomainEvent {
	return ar.domainEvents
}

// ClearDomainEvents clears all domain events
func (ar *AggregateRoot) ClearDomainEvents() {
	ar.domainEvents = nil
}

// BaseEntity represents the base entity with common fields
type BaseEntity struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"` // Soft delete field
}

// IsDeleted checks if the entity is soft deleted
func (be *BaseEntity) IsDeleted() bool {
	return be.DeletedAt != nil
}

// SoftDelete marks the entity as deleted
func (be *BaseEntity) SoftDelete() {
	now := time.Now()
	be.DeletedAt = &now
}

// Restore restores a soft deleted entity
func (be *BaseEntity) Restore() {
	be.DeletedAt = nil
}

// User represents the aggregate root for user domain
type User struct {
	BaseEntity             // Embedded base entity with soft delete
	AggregateRoot          // Embedded aggregate root for domain events
	Email         string   `json:"email" gorm:"uniqueIndex;not null"`
	Username      string   `json:"username" gorm:"uniqueIndex;not null"`
	FirstName     string   `json:"first_name" gorm:"not null"`
	LastName      string   `json:"last_name" gorm:"not null"`
	IsActive      bool     `json:"is_active" gorm:"default:true"`
	Profile       *Profile `json:"profile,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// UserCreatedEvent represents a user created domain event
type UserCreatedEvent struct {
	UserID     uuid.UUID `json:"user_id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	OccurredAt time.Time `json:"occurred_at"`
}

// OccurredOn returns when the event occurred
func (e UserCreatedEvent) OccurredOn() time.Time {
	return e.OccurredAt
}

// EventType returns the event type
func (e UserCreatedEvent) EventType() string {
	return "UserCreated"
}

// UserDeletedEvent represents a user deleted domain event
type UserDeletedEvent struct {
	UserID     uuid.UUID `json:"user_id"`
	OccurredAt time.Time `json:"occurred_at"`
}

// OccurredOn returns when the event occurred
func (e UserDeletedEvent) OccurredOn() time.Time {
	return e.OccurredAt
}

// EventType returns the event type
func (e UserDeletedEvent) EventType() string {
	return "UserDeleted"
}

// NewUser creates a new user aggregate
func NewUser(email, username, firstName, lastName string) *User {
	user := &User{
		BaseEntity: BaseEntity{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email:     email,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true,
	}

	// Add domain event
	user.AddDomainEvent(UserCreatedEvent{
		UserID:     user.ID,
		Email:      user.Email,
		Username:   user.Username,
		OccurredAt: time.Now(),
	})

	return user
}

// Delete soft deletes the user and raises domain event
func (u *User) Delete() {
	u.SoftDelete()
	u.IsActive = false

	// Add domain event
	u.AddDomainEvent(UserDeletedEvent{
		UserID:     u.ID,
		OccurredAt: time.Now(),
	})
}

// TableName returns the table name for GORM
func (u *User) TableName() string {
	return "users"
}

// Profile represents user profile information (child entity of User aggregate)
type Profile struct {
	BaseEntity             // Embedded base entity with soft delete
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	Bio         string     `json:"bio"`
	Avatar      string     `json:"avatar"`
	DateOfBirth *time.Time `json:"date_of_birth"`
}

// NewProfile creates a new profile
func NewProfile(userID uuid.UUID, bio, avatar string, dateOfBirth *time.Time) *Profile {
	return &Profile{
		BaseEntity: BaseEntity{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:      userID,
		Bio:         bio,
		Avatar:      avatar,
		DateOfBirth: dateOfBirth,
	}
}

// TableName returns the table name for GORM
func (p *Profile) TableName() string {
	return "profiles"
}

// Product represents a product aggregate root
type Product struct {
	BaseEntity              // Embedded base entity with soft delete
	AggregateRoot           // Embedded aggregate root for domain events
	Name          string    `json:"name" gorm:"not null;index"`
	Description   string    `json:"description"`
	Price         float64   `json:"price" gorm:"not null"`
	SKU           string    `json:"sku" gorm:"uniqueIndex;not null"`
	Category      string    `json:"category" gorm:"index"`
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	CreatedBy     uuid.UUID `json:"created_by" gorm:"type:uuid;not null"`
}

// ProductCreatedEvent represents a product created domain event
type ProductCreatedEvent struct {
	ProductID  uuid.UUID `json:"product_id"`
	Name       string    `json:"name"`
	SKU        string    `json:"sku"`
	CreatedBy  uuid.UUID `json:"created_by"`
	OccurredAt time.Time `json:"occurred_at"`
}

// OccurredOn returns when the event occurred
func (e ProductCreatedEvent) OccurredOn() time.Time {
	return e.OccurredAt
}

// EventType returns the event type
func (e ProductCreatedEvent) EventType() string {
	return "ProductCreated"
}

// ProductDeletedEvent represents a product deleted domain event
type ProductDeletedEvent struct {
	ProductID  uuid.UUID `json:"product_id"`
	OccurredAt time.Time `json:"occurred_at"`
}

// OccurredOn returns when the event occurred
func (e ProductDeletedEvent) OccurredOn() time.Time {
	return e.OccurredAt
}

// EventType returns the event type
func (e ProductDeletedEvent) EventType() string {
	return "ProductDeleted"
}

// NewProduct creates a new product aggregate
func NewProduct(name, description, sku, category string, price float64, createdBy uuid.UUID) *Product {
	product := &Product{
		BaseEntity: BaseEntity{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        name,
		Description: description,
		Price:       price,
		SKU:         sku,
		Category:    category,
		IsActive:    true,
		CreatedBy:   createdBy,
	}

	// Add domain event
	product.AddDomainEvent(ProductCreatedEvent{
		ProductID:  product.ID,
		Name:       product.Name,
		SKU:        product.SKU,
		CreatedBy:  product.CreatedBy,
		OccurredAt: time.Now(),
	})

	return product
}

// Delete soft deletes the product and raises domain event
func (p *Product) Delete() {
	p.SoftDelete()
	p.IsActive = false

	// Add domain event
	p.AddDomainEvent(ProductDeletedEvent{
		ProductID:  p.ID,
		OccurredAt: time.Now(),
	})
}

// TableName returns the table name for GORM
func (p *Product) TableName() string {
	return "products"
}

// Order represents an order aggregate root
type Order struct {
	BaseEntity                // Embedded base entity with soft delete
	AggregateRoot             // Embedded aggregate root for domain events
	UserID        uuid.UUID   `json:"user_id" gorm:"type:uuid;not null;index"`
	Status        OrderStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	TotalPrice    float64     `json:"total_price" gorm:"not null"`
	Items         []OrderItem `json:"items" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

// OrderCreatedEvent represents an order created domain event
type OrderCreatedEvent struct {
	OrderID    uuid.UUID `json:"order_id"`
	UserID     uuid.UUID `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	ItemCount  int       `json:"item_count"`
	OccurredAt time.Time `json:"occurred_at"`
}

// OccurredOn returns when the event occurred
func (e OrderCreatedEvent) OccurredOn() time.Time {
	return e.OccurredAt
}

// EventType returns the event type
func (e OrderCreatedEvent) EventType() string {
	return "OrderCreated"
}

// OrderCancelledEvent represents an order cancelled domain event
type OrderCancelledEvent struct {
	OrderID    uuid.UUID `json:"order_id"`
	OccurredAt time.Time `json:"occurred_at"`
}

// OccurredOn returns when the event occurred
func (e OrderCancelledEvent) OccurredOn() time.Time {
	return e.OccurredAt
}

// EventType returns the event type
func (e OrderCancelledEvent) EventType() string {
	return "OrderCancelled"
}

// NewOrder creates a new order aggregate
func NewOrder(userID uuid.UUID, items []OrderItem) *Order {
	var totalPrice float64
	for _, item := range items {
		totalPrice += item.Price * float64(item.Quantity)
	}

	order := &Order{
		BaseEntity: BaseEntity{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:     userID,
		Status:     OrderStatusPending,
		TotalPrice: totalPrice,
		Items:      items,
	}

	// Set order ID for all items
	for i := range order.Items {
		order.Items[i].OrderID = order.ID
		if order.Items[i].ID == uuid.Nil {
			order.Items[i].BaseEntity.ID = uuid.New()
			order.Items[i].BaseEntity.CreatedAt = time.Now()
			order.Items[i].BaseEntity.UpdatedAt = time.Now()
		}
	}

	// Add domain event
	order.AddDomainEvent(OrderCreatedEvent{
		OrderID:    order.ID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
		ItemCount:  len(order.Items),
		OccurredAt: time.Now(),
	})

	return order
}

// Cancel cancels the order and raises domain event
func (o *Order) Cancel() {
	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()

	// Add domain event
	o.AddDomainEvent(OrderCancelledEvent{
		OrderID:    o.ID,
		OccurredAt: time.Now(),
	})
}

// TableName returns the table name for GORM
func (o *Order) TableName() string {
	return "orders"
}

// OrderItem represents an order item entity (child entity of Order aggregate)
type OrderItem struct {
	BaseEntity           // Embedded base entity with soft delete
	OrderID    uuid.UUID `json:"order_id" gorm:"type:uuid;not null;index"`
	ProductID  uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	Quantity   int       `json:"quantity" gorm:"not null"`
	Price      float64   `json:"price" gorm:"not null"`
}

// NewOrderItem creates a new order item
func NewOrderItem(orderID, productID uuid.UUID, quantity int, price float64) *OrderItem {
	return &OrderItem{
		BaseEntity: BaseEntity{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  quantity,
		Price:     price,
	}
}

// TableName returns the table name for GORM
func (oi *OrderItem) TableName() string {
	return "order_items"
}

// OrderStatus represents order status value object
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// IsValid checks if the order status is valid
func (os OrderStatus) IsValid() bool {
	switch os {
	case OrderStatusPending, OrderStatusConfirmed, OrderStatusShipped, OrderStatusDelivered, OrderStatusCancelled:
		return true
	default:
		return false
	}
}
