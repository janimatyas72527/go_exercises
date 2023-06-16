package exercise_4

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestCollectContent(t *testing.T) {
	// Arrange
	var content_nodes [3]*html.Node

	// Create Text Nodes
	for i := range content_nodes {
		content_nodes[i] = &html.Node{
			Type: html.TextNode,
			Data: fmt.Sprint("\r\n\tTest text \r\n\t#", i),
		}
	}
	// Link nodes together to form a small hierarchy
	content_nodes[0].NextSibling = content_nodes[2]
	content_nodes[0].FirstChild = content_nodes[1]
	content_nodes[0].LastChild = content_nodes[1]
	content_nodes[1].Parent = content_nodes[0]
	content_nodes[2].PrevSibling = content_nodes[0]

	// Create Anchor Node
	anchor_node := &html.Node{
		Type:       html.ElementNode,
		Data:       "a",
		Attr:       []html.Attribute{{Key: "href", Val: "/test"}, {Key: "data-key", Val: "value"}},
		FirstChild: content_nodes[0],
		LastChild:  content_nodes[2],
	}

	// Set parent by context nodes
	content_nodes[0].Parent = anchor_node
	content_nodes[2].Parent = anchor_node

	// Act

	var sb strings.Builder

	CollectContent(anchor_node, &sb)

	// Assert

	expected_content := "Test text #0 Test text #1 Test text #2"
	actual_content := sb.String()

	assert.Equal(t, expected_content, actual_content)
}

func TestGetHref(t *testing.T) {
	// Arrange
	const expected_href_1 = "/test"

	// Act

	href_1 := GetHref([]html.Attribute{{Key: "href", Val: expected_href_1}, {Key: "data-key", Val: "value"}})
	href_2 := GetHref([]html.Attribute{})

	// Assert

	assert.Equal(t, expected_href_1, href_1)
	assert.Empty(t, href_2)
}

func TestWriteData(t *testing.T) {
	// Arrange

	fs := afero.NewMemMapFs()
	f, f_err := fs.OpenFile("test_file", os.O_CREATE+os.O_APPEND, os.ModeAppend)

	if f_err != nil {
		panic(f_err)
	}

	test_link := "/test_link"
	test_text := "test_text"
	test_result := fmt.Sprint("Link {\n\tHref: ", test_link, "\n\tText:", test_text, "\n}\n")
	test_buf_length := len(test_result)
	test_buf := make([]byte, test_buf_length)

	// Act

	WriteData(test_link, test_text, f)
	f.Seek(0, 0)
	f.Read(test_buf)
	f.Close()

	read_data := string(test_buf)

	// Assert

	assert.Equal(t, test_result, read_data)
}

func TestProcessNode(t *testing.T) {
	// Arrange

	list_item_node := &html.Node{
		Type: html.ElementNode,
		Data: "li",
	}

	anchor_node := &html.Node{
		Type: html.ElementNode,
		Data: "a",
	}

	anchor_node.Parent = list_item_node
	list_item_node.FirstChild = anchor_node

	fs := afero.NewMemMapFs()
	f, f_err := fs.OpenFile("test_file", os.O_CREATE+os.O_APPEND, os.ModeAppend)

	if f_err != nil {
		panic(f_err)
	}

	test_result := "Link {\n\tHref: \n\tText:\n}\n"
	test_buf_length := len(test_result)
	test_buf := make([]byte, test_buf_length)

	// Act

	ProcessNode(list_item_node, f)
	f.Seek(0, 0)
	f.Read(test_buf)
	f.Close()

	read_data := string(test_buf)

	// Assert

	assert.Equal(t, test_result, read_data)
}

func TestExtractLinks_Standard_Path(t *testing.T) {
	// Arrange

	content, content_err := os.ReadFile("ex1.html")

	if content_err != nil {
		panic(content_err)
	}

	content_string := string(content)
	test_server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, content_string)
			}))

	defer test_server.Close()

	test_result := "Link {\n\tHref: /other-page\n\tText:A link to another page\n}\n"
	test_buf_length := len(test_result)
	test_buf := make([]byte, test_buf_length)
	fs := afero.NewMemMapFs()

	// Act

	ExtractLinks(test_server.URL, "test_file", fs)

	f, f_err := fs.Open("test_file")

	if f_err != nil {
		panic(f_err)
	}
	f.Seek(0, 0)
	f.Read(test_buf)
	f.Close()

	read_data := string(test_buf)

	// Assert
	assert.Equal(t, test_result, read_data)
}

func TestExtractLinks_HTTP_Error(t *testing.T) {
	// Arrange

	// Act-Assert

	assert.Panics(t, func() {
		ExtractLinks("NonStandardURL", "target_file", nil)
	})
}

func TestExtractLinks_File_System_Error(t *testing.T) {
	// Arrange

	test_server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "<html></html>")
			}))

	defer test_server.Close()

	// Make file system read-only
	fs := afero.NewReadOnlyFs(afero.NewMemMapFs())

	// Act-Assert

	assert.Panics(t, func() {
		ExtractLinks(test_server.URL, "target_file", fs)
	})
}
