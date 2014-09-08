package main

func contains(set []string, value string) bool {
    for _, i := range set {
      if i == value {
        return true
      }
    }

    return false
}

func indexOf(set []string, value string) int {
    for idx, i := range set {
      if i == value {
        return idx
      }
    }

    return -1
}