package main

import (
	"fmt"
	"io"
	"os"
	"encoding/xml"
	"strings"
	"regexp"
)

// this is a comment

func readFile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	text := ""
	for {
		buf := make([]byte, 1024)
		// read a chunk
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		text += fmt.Sprintf("%s", buf)
	}
	return text
}

type Elem struct {
	Name     string
	Text     string
	Attributes []xml.Attr
	Children []*Elem
	isTextNode bool
}

func parse(html string) *Elem {
	var tree *Elem
	var cursor *Elem
	var stack []*Elem
	decoder := xml.NewDecoder(strings.NewReader(html))
	decoder.Strict = false
	decoder.AutoClose = xml.HTMLAutoClose
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}

		switch entity := token.(type) {
		case xml.StartElement:
			cursor = new(Elem)
			cursor.Name = entity.Name.Local
			cursor.Attributes = entity.Attr
			if (tree != nil) {
				tree.Children = append(tree.Children, cursor)
			}
			tree = cursor
			stack = append(stack, cursor)
		case xml.EndElement:
			if tree.Name != entity.Name.Local {
				panic("Tag names mismatch: got " + entity.Name.Local + ", expected " + tree.Name)
			}
			stack = stack[:len(stack) - 1]
			if len(stack) > 0 {
				tree = stack[len(stack) - 1]
			}
		case xml.CharData:
			if (tree != nil) {
				cursor = new(Elem)
				cursor.isTextNode = true
				cursor.Text = fmt.Sprintf("%s", entity)
				tree.Children = append(tree.Children, cursor)
			}
		default:
			//fmt.Printf("%#v\n", token)
		}
	}
	return tree
}

func getAttr(attributes []xml.Attr, name string) string {
	name = strings.ToLower(name)
	for _, elem := range attributes {
		if strings.ToLower(elem.Name.Local) == name {
			return elem.Value
		}
	}
	return ""
}

func renderMarkdown(tree *Elem) string {

	return renderRecursive(tree, false, 0)
}

func renderRecursive(tree *Elem, inBody bool, listOrder int) string {

	if tree.isTextNode {
		if inBody {
			return minify(tree.Text)
		}
		return ""
	}

	goDeeper := true
	template := "%s"
	switch tree.Name {
	case "body":
		inBody = true
	case "a":
		template = "[%s](" + getAttr(tree.Attributes, "href") + ")"
	case "hr":
		template = "* * *\n\n%s"
		goDeeper = false
	case "p":
		template = "%s\n\n";
	case "s":
		template = "~~%s~~";
	case "i":
		template = "*%s*";
	case "em":
		template = "*%s*";
	case "b":
		template = "**%s**";
	case "strong":
		template = "**%s**";
	case "ul":
		listOrder = 0
		template = "%s\n\n"
	case "ol":
		listOrder = 1
		template = "%s\n\n"
	case "li":
		if listOrder > 0 {
			template = fmt.Sprintf("%d. %%s\n", listOrder)
		} else {
			template = "* %s\n"
		}

	case "h1":
		template = "# %s\n\n";
	case "h2":
		template = "## %s\n\n";
	case "h3":
		template = "### %s\n\n";
	case "h4":
		template = "#### %s\n\n";
	case "h5":
		template = "##### %s\n\n";
	case "h6":
		template = "##### %s\n\n";
	case "h7":
		template = "###### %s\n\n";
	}

	content := ""
	if goDeeper {
		index := 1
		for _, elem := range tree.Children {
			content += renderRecursive(elem, inBody, listOrder * index)
			if !elem.isTextNode {
				index++
			}
		}
	}

	return trim(fmt.Sprintf(template, content))
}

func trim(text string) string {
	r , _ := regexp.Compile("(?m:^ *(.+?) *$)")
	return r.ReplaceAllString(text, "$1")
}

func minify(text string) string {
	r , _ := regexp.Compile("\\s+")
	return r.ReplaceAllString(text, " ")
}

func main() {
	if len(os.Args) <= 1 {
		os.Exit(2)
	}
	html := readFile(os.Args[1])
	tree := parse(html)
	fmt.Printf("%s\n", renderMarkdown(tree))
}
