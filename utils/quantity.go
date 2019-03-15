package utils

import ("strings"
"log"
)

func ParseQuantity(quantity string) string {
	q := strings.Split(quantity, " ")[0]
	p :=  strings.Split(q, ".")
log.Printf("p[0] is %v, p[1] is %v\n", p[0], p[1])
	qs := p[0]+p[1]
log.Printf("qs is %v\n", qs)
	return p[0]+p[1]
}
