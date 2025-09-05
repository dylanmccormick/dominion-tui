package assert

import "log"

var assertData map[string]any = map[string]any{}

func runAssert(msg string) {
	for k, v := range assertData {
		log.Println("context", "key", k, "value", v)
	}
	log.Fatal(msg)
}

func Assert(truth bool, msg string) {
	if !truth {
		runAssert(msg)
	}
}

