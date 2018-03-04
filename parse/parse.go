
package parse

func Text(js map[string]interface{}, key string) string {
    str, ok := js[key].(string)
    if ok {
        return str
    }
    return ""
}

func Number(js map[string]interface{}, key string) float64 {
    num, ok := js[key].(float64)
    if ok {
        return num
    }
    return float64(0.0)
}

func Integer(js map[string]interface{}, key string) int64 {
    return int64(Number(js, key))
}

func StringList(js map[string]interface{}, key string) []string {
    list := make([]string, 0)
    raw, ok := js[key].([]interface{})
    if !ok {
        return list
    }
    for i := 0; i < len(raw); i++ {
        str, ok := raw[i].(string)
        if !ok {
            continue
        }
        list = append(list, str)
    }
    return list
}
