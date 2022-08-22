/**
 * @Author raven
 * @Description
 * @Date 2022/8/18
 **/
package internal

func StringSlice2InterfaceSlice(stringSlice []string) []interface{} {
	interfaceSlice := make([]interface{}, 0, len(stringSlice))
	for _, item := range stringSlice {
		interfaceSlice = append(interfaceSlice, item)
	}
	return interfaceSlice
}
