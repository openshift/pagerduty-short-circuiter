package utils

import (
	"fmt"

	"github.com/rivo/tview"
	"golang.org/x/net/html"
)

var numLinks int = 0

func TraverseHTMLDoc(n *html.Node, textView *tview.TextView) int {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "a":
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					fmt.Fprintf(textView, `["%d"][blue]%s - [white][""]`, numLinks, attr.Val)
					numLinks++
				}
			}
		case "img":
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					fmt.Fprintf(textView, `["%d"]%s[""]`, numLinks, attr.Val)
					numLinks++
				}
			}
		case "h1", "h2", "h3", "h4", "h5", "h6":
			fmt.Fprintf(textView, `[yellow]`)
		default:
			fmt.Fprintf(textView, `[white]`)
		}
	} else if n.Type == html.TextNode {
		fmt.Fprintf(textView, "%s", n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TraverseHTMLDoc(c, textView)
	}
	return numLinks
}

func FetchHTMLContent(URL string, textView *tview.TextView) {
	textView.Clear()
	numLinks = 0
	owner, repo, path := getGitHubMdURL(URL)
	contents, err := GetGHReadme(owner, repo, path)
	if (err) != nil {
		ErrorLogger.Printf("Error while fetching readme contents. The error message was : %s", err)
	}
	body := ConvertMarkdownToHTML(contents)

	// Parse the HTML file
	doc, err := html.Parse(body)
	if (err) != nil {
		ErrorLogger.Printf("Error while fetching readme contents. The error message was : %s", err)
	}
	TraverseHTMLDoc(doc, textView)
	// TODO : Work on Parsing Non Markdown Files
	// doc, err := html.Parse(body)
	// if err != nil {
	// 	panic(err)
	// }
	// traverseHTMLDoc(doc, textView)
	// text, err := html2text.FromHTMLNode(doc, html2text.Options{})
	// fmt.Println(text)
	// fmt.Fprint(textView, text)
	// if err != nil {
	// 	panic(err)
	// }
}
