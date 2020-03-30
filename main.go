package main

import (
	"encoding/xml"
	"flag"
	"io/ioutil"
	"os"
)

type opml struct {
	XMLName xml.Name `xml:"opml"`
	Body    *body    `xml:"body"`
}

type body struct {
	XMLName  xml.Name  `xml:"body"`
	Outlines []outline `xml:"outline"`
}

type outline struct {
	XMLName xml.Name `xml:"outline"`
	Type    string   `xml:"type,attr"`
	Text    string   `xml:"text,attr"`
	XMLUrl  string   `xml:"xmlUrl,attr"`
}

func main() {
	mergedTo := flag.String("p", "./feedlist.opml", "The opml file with higher priority")
	merged := flag.String("n", "./podcast_republic_podcasts.opml", "The opml file with lower priority")
	flag.Parse()
	src1, err := os.Open(*mergedTo)
	if err != nil {
		panic(err)
	}
	src2, err := os.Open(*merged)
	if err != nil {
		panic(err)
	}
	defer src1.Close()
	defer src2.Close()

	text1, _ := ioutil.ReadAll(src1)
	text2, _ := ioutil.ReadAll(src2)

	var v1, v2 opml
	xml.Unmarshal(text1, &v1)
	xml.Unmarshal(text2, &v2)
	v1.Body.Outlines = union(v1.Body.Outlines, v2.Body.Outlines)
	f, _ := xml.MarshalIndent(v1, "", "    ")
	output, _ := os.Create("merge.opml")
	output.Write(f)
	output.Sync()
	output.Close()
}

func union(l1, l2 []outline) []outline {
	merged := append(l1, make([]outline, 0)...)

	for _, i := range l2 {
		if !exist(i, merged) {
			merged = append(merged, i)
		}
	}

	return merged
}

func exist(l outline, L []outline) bool {
	for _, i := range L {
		if l.XMLUrl == i.XMLUrl {
			return true
		}
	}
	return false
}
