package main

import (
	"fmt"
	"regexp"
)

var r, re *regexp.Regexp

func init() {
	//capture arguments including entire sections wrapped in `"` or `'`
	r = regexp.MustCompile(`\\"(.*?)\\"|\"(.*?)\"|'(.*?)'|\S+`)
	//replace escape character before `"`
	re = regexp.MustCompile(`\\"|\"|'`)
}

func testParse() {
	input := []string{
		`bash -c \"sleep 2; cat 'go.mod'\"`,
		`bash -c 'sleep 2; cat go.mod'`,
		`bash -c "sleep 2; cat go.mod"`,
		`bash -c sleep\ 2;cat\ go.mod`,
		`ping google.com`,
		`ping -r \"google.com\"`,
		`ping -r 'google.com'`,
		`ping a`,
		`ping --as -a b.a.a`,
	}
	procc := make([][]string, 9)
	procc[0] = r.FindAllString(input[0], -1)
	procc[1] = r.FindAllString(input[1], -1)
	procc[2] = r.FindAllString(input[2], -1)
	procc[3] = r.FindAllString(input[3], -1)
	procc[4] = r.FindAllString(input[4], -1)
	procc[5] = r.FindAllString(input[5], -1)
	procc[6] = r.FindAllString(input[6], -1)
	procc[7] = r.FindAllString(input[7], -1)
	procc[8] = r.FindAllString(input[8], -1)

	for _, item := range procc {
		for i, str := range item {
			//deal with `\"`
			item[i] = re.ReplaceAllString(str, ``)
		}
		fmt.Printf("%#v\n", item)
	}
}

// parseInput(input []string){

// }
