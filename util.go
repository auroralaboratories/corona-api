package main

import "strings"

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

func compact(set []string) []string {
  rv := make([]string, 0, len(set))

  for _, v := range set {
    if strings.TrimSpace(v) == "" {
      continue
    }

    rv = append(rv, v)
  }

  return rv
}

func Stosl(s string)[]string{
  return strings.Split(s, ";")
}

func Stob(s string) bool {
  if s == "true"{
    return true
  }
  return false
}