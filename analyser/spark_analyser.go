package analyser

import (
	"fmt"
	"errors"
	"strings"
	"bytes"
	"io"
	"golang.org/x/net/html"
	"htmlparser/models"
)

func ParseSparkDashboard(htmlContent string) {
	doc, _ := html.Parse(strings.NewReader(htmlContent))
	bn, err := GetTableBody(doc)
	if err != nil {
		return
	}

	tbody, err := findChild(bn, "tbody")

	if err != nil {
		return
	}

	_, err = browseTr(tbody)

	if err != nil {
		return
	}

	renderNode(tbody)
	//body := renderNode(tbody)
	//fmt.Println(body)
}


func GetTableBody(doc *html.Node) (*html.Node, error) {
	var res *html.Node

	var f func(*html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "table" {

			for i:=0 ; i< len(node.Attr); i++ {
				if node.Attr[i].Key == "id" {

					if node.Attr[i].Val == "completed-batches-table" {
						res = node
						fmt.Println(node.Attr[i].Val)

					}
				}
			}


		}

		if res == nil {
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				f(child)
			}
		}
	}
	f(doc)
	if res != nil {
		return res, nil
	}
	return nil, errors.New("Missing <table> in the node tree")
}

func findChild(doc *html.Node, tagName string) (*html.Node, error) {
	var res *html.Node
	var f func(*html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == tagName {
			res = node
		}

		if res == nil {
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				f(child)
			}
		}
	}
	f(doc)
	if res != nil {
		return res, nil
	}
	return nil, errors.New("Missing <" + tagName + "> in the node tree")
}
func browseTr(tr *html.Node) (*html.Node, error) {

	lignes := -1
	var batches []models.Batch

	for child := tr.FirstChild; child != nil;  child = child.NextSibling  {

		if child.Data == "tr" {

			td, err := findChild(child, "td")

			if err != nil {
				return nil, errors.New("Cannot find td")
			}

			batch, err := browseTd(td)

			if err != nil {
				return nil, nil
			}
			batches = append(batches, batch)
			lignes++
			batch = models.Batch{}
		}

		if lignes  == 0 {
			fmt.Println("Nb lignes :", lignes+1)
			fmt.Println(batches)
			return nil, nil
		}

	}
	fmt.Println("Nb lignes :", lignes+1)
	fmt.Println(batches)
	return nil, nil
}

func browseTd(td *html.Node) (models.Batch, error) {
	cols := 0
	batch := models.Batch{}

	for child := td; child != nil; child = child.NextSibling {

		if child.Data == "td"  {

			switch cols {
			case 0:
				val := renderNode(child.FirstChild)
				fmt.Println("ici", val)
				batch.BatchTime = val
			case 1:
				//val, _ := strconv.Atoi(strings.Trim(renderNode(child.FirstChild)," "))
				val := strings.Trim(renderNode(child.FirstChild)," ")

				fmt.Println("la", val)
				batch.InputSize = val + ""
			case 2:
				batch.SchedulingDelay = renderNode(child.FirstChild)
			case 3:
				batch.ProcessingTime = renderNode(child.FirstChild)
			case 4:
				batch.TotalDelay = renderNode(child.FirstChild)
			}
			fmt.Println("cols", cols)
			cols++
		} else {

			fmt.Println("lost", child.Attr)
		}
	}
	fmt.Println(batch)
	return batch, nil
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}
