package utils

import "strings"

type MultiValueBuilder struct {
	Sep    string
	Values []string
}

func (b *MultiValueBuilder) Add(val string) {
	b.Values = append(b.Values, val)
}

func (b *MultiValueBuilder) Len() int {
	return len(b.Values)
}

func (b *MultiValueBuilder) Encode() string {
	return strings.Join(b.Values, b.Sep)
}
