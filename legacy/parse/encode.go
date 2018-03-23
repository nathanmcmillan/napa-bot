package parse

import "bytes"

// Begin starts json
func Begin(b *bytes.Buffer) {
	b.WriteString(`{`)
}

// End ends json
func End(b *bytes.Buffer) {
	b.WriteString(`}`)
}

// First first element of json
func First(b *bytes.Buffer, name, value string) {
	b.WriteString(`"`)
	b.WriteString(name)
	b.WriteString(`"`)
	b.WriteString(`:`)
	b.WriteString(`"`)
	b.WriteString(value)
	b.WriteString(`"`)
}

// Append second or later element of json
func Append(b *bytes.Buffer, name, value string) {
	b.WriteString(`, `)
	First(b, name, value)
}
