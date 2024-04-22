package html

import (
	"io"
	"net/http"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

const userAgent = "cgb-parser/0.0.1"

func FetchURL(url string) (io.ReadCloser, error) {
	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func findAttr(node *html.Node, attrName string) string {
	for _, a := range node.Attr {
		if a.Key == attrName {
			return a.Val
		}
	}

	return ""
}

func GetClasses(node *html.Node) []string {
	return strings.Split(findAttr(node, "class"), " ")
}

func HasClass(node *html.Node, class string) bool {
	classes := GetClasses(node)
	return slices.Contains(classes, class)
}

func GetId(node *html.Node) string {
	return findAttr(node, "id")
}

func FindNodes(root *html.Node, match func(*html.Node) bool) []*html.Node {
	matches := []*html.Node{}
	for c := root.FirstChild; c != nil; c = c.NextSibling {
		if match(c) {
			matches = append(matches, c)
		}
		childMatches := FindNodes(c, match)
		if childMatches != nil {
			matches = append(matches, childMatches...)
		}
	}

	if len(matches) == 0 {
		return nil
	}
	return matches
}

func FindNode(root *html.Node, match func(*html.Node) bool) *html.Node {
	for c := root.FirstChild; c != nil; c = c.NextSibling {
		if match(c) {
			return c
		}
		if n := FindNode(c, match); n != nil {
			return n
		}
	}

	return nil
}

func DumpNode(node *html.Node) (string, error) {
	strBuilder := &strings.Builder{}
	if err := html.Render(strBuilder, node); err != nil {
		return "", err
	}
	return strBuilder.String(), nil
}
