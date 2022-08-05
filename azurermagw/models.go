package azurermagw

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WebappBinding -
type WebappBinding struct {
	Name                 types.String         `tfsdk:"name"`
	Agw_name             types.String         `tfsdk:"agw_name"`
	Agw_rg               types.String         `tfsdk:"agw_rg"`
	Backend_address_pool Backend_address_pool `tfsdk:"backend_address_pool"`
	//Backend_http_settings   Backend_http_settings	`tfsdk:"backend_http_settings"`
}
type Backend_address_pool struct {
	Name         types.String   `tfsdk:"name"`
	Id           types.String   `tfsdk:"id"`
	Fqdns        []types.String `tfsdk:"fqdns"`
	Ip_addresses []types.String `tfsdk:"ip_addresses"`
}
 /*
type Backend_http_settings struct {
	Name		types.String	`tfsdk:"name"`
	Protocol	types.String	`tfsdk:"protocol"`
}*/




// Order -
type Order struct {
	ID          types.String `tfsdk:"id"`
	Items       []OrderItem  `tfsdk:"items"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// OrderItem -
type OrderItem struct {
	Coffee   Coffee `tfsdk:"coffee"`
	Quantity int    `tfsdk:"quantity"`
}

// Coffee -
// This Coffee struct is for Order.Items[].Coffee which does not have an
// ingredients field in the schema defined by plugin framework. Since the
// resource schema must match the struct exactly (extra field will return an
// error). This struct has Ingredients commented out.
type Coffee struct {
	ID          int          `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Teaser      types.String `tfsdk:"teaser"`
	Description types.String `tfsdk:"description"`
	Price       types.Number `tfsdk:"price"`
	Image       types.String `tfsdk:"image"`
	// Ingredients []Ingredient   `tfsdk:"ingredients"`
}

// Ingredient -
type Ingredient struct {
	ID       int    `tfsdk:"ingredient_id"`
	Name     string `tfsdk:"name"`
	Quantity int    `tfsdk:"quantity"`
	Unit     string `tfsdk:"unit"`
}

//
// Coffee Data Source specific structs
//

// CoffeeIngredients
type CoffeeIngredients struct {
	ID          int            `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	Teaser      types.String   `tfsdk:"teaser"`
	Description types.String   `tfsdk:"description"`
	Price       types.Number   `tfsdk:"price"`
	Image       types.String   `tfsdk:"image"`
	Ingredient  []IngredientID `tfsdk:"ingredients"`
}

// Ingredient -
type IngredientID struct {
	ID int `tfsdk:"id"`
}