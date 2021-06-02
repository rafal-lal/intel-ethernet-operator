package utils

import (
	"testing"

	"github.com/golang/protobuf/ptypes"
	flowapi "github.com/otcshare/intel-ethernet-operator/pkg/rpc/v1/flow"
)

func TestGetFlowActionAny(t *testing.T) {

	actionData := []struct {
		Type string
		Conf []byte
	}{
		{
			Type: "RTE_FLOW_ACTION_TYPE_VF",
			Conf: []byte(`
				 { "id": 1 }
				`),
		},
		{
			Type: "RTE_FLOW_ACTION_TYPE_END",
			Conf: []byte(`{}`),
		},
	}

	for _, item := range actionData {
		any, err := GetFlowActionAny(item.Type, item.Conf)
		if err != nil {
			t.Errorf("%v", err)
		}
		_ = any
	}
}

func TestGetFlowItemAny(t *testing.T) {

	itemData := []struct {
		Type string
		Spec []byte
	}{
		{
			Type: "RTE_FLOW_ITEM_TYPE_IPV4",
			Spec: []byte(`{
				"hdr": {
					"dst_addr": "1.1.1.1"
				}
			}`),
		},
		{
			Type: "RTE_FLOW_ITEM_TYPE_END",
			Spec: []byte(`{}`),
		},
	}

	for _, item := range itemData {
		any, err := GetFlowItemAny(item.Type, item.Spec)
		if err != nil {
			t.Errorf("%v", err)
		}
		_ = any
		//fmt.Printf("marshalled any object: %+v\n", any)
	}
}

func TestItemAnyObjects(t *testing.T) {
	item := struct {
		Type string
		Spec []byte
	}{
		Type: "RTE_FLOW_ITEM_TYPE_IPV4",
		Spec: []byte(`{
						"hdr": {
							"dst_addr": "1.1.1.1"
						}
				}`),
	}

	any, err := GetFlowItemAny(item.Type, item.Spec)
	if err != nil {
		t.Errorf("%v", err)
	}
	//fmt.Printf("any.Value: %s\n", string(any.Value))
	ipv4 := &flowapi.RteFlowItemIpv4{}

	if err := ptypes.UnmarshalAny(any, ipv4); err != nil {
		t.Errorf("%v", err)
	}

	//fmt.Printf("unmarshalled any object: %+v\n", ipv4)
}
