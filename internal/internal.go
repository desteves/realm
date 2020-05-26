// Package internal converts query or mutation to a string
package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// $GOPATH/src/github.com/shurcoo!l/graphql@v0.0.0-20181231061246-d48a9a75455f/query.go

// Name is an identifier name, broken up into individual words.
type Name []string

// ConstructQuery takes the query interface along with any variables to produce a valid query in string format.
func ConstructQuery(v interface{}, variables map[string]interface{}) string {
	query := query(v)
	if variables == nil || len(variables) == 0 {
		return query
	}
	return "query(" + queryArguments(variables) + ")" + query

}

// func ConstructMutation(v interface{}, variables map[string]interface{}) string {
// 	query := query(v)
// 	if len(variables) > 0 {
// 		return "mutation(" + queryArguments(variables) + ")" + query
// 	}
// 	return "mutation" + query
// }

// queryArguments constructs a minified arguments string for variables.
//
// E.g., map[string]interface{}{"a": Int(123), "b": NewBoolean(true)} -> "$a:Int!$b:Boolean".
func queryArguments(variables map[string]interface{}) string {
	// Sort keys in order to produce deterministic output for testing purposes.
	// TODO: If tests can be made to work with non-deterministic output, then no need to sort.
	keys := make([]string, 0, len(variables))
	for k := range variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		io.WriteString(&buf, "$")
		io.WriteString(&buf, k)
		io.WriteString(&buf, ":")
		writeArgumentType(&buf, reflect.TypeOf(variables[k]), true)
		// Don't insert a comma here.
		// Commas in GraphQL are insignificant, and we want minified output.
		// See https://facebook.github.io/graphql/October2016/#sec-Insignificant-Commas.
	}
	return buf.String()
}

// writeArgumentType writes a minified GraphQL type for t to w.
// value indicates whether t is a value (required) type or pointer (optional) type.
// If value is true, then "!" is written at the end of t.
func writeArgumentType(w io.Writer, t reflect.Type, value bool) {
	if t.Kind() == reflect.Ptr {
		// Pointer is an optional type, so no "!" at the end of the pointer's underlying type.
		writeArgumentType(w, t.Elem(), false)
		return
	}

	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		// List. E.g., "[Int]".
		io.WriteString(w, "[")
		writeArgumentType(w, t.Elem(), true)
		io.WriteString(w, "]")
	default:
		// Named type. E.g., "Int".
		name := t.Name()
		if name == "string" { // HACK: Workaround for https://github.com/shurcooL/githubv4/issues/12.
			name = "ID"
		}
		io.WriteString(w, name)
	}

	if value {
		// Value is a required type, so add "!" to the end.
		io.WriteString(w, "!")
	}
}

// query uses writeQuery to recursively construct
// a minified query string from the provided struct v.
//
// E.g., struct{Foo Int, BarBaz *Boolean} -> "{foo,barBaz}".
func query(v interface{}) string {
	var buf bytes.Buffer
	writeQuery(&buf, reflect.TypeOf(v), false)
	return buf.String()
}

// writeQuery writes a minified query for t to w.
// If inline is true, the struct fields of t are inlined into parent struct.
func writeQuery(w io.Writer, t reflect.Type, inline bool) {
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice:
		writeQuery(w, t.Elem(), false)
	case reflect.Struct:
		// If the type implements json.Unmarshaler, it's a scalar. Don't expand it.
		if reflect.PtrTo(t).Implements(jsonUnmarshaler) {
			return
		}
		if !inline {
			io.WriteString(w, "{")
		}
		for i := 0; i < t.NumField(); i++ {
			if i != 0 {
				io.WriteString(w, ",")
			}
			f := t.Field(i)
			value, ok := f.Tag.Lookup("graphql")
			inlineField := f.Anonymous && !ok
			if !inlineField {
				if ok {
					io.WriteString(w, value)
				} else {
					io.WriteString(w, parseMixedCaps(f.Name).toLowerCamelCase())
				}
			}
			writeQuery(w, f.Type, inlineField)
		}
		if !inline {
			io.WriteString(w, "}")
		}
	}
}

// parseMixedCaps parses a MixedCaps identifier name.
//
// E.g., "ClientMutationID" -> {"Client", "Mutation", "ID"}.
func parseMixedCaps(name string) Name {
	var words Name

	// Split name at any lower -> Upper or Upper -> Upper,lower transitions.
	// Check each word for initialisms.
	runes := []rune(name)
	w, i := 0, 0 // Index of start of word, scan.
	for i+1 <= len(runes) {
		eow := false // Whether we hit the end of a word.
		if i+1 == len(runes) {
			eow = true
		} else if unicode.IsLower(runes[i]) && unicode.IsUpper(runes[i+1]) {
			// lower -> Upper.
			eow = true
		} else if i+2 < len(runes) && unicode.IsUpper(runes[i]) && unicode.IsUpper(runes[i+1]) && unicode.IsLower(runes[i+2]) {
			// Upper -> Upper,lower. End of acronym, followed by a word.
			eow = true

			if string(runes[i:i+3]) == "IDs" { // Special case, plural form of ID initialism.
				eow = false
			}
		}
		i++
		if !eow {
			continue
		}

		// [w, i) is a word.
		word := string(runes[w:i])
		if initialism, ok := isInitialism(word); ok {
			words = append(words, initialism)
		} else if i1, i2, ok := isTwoInitialisms(word); ok {
			words = append(words, i1, i2)
		} else {
			words = append(words, word)
		}
		w = i
	}
	return words
}

// isInitialism reports whether word is an initialism.
func isInitialism(word string) (string, bool) {
	initialism := strings.ToUpper(word)
	_, ok := initialisms[initialism]
	return initialism, ok
}

// isTwoInitialisms reports whether word is two initialisms.
func isTwoInitialisms(word string) (string, string, bool) {
	word = strings.ToUpper(word)
	for i := 2; i <= len(word)-2; i++ { // Shortest initialism is 2 characters long.
		_, ok1 := initialisms[word[:i]]
		_, ok2 := initialisms[word[i:]]
		if ok1 && ok2 {
			return word[:i], word[i:], true
		}
	}
	return "", "", false
}

// initialisms is the set of initialisms in the MixedCaps naming convention.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var initialisms = map[string]struct{}{
	// These are the common initialisms from golint. Keep them in sync
	// with https://gotools.org/github.com/golang/lint#commonInitialisms.
	"ACL":   {},
	"API":   {},
	"ASCII": {},
	"CPU":   {},
	"CSS":   {},
	"DNS":   {},
	"EOF":   {},
	"GUID":  {},
	"HTML":  {},
	"HTTP":  {},
	"HTTPS": {},
	"ID":    {},
	"IP":    {},
	"JSON":  {},
	"LHS":   {},
	"QPS":   {},
	"RAM":   {},
	"RHS":   {},
	"RPC":   {},
	"SLA":   {},
	"SMTP":  {},
	"SQL":   {},
	"SSH":   {},
	"TCP":   {},
	"TLS":   {},
	"TTL":   {},
	"UDP":   {},
	"UI":    {},
	"UID":   {},
	"UUID":  {},
	"URI":   {},
	"URL":   {},
	"UTF8":  {},
	"VM":    {},
	"XML":   {},
	"XMPP":  {},
	"XSRF":  {},
	"XSS":   {},

	// Additional common initialisms.
	"RSS": {},
}

// toLowerCamelCase expresses identifer name in lowerCamelCase naming convention.
//
// E.g., "clientMutationId".
func (n Name) toLowerCamelCase() string {
	for i, word := range n {
		if i == 0 {
			n[i] = strings.ToLower(word)
			continue
		}
		r, size := utf8.DecodeRuneInString(word)
		n[i] = string(unicode.ToUpper(r)) + strings.ToLower(word[size:])
	}
	return strings.Join(n, "")
}

var jsonUnmarshaler = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
