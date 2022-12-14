package azurermagw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type stringDefaultModifier struct {
    Default string
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m stringDefaultModifier) Description(ctx context.Context) string {
    return fmt.Sprintf("If value is not configured, defaults to %s", m.Default)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m stringDefaultModifier) MarkdownDescription(ctx context.Context) string {
    return fmt.Sprintf("If value is not configured, defaults to `%s`", m.Default)
}

// Modify runs the logic of the plan modifier.
// Access to the configuration, plan, and state is available in `req`, while
// `resp` contains fields for updating the planned value, triggering resource
// replacement, and returning diagnostics.
func (m stringDefaultModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
    // types.String must be the attr.Value produced by the attr.Type in the schema for this attribute
    // for generic plan modifiers, use
    // https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#ConvertValue
    // to convert into a known type.
    var str types.String
    diags := tfsdk.ValueAs(ctx, req.AttributePlan, &str)
    resp.Diagnostics.Append(diags...)
    if diags.HasError() {
        return
    }

    if (!str.Null)&&(str.Value !="") {
		return
    }
	    
    resp.AttributePlan = types.String{Value: m.Default}
}