
package decoding

func GetString(js map[string]interface{}, key string) string {
    str, ok := js[key].(string)
    if ok {
        return str
    }
    return ""
}

func GetNumber(js map[string]interface{}, key string) float64 {
    num, ok := js[key].(float64)
    if ok {
        return num
    }
    return float64(0.0)
}

func GetInteger(js map[string]interface{}, key string) int64 {
    return int64(GetNumber(js, key))
}