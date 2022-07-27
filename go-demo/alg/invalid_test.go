package alg

import "testing"

func TestInvalid(t *testing.T) {
	t.Log(isValid("()"))
}

//func isValid(s string) bool {
//	var stack = make([]string, 0)
//	for i := 0; i < len(s); i++ {
//		c := string(s[i])
//		if c == "(" {
//			stack = append(stack, ")")
//		} else if c == "[" {
//			stack = append(stack, "]")
//		} else if c == "{" {
//			stack = append(stack, "}")
//		} else if len(stack) != 0 && pop(&stack) != c {
//			return false
//		} else if len(stack) == 0 {
//			return false
//		}
//	}
//
//	return len(stack) == 0
//}
//
//func pop(s *[]string) string {
//	l := len(*s)
//	if l == 0 {
//		return ""
//	}
//
//	last := (*s)[l-1]
//	*s = (*s)[0 : l-1]
//	return last
//}

//type stack []string
//
//func (s stack) Pop() string {
//	return s[len(s)-1]
//}

func isValid(s string) bool {
	n := len(s)
	if n%2 == 1 {
		return false
	}
	pairs := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}
	stack := []byte{}
	for i := 0; i < n; i++ {
		if pairs[s[i]] > 0 {
			if len(stack) == 0 || stack[len(stack)-1] != pairs[s[i]] {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, s[i])
		}
	}
	return len(stack) == 0
}
