package azurermagw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type intDefaultModifier struct {
    Default int64
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m intDefaultModifier) Description(ctx context.Context) string {
    return fmt.Sprintf("If value is not configured, defaults to %v", m.Default)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m intDefaultModifier) MarkdownDescription(ctx context.Context) string {
    return fmt.Sprintf("If value is not configured, defaults to `%v`", m.Default)
}

// Modify runs the logic of the plan modifier.
// Access to the configuration, plan, and state is available in `req`, while
// `resp` contains fields for updating the planned value, triggering resource
// replacement, and returning diagnostics.
func (m intDefaultModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
    // types.int must be the attr.Value produced by the attr.Type in the schema for this attribute
    // for generic plan modifiers, use
    // https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#ConvertValue
    // to convert into a known type.
    var str types.Int64
    diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &str)
    resp.Diagnostics.Append(diags...)
    if diags.HasError() {
        return
    }
    if (str.Value!=0){
		return
    }
	resp.AttributePlan = types.Int64{Value: m.Default}
}