package vmWareGo

import (
	"context"
	"encoding/json"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"strconv"
	"strings"
)

func getCustomFields(customFields []types.BaseCustomFieldValue, customFieldsMap map[int32]string) MapStr {
	outputFields := MapStr{}
	for _, v := range customFields {
		customFieldString := v.(*types.CustomFieldStringValue)
		key, ok := customFieldsMap[v.GetCustomFieldValue().Key]
		if ok {
			// If key has '.', is replaced with '_' to be compatible with ES2.x.
			fmtKey := strings.Replace(key, ".", "_", -1)
			outputFields.Put(fmtKey, customFieldString.Value)
		}
	}

	return outputFields
}

func setCustomFieldsMap(ctx context.Context, c *vim25.Client) (map[int32]string, error) {
	customFieldsMap := make(map[int32]string)

	customFieldsManager, err := object.GetCustomFieldsManager(c)

	if err != nil {
		return nil, err
	} else {
		field, err := customFieldsManager.Field(ctx)
		if err != nil {
			return nil, err
		}

		for _, def := range field {
			customFieldsMap[def.Key] = def.Name
		}
	}

	return customFieldsMap, nil
}

func getTmplTag(vm mo.VirtualMachine, c *govmomi.Client) (bool, error) {
	ctx := context.Background()
	customFieldsMap := make(map[int32]string)
	tmpl := false

	customFieldsMap, err := setCustomFieldsMap(ctx, c.Client)
	if err != nil {
		return false, err
	}

	var tags Tags

	customFields := getCustomFields(vm.Summary.CustomValue, customFieldsMap)

	if len(customFields) > 0 {
		err := json.Unmarshal([]byte(customFields.String()), &tags)
		if err != nil {
			return false, err
		}
		tmpl, _ = strconv.ParseBool(tags.Template)
	}

	return tmpl, nil
}