package google

import (
	"fmt"
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

const ServiceUsageAssetType string = "serviceusage.googleapis.com/Service"

func resourceConverterServiceUsage() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType: ServiceUsageAssetType,
		Convert:   GetServiceUsageCaiObject,
	}
}

func GetServiceUsageCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	name, err := tpgresource.AssetName(d, config, "//serviceusage.googleapis.com/projects/{{project}}/services/{{service}}")
	if err != nil {
		return []tpgresource.Asset{}, err
	}
	if obj, err := GetServiceUsageApiObject(d, config); err == nil {
		return []tpgresource.Asset{{
			Name: name,
			Type: ServiceUsageAssetType,
			Resource: &tpgresource.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/serviceusage/v1/rest",
				DiscoveryName:        "Service",
				Data:                 obj,
			}},
		}, nil
	}
	return []tpgresource.Asset{}, err
}

func GetServiceUsageApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	parentProjectProp, err := expandServiceUsageParentProject(d.Get("project"), d, config)
	if err != nil {
		return nil, err
	}
	obj["parent"] = parentProjectProp

	serviceNameProp, err := expandServiceUsageServiceName(d.Get("service"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("service"); !tpgresource.IsEmptyValue(reflect.ValueOf(serviceNameProp)) && (ok || !reflect.DeepEqual(v, serviceNameProp)) {
		obj["name"] = serviceNameProp
	}

	obj["state"] = "ENABLED"

	return obj, nil
}

func expandServiceUsageParentProject(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil || v.(string) == "" {
		// It does not try to construct anything from empty.
		return "", nil
	}
	// Ideally we should use project_number, but since that is generated server-side,
	// we substitute project_id.
	return fmt.Sprintf("projects/%s", v.(string)), nil
}

func expandServiceUsageServiceName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
