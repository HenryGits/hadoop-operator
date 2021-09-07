/**
 @author: ZHC
 @date: 2021-09-07 16:47:40
 @description:
**/
package tools

// Helper functions to check and remove string from a slice of strings.
func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
