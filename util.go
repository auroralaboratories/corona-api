package main

func contains(set []string, value string) bool {
    for _, i := range set {
      if i == value {
        return true
      }
    }

    return false
}