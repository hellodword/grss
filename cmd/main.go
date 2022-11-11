package main

import (
	"fmt"
	"github.com/hellodword/grss"
	"os"
	"strings"
)

func main() {
	s := `
<?xml version="1.0"?>
<rdf:RDF
xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
xmlns="http://channel.netscape.com/rdf/simple/0.9/">

  <channel>
    <title>Mozilla Dot Org</title>
    <link>http://www.mozilla.org</link>
    <description>the Mozilla Organization
      web site</description>
  </channel>

  <image>
    <title>Mozilla</title>
    <url>http://www.mozilla.org/images/moz.gif</url>
    <link>http://www.mozilla.org</link>
  </image>

  <item>
    <title>New Status Updates</title>
    <link>http://www.mozilla.org/status/</link>
  </item>

  <item>
    <title>Bugzilla Reorganized</title>
    <link>http://www.mozilla.org/bugs/</link>
  </item>

  <item>
    <title>Mozilla Party, 2.0!</title>
    <link>http://www.mozilla.org/party/1999/</link>
  </item>

  <item>
    <title>Unix Platform Parity</title>
    <link>http://www.mozilla.org/build/unix.html</link>
  </item>

  <item>
    <title>NPL 1.0M published</title>
    <link>http://www.mozilla.org/NPL/NPL-1.0M.html</link>
  </item>

</rdf:RDF>
`

	_, f, err := grss.Parse(strings.NewReader(s))
	if err != nil {
		panic(err)
	}

	{
		a, err := f.ToAtom()
		if err != nil {
			panic(err)
		}
		if a != nil {
			fmt.Println("atom")
			a.WriteOut(os.Stdout)
		}
	}

	{
		a, err := f.ToJSON()
		if err != nil {
			panic(err)
		}
		if a != nil {
			fmt.Println("json")
			a.WriteOut(os.Stdout)
		}
	}

	{
		a, err := f.ToRss2()
		if err != nil {
			panic(err)
		}
		if a != nil {
			fmt.Println("rss")
			a.WriteOut(os.Stdout)
		}
	}

}
