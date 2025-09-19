# Aggregate Root and Soft Delete Patterns in Go Clean Architecture

## Overview

This implementation provides Enterprise-grade Aggregate Root and Soft Delete patterns similar to what you'd find in C# Entity Framework, but adapted for Go's idioms and your Clean Architecture setup.

## 🏗️ Key Components

### 1. **Aggregate Root Pattern**

```go
// Base aggregate root with domain events
type AggregateRoot struct {
    domainEvents []DomainEvent
}

// User as an aggregate root
type User struct {
    BaseEntity                    // Common fields + soft delete
    AggregateRoot                 // Domain events capability
    Email     string
    Username  string
    // ... other fields
}
```

**Benefits:**
- ✅ Encapsulates business logic and invariants
- ✅ Manages domain events for eventual consistency
- ✅ Controls transaction boundaries
- ✅ Ensures data consistency within the aggregate

### 2. **Soft Delete Pattern**

```go
// Base entity with soft delete capability
type BaseEntity struct {
    ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
    CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"` // Soft delete field
}

// Soft delete methods
func (be *BaseEntity) SoftDelete() {
    now := time.Now()
    be.DeletedAt = &now
}

func (be *BaseEntity) IsDeleted() bool {
    return be.DeletedAt != nil
}

func (be *BaseEntity) Restore() {
    be.DeletedAt = nil
}
```

**Benefits:**
- ✅ Data preservation for audit trails
- ✅ Easy restoration of deleted records
- ✅ Compliance with data retention policies
- ✅ Safe deletion without losing referential integrity

## 🔧 Implementation Details

### Domain Events System

```go
// Domain events are raised when business operations occur
func (u *User) Delete() {
    u.SoftDelete()
    u.IsActive = false
    
    // Raise domain event
    u.AddDomainEvent(UserDeletedEvent{
        UserID:     u.ID,
        OccurredAt: time.Now(),
    })
}
```

### Repository Pattern with Soft Delete

```go
type UserRepository interface {
    // Standard operations (exclude soft deleted)
    GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
    List(ctx context.Context, offset, limit int) ([]*entities.User, error)
    
    // Soft delete operations
    SoftDelete(ctx context.Context, id uuid.UUID) error
    Restore(ctx context.Context, id uuid.UUID) error
    
    // Include deleted operations
    GetByIDIncludeDeleted(ctx context.Context, id uuid.UUID) (*entities.User, error)
    ListIncludeDeleted(ctx context.Context, offset, limit int) ([]*entities.User, error)
    ListDeleted(ctx context.Context, offset, limit int) ([]*entities.User, error)
}
```

### GORM Implementation

```go
// Exclude soft deleted (default GORM behavior)
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
    var user entities.User
    err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
    return &user, err
}

// Include soft deleted
func (r *userRepository) GetByIDIncludeDeleted(ctx context.Context, id uuid.UUID) (*entities.User, error) {
    var user entities.User
    err := r.db.WithContext(ctx).Unscoped().First(&user, "id = ?", id).Error
    return &user, err
}

// Soft delete
func (r *userRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&entities.User{}, "id = ?", id).Error
}

// Restore
func (r *userRepository) Restore(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Model(&entities.User{}).Unscoped().
        Where("id = ? AND deleted_at IS NOT NULL", id).
        Update("deleted_at", nil).Error
}
```

## 🚀 Usage Examples

### 1. **Creating Aggregates with Domain Events**

```go
// Create user aggregate
user := entities.NewUser("john@example.com", "john_doe", "John", "Doe")

// Save to repository
err := userRepo.Create(ctx, user)

// Dispatch domain events
err = eventDispatcher.DispatchEvents(ctx, &user.AggregateRoot)
```

### 2. **Soft Delete Operations**

```go
// Soft delete through aggregate service
err := userAggregateService.SoftDeleteUser(ctx, userID)

// Restore soft deleted user
err := userAggregateService.RestoreUser(ctx, userID)

// Query deleted users
deletedUsers, err := userRepo.ListDeleted(ctx, 0, 10)
```

### 3. **REST API Endpoints**

```bash
# Soft delete a user
DELETE /api/users/{id}/soft-delete

# Restore a user
POST /api/users/{id}/restore

# Get deleted users
GET /api/users/deleted?offset=0&limit=10

# Get user including deleted
GET /api/users/{id}/with-deleted
```

### 4. **Domain Event Handling**

```go
// Event handler for user creation
type UserCreatedEventHandler struct {
    logger *logger.Logger
}

func (h *UserCreatedEventHandler) Handle(ctx context.Context, event entities.DomainEvent) error {
    userEvent := event.(entities.UserCreatedEvent)
    
    // Send welcome email
    // Create default settings
    // Update analytics
    
    return nil
}
```

## 🎯 Business Benefits

### **Aggregate Root Pattern**
1. **Data Consistency**: Ensures business invariants are maintained
2. **Transaction Boundaries**: Clear boundaries for database transactions
3. **Domain Events**: Enables eventual consistency and integration patterns
4. **Encapsulation**: Business logic stays within the domain layer

### **Soft Delete Pattern**
1. **Data Recovery**: Easy restoration of accidentally deleted data
2. **Audit Trail**: Maintains complete history of all operations
3. **Compliance**: Meets regulatory requirements for data retention
4. **Referential Integrity**: Preserves relationships even after deletion
5. **Performance**: Query optimization for active vs. deleted data

## 🔍 Advanced Features

### **Aggregate Invariants**
```go
func (o *Order) AddItem(productID uuid.UUID, quantity int, price float64) error {
    // Business rule: Maximum 10 items per order
    if len(o.Items) >= 10 {
        return errors.New("maximum items limit exceeded")
    }
    
    item := NewOrderItem(o.ID, productID, quantity, price)
    o.Items = append(o.Items, *item)
    o.RecalculateTotal()
    
    return nil
}
```

### **Cascade Soft Delete**
```go
func (u *User) Delete() {
    u.SoftDelete()
    u.IsActive = false
    
    // Cascade soft delete to related entities
    if u.Profile != nil {
        u.Profile.SoftDelete()
    }
    
    u.AddDomainEvent(UserDeletedEvent{...})
}
```

### **Query Scopes**
```go
// Active users only (default)
users := userRepo.List(ctx, 0, 10)

// All users including deleted
allUsers := userRepo.ListIncludeDeleted(ctx, 0, 10)

// Only deleted users
deletedUsers := userRepo.ListDeleted(ctx, 0, 10)
```

## 🛡️ Best Practices

1. **Always use aggregate methods** for business operations
2. **Handle domain events** for cross-aggregate communication
3. **Validate business rules** before state changes
4. **Use soft delete by default** for user-facing entities
5. **Provide restoration capabilities** for soft deleted data
6. **Implement proper authorization** for delete/restore operations
7. **Monitor deleted data** for cleanup policies
8. **Use transactions** for complex aggregate operations

## 🔄 Migration from Existing Code

To adopt these patterns in your existing codebase:

1. **Add BaseEntity** to existing entities
2. **Update repositories** to support soft delete methods
3. **Migrate existing delete operations** to soft delete
4. **Add domain events** to critical business operations
5. **Create aggregate services** for complex business logic
6. **Update API endpoints** to support soft delete operations

This implementation provides enterprise-grade patterns that are production-ready and maintainable while following Go idioms and Clean Architecture principles.