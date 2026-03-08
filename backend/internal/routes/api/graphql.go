package api

import (
	"backend/internal/config"
	"backend/internal/graph"
	"backend/internal/services"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

// GraphQLHandler handles GraphQL requests
// @Summary Unified GraphQL API
// @Description Execute complex queries and mutations for properties and background jobs.
// @Description
// @Description This endpoint provides a flexible way to fetch exactly the data you need.
// @Description It supports filtering, pagination, and nested data structures.
// @Description
// @Description ### Authentication
// @Description This endpoint uses Optional JWT Authentication.
// @Description - Public users can query property listings.
// @Description - Authenticated users can access additional details.
// @Description - Admin users (identified by JWT) can perform mutations like creating, updating, or deleting properties.
// @Description
// @Description ### Example Property Query
// @Description ```graphql
// @Description query {
// @Description   properties(limit: 10, minPrice: 200000, townCity: "London") {
// @Description     items {
// @Description       id
// @Description       price
// @Description       street
// @Description       postcode
// @Description     }
// @Description     total
// @Description   }
// @Description }
// @Description ```
// @Description
// @Description ### Example Job Query (Admin only)
// @Description ```graphql
// @Description query {
// @Description   jobs(limit: 5) {
// @Description     items {
// @Description       id
// @Description       status
// @Description       progress
// @Description       taskType
// @Description     }
// @Description     total
// @Description   }
// @Description }
// @Description ```
// @Description
// @Description ### Example Mutations (Admin only)
// @Description
// @Description **Create Property:**
// @Description ```graphql
// @Description mutation {
// @Description   createProperty(input: {
// @Description     price: 250000,
// @Description     dateOfTransfer: "2024-03-07T00:00:00Z",
// @Description     postcode: "SW1A 1AA",
// @Description     propertyType: "D",
// @Description     oldNew: "N",
// @Description     duration: "F",
// @Description     paon: "10",
// @Description     saon: "",
// @Description     street: "Downing Street",
// @Description     locality: "Westminster",
// @Description     townCity: "London",
// @Description     district: "Greater London",
// @Description     county: "London",
// @Description     ppdCategoryType: "A",
// @Description     recordStatus: "A"
// @Description   }) { id price street }
// @Description }
// @Description ```
// @Description
// @Description **Update Property:**
// @Description ```graphql
// @Description mutation {
// @Description   updateProperty(id: "UUID-HERE", input: {
// @Description     price: 260000,
// @Description     # ... other fields
// @Description   }) { id price }
// @Description }
// @Description ```
// @Description
// @Description **Delete Property:**
// @Description ```graphql
// @Description mutation {
// @Description   deleteProperty(id: "UUID-HERE")
// @Description }
// @Description ```
// @Tags GraphQL
// @Accept json
// @Produce json
// @Param request body GraphQLRequest true "GraphQL Operation (Query or Mutation)"
// @Success 200 {object} GraphQLResponse
// @Failure 400 {object} ErrorResponse "Invalid GraphQL request"
// @Failure 401 {object} ErrorResponse "Unauthorized (if authentication is required for specific operation)"
// @Router /api/query [post]
func GraphQLHandler(svcs *services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := c.MustGet("config").(*config.Config)
		h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
			Config:          cfg,
			PropertyService: svcs.Property,
			JobService:      svcs.Job,
		}}))

		h.ServeHTTP(c.Writer, c.Request)
	}
}

// PlaygroundHandler serves the GraphQL Playground
func PlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/api/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
