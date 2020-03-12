package vmWareGo

import (
	"context"
	"encoding/json"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
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

func (vm *Vmware) VmCustomFields(targetVm mo.VirtualMachine, tagsStruct interface{}) error {
	ctx := context.Background()
	customFieldsMap := make(map[int32]string)

	customFieldsMap, err := setCustomFieldsMap(ctx, vm.Client.Client)
	if err != nil {
		return err
	}

	customFields := getCustomFields(targetVm.Summary.CustomValue, customFieldsMap)

	if len(customFields) > 0 {
		err := json.Unmarshal([]byte(customFields.String()), &tagsStruct)
		if err != nil {
			return err
		}
	}

	return nil
}
