package exercise_4

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Collect content recursively
func CollectContent(node *html.Node, sb *strings.Builder) {
	if node.Type == html.TextNode {
		sb.WriteString(" ") // Space before data

		// Remove newline and tabs from text
		text := strings.ReplaceAll(node.Data, "\r", "")
		text = strings.ReplaceAll(text, "\n", "")
		text = strings.ReplaceAll(text, "\t", "")

		sb.WriteString(text)
	}
	// Collect content from child nodes, too
	child := node.FirstChild
	for child != nil {
		CollectContent(child, sb)
		child = child.NextSibling
	}
}

func GetHref(attr []html.Attribute) string {
	for i := range attr {
		if attr[i].Key == "href" {
			return attr[i].Val
		}
	}
	return ""
}

func WriteData(href string, content string, tf *os.File) {
	tf.WriteString(fmt.Sprint("Link {\n\tHref: ", href, "\n\tText:", content, "\n}\n"))
}

// Process nodes recursively
func ProcessNode(node *html.Node, tf *os.File) {
	if node.Type == html.ElementNode && node.Data == "a" {
		var sb strings.Builder

		CollectContent(node, &sb)

		href := GetHref(node.Attr)

		WriteData(href, sb.String(), tf)
	}
	// Process child nodes, too
	child := node.FirstChild
	for child != nil {
		ProcessNode(child, tf)
		child = child.NextSibling
	}
}

func ExtractLinks(source_url string, target_file string) {
	response, response_err := http.Get(source_url)

	if response_err != nil {
		panic(response_err)
	}

	document, document_err := html.Parse(response.Body)

	response.Body.Close()
	if document_err != nil {
		panic(document_err)
	}

	tf, tf_err := os.OpenFile(target_file, os.O_CREATE+os.O_APPEND, os.ModeAppend)

	if tf_err != nil {
		panic(tf_err)
	}
	ProcessNode(document, tf)
	tf.Close()
}
