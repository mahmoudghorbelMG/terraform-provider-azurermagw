package azurermagw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type stringArrayDefaultModifier struct {
    Default []string
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m stringArrayDefaultModifier) Description(ctx context.Context) string {
    return fmt.Sprintf("If value is not configured, defaults to %s", m.Default)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m stringArrayDefaultModifier) MarkdownDescription(ctx context.Context) string {
    return fmt.Sprintf("If value is not configured, defaults to `%s`", m.Default)
}

// Modify runs the logic of the plan modifier.
// Access to the configuration, plan, and state is available in `req`, while
// `resp` contains fields for updating the planned value, triggering resource
// replacement, and returning diagnostics.
func (m stringArrayDefaultModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
    // types.String must be the attr.Value produced by the attr.Type in the schema for this attribute
    // for generic plan modifiers, use
    // https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#ConvertValue
    // to convert into a known type.
    var str []types.String
    diags := tfsdk.ValueAs(ctx, req.AttributePlan, &str)
    resp.Diagnostics.Append(diags...)
    if diags.HasError() {
        return
    }

    if (len(str)!= 0) {
		return
    }
	var array []types.String
    for i := 0; i < len(m.Default); i++ {
        array[i] = types.String{Value: m.Default[i]}
    }
    //resp.AttributePlan = types.List(array)
}