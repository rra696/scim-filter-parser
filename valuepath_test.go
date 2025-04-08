package filter

import "fmt"

func ExampleParseValuePath() {
	fmt.Println(ParseValuePath([]byte("emails[type eq \"work\"]")))
	fmt.Println(ParseValuePath([]byte("emails[not (type eq \"work\")]")))
	fmt.Println(ParseValuePath([]byte("emails[type eq \"work\" and value co \"@example.com\"]")))
	fmt.Println(ParseValuePath([]byte("emails[type eq \"work\" and value co \"@example.com\" or primary eq true]")))
	// Output:
	// emails[type eq "work"] <nil>
	// emails[not(type eq "work")] <nil>
	// emails[type eq "work" and value co "@example.com"] <nil>
	// emails[type eq "work" and value co "@example.com" or primary eq true] <nil>
}
