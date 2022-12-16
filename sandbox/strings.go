//
//
//

/*
 */
package main

func Concat(a, b string) string {
	c := len([]byte(a))
	c += len([]byte(b))
	m := make([]byte, 0, c)
	m = append(m, a...)
	m = append(m, b...)
	return string(m)
}
