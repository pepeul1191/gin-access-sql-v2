// pkg/utils/templates.go
package utils

func FirstFlashOrEmpty(flashes []interface{}) string {
	if len(flashes) > 0 {
		return flashes[0].(string)
	}
	return ""
}
