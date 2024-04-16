// SPDX-License-Identifier: Apache-2.0

package router

import (
	"api-gateway/v2/modules/lura/v2/config"
)

func IsValidSequentialEndpoint(_ *config.EndpointConfig) bool {
	// if endpoint.ExtraConfig[proxy.Namespace] == nil {
	// 	return false
	// }

	// proxyCfg := endpoint.ExtraConfig[proxy.Namespace].(map[string]interface{})
	// if proxyCfg["sequential"] == false {
	// 	return false
	// }

	// for i, backend := range endpoint.Backend {
	// 	if backend.Method != http.MethodGet && (i+1) != len(endpoint.Backend) {
	// 		return false
	// 	}
	// }

	return true
}
