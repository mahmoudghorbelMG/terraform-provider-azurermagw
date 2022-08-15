package main

import (
	"context"
	//"fmt"
	//"log"
	//"reflect"
	"terraform-provider-azurermagw/azurermagw"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)


func main() {
	/*type test struct {
		A bool
		B bool
		C bool
	}
	v := test{
		A: false,
		B: false,
	}
	
	var h *azurermagw.Http_listener
	h = &azurermagw.Http_listener{}
	if h != nil {
		fmt.Println("Field nzerazerazer ", "jjj"+h.Name.Value+" jjj")
	}*/

	//log.Printf("Field %s rrrrrrrrrrrrrrrrrrrr ", hasField(v,"Z"))
	/*metaValue := reflect.ValueOf(v).Elem()

	for _, name := range []string{"A", "B", "Z"} {
		field := metaValue.FieldByName(name)
		if field == (reflect.Value{}) {
			log.Printf("Field %s not exist in struct", name)
		}
	}*/
	
	/*probe_string := "/subscriptions/__AZURE_SUBSCRIPTION_ID__/resourceGroups/__rg_name__/providers/Microsoft.Network/applicationGateways/__agw_name__/probes/"
	str := probe_string + "mahmoud-probe"
	fmt.Println("\n-- probe_string: ",str)
	splitted_list := strings.Split(str,"/")
	probe_string1 := splitted_list[len(splitted_list)-1]
	fmt.Println("\n-- probe_string1: ",probe_string1)*/
	

	/*for i := 0; i < 3; i++ {
		fmt.Print("##################################")
	}*/
	tfsdk.Serve(context.Background(), azurermagw.New, tfsdk.ServeOpts{
		Name: "azurermagw",	})
}
