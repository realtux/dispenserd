package main

func UtilInArray(a string, arr []string) bool {
    for _, b := range arr {
        if b == a {
            return true
        }
    }
    return false
}
