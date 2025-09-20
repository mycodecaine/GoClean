package handlers

import (
	"context"
	"goclean/internal/application/commands"
	"goclean/internal/application/dto"
	"goclean/internal/application/queries"
	"goclean/internal/infrastructure/auth"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ProductHandler handles product-related HTTP requests
type ProductHandler struct {
	productCommandHandler *commands.ProductCommandHandler
	productQueryHandler   *queries.ProductQueryHandler
}

// NewProductHandler creates a new product handler
func NewProductHandler(
	productCommandHandler *commands.ProductCommandHandler,
	productQueryHandler *queries.ProductQueryHandler,
) *ProductHandler {
	return &ProductHandler{
		productCommandHandler: productCommandHandler,
		productQueryHandler:   productQueryHandler,
	}
}

// CreateProduct creates a new product
// @Summary Create a new product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.CreateProductRequest true "Product data"
// @Success 201 {object} dto.ProductAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 401 {object} dto.ErrorAPIResponse
// @Failure 500 {object} dto.ErrorAPIResponse
// @Router /api/v1/products [post]
// @Security BearerAuth
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var req dto.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse[interface{}]{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Get current user from context
	claims, ok := c.Get("user_claims").(*auth.UserClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.APIResponse[interface{}]{
			Success: false,
			Error:   "User not authenticated",
		})
	}

	createdBy, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse[interface{}]{
			Success: false,
			Error:   "Invalid user ID",
		})
	}

	// Convert to command
	cmd := commands.CreateProductCommand{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SKU:         req.SKU,
		Category:    req.Category,
		CreatedBy:   createdBy,
	}

	// Execute command
	if err := h.productCommandHandler.Handle(context.Background(), cmd); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse[interface{}]{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.APIResponse[interface{}]{
		Success: true,
		Message: "Product created successfully",
	})
}

// GetProduct retrieves a product by ID
// @Summary Get product by ID
// @Description Get product information by product ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductAPIResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Failure 404 {object} dto.ErrorAPIResponse
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse[interface{}]{
			Success: false,
			Error:   "Invalid product ID",
		})
	}

	query := queries.GetProductByIDQuery{ID: id}
	result, err := h.productQueryHandler.Handle(context.Background(), query)
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.APIResponse[interface{}]{
			Success: false,
			Error:   "Product not found",
		})
	}

	// Convert to DTO
	productDTO := &dto.ProductDTO{
		ID:          result.Product.ID,
		Name:        result.Product.Name,
		Description: result.Product.Description,
		Price:       result.Product.Price,
		SKU:         result.Product.SKU,
		Category:    result.Product.Category,
		IsActive:    result.Product.IsActive,
		CreatedBy:   result.Product.CreatedBy,
		CreatedAt:   result.Product.CreatedAt,
		UpdatedAt:   result.Product.UpdatedAt,
	}

	return c.JSON(http.StatusOK, dto.APIResponse[*dto.ProductDTO]{
		Success: true,
		Data:    productDTO,
	})
}

// ListProducts retrieves products with pagination
// @Summary List products
// @Description Get list of products with pagination
// @Tags products
// @Produce json
// @Param offset query int false "Offset" default(0)
// @Param limit query int false "Limit" default(10)
// @Param category query string false "Category filter"
// @Param search query string false "Search query"
// @Success 200 {object} dto.ProductsListResponse
// @Failure 400 {object} dto.ErrorAPIResponse
// @Router /api/v1/products [get]
func (h *ProductHandler) ListProducts(c echo.Context) error {
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}

	category := c.QueryParam("category")
	search := c.QueryParam("search")

	var result *queries.ProductsResult
	var err error

	if search != "" {
		query := queries.SearchProductsQuery{
			Query:  search,
			Offset: offset,
			Limit:  limit,
		}
		result, err = h.productQueryHandler.HandleSearch(context.Background(), query)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.APIResponse[interface{}]{
				Success: false,
				Error:   err.Error(),
			})
		}
	} else if category != "" {
		query := queries.ListProductsByCategoryQuery{
			Category: category,
			Offset:   offset,
			Limit:    limit,
		}
		result, err = h.productQueryHandler.HandleByCategory(context.Background(), query)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.APIResponse[interface{}]{
				Success: false,
				Error:   err.Error(),
			})
		}
	} else {
		query := queries.ListProductsQuery{
			Offset: offset,
			Limit:  limit,
		}
		result, err = h.productQueryHandler.HandleList(context.Background(), query)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.APIResponse[interface{}]{
				Success: false,
				Error:   err.Error(),
			})
		}
	}

	// Convert to DTOs
	productDTOs := make([]dto.ProductDTO, len(result.Products))
	for i, product := range result.Products {
		productDTOs[i] = dto.ProductDTO{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			SKU:         product.SKU,
			Category:    product.Category,
			IsActive:    product.IsActive,
			CreatedBy:   product.CreatedBy,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	return c.JSON(http.StatusOK, dto.PaginatedResponse[[]dto.ProductDTO]{
		APIResponse: dto.APIResponse[[]dto.ProductDTO]{
			Success: true,
			Data:    productDTOs,
		},
		Pagination: dto.PaginationInfo{
			Offset: offset,
			Limit:  limit,
			Total:  result.Total,
		},
	})
}
