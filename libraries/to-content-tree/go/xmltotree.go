package xmltotree

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/beevik/etree"
	"io"
	"strings"

	ct "github.com/Financial-Times/content-tree"
)

type XmlToTree struct {
}

func NewXmlToTree() *XmlToTree {
	return &XmlToTree{}
}

func (xt *XmlToTree) Transform(bodyXML string) (string, error) {
	return FromXMLString(bodyXML)
}

// FromXMLString parses XML from a string and returns a ContentTree root.
func FromXMLString(s string) (string, error) {
	return FromETreeReader(strings.NewReader(s))
}

// FromXMLReader parses XML from an io.Reader and returns a ContentTree root.
func FromXMLReader(r io.Reader) (*ct.Root, error) {
	dec := xml.NewDecoder(r)
	for {
		token, err := dec.Token()
		if err != nil {
			return nil, err
		}
		if token == nil {
			break
		}

		switch t := token.(type) {
		case xml.StartElement:
			/*
				let transformer =
				        xmlnode.name == "content" || xmlnode.name == "ft-content"
				          ? String(xmlnode.attributes.type)
				          : xmlnode.name;
			*/
			transformerName := t.Name.Local
			if t.Name.Local == "content" || t.Name.Local == "ft-content" {
				for _, attr := range t.Attr {
					if attr.Name.Local == "type" {
						transformerName = attr.Value
						break
					}
				}
			}
			_, ok := DefaultTransformers[transformerName]
			if !ok {
				// return { type: "__UNKNOWN__", data: xmlnode }; -- TODO
				return nil, fmt.Errorf("unknown transformer: %s", transformerName)
			}

		case xml.EndElement:
			fmt.Printf("EndElement: name=%s\n", t.Name.Local)
		case xml.CharData:
			// trim spaces so you don't print a lot of whitespace-only data
			data := strings.TrimSpace(string(t))
			if data != "" {
				fmt.Printf("CharData: %q\n", data)
			}
		case xml.Comment:
			fmt.Printf("Comment: %q\n", string(t))
		case xml.ProcInst:
			fmt.Printf("ProcInst: target=%s\n", t.Target)
		case xml.Directive:
			fmt.Printf("Directive: %q\n", string(t))
		default:
			fmt.Println("Unknown token type")
		}
	}
	return nil, nil
}

func FromETreeReader(r io.Reader) (string, error) {
	doc := etree.NewDocument()
	_, err := doc.ReadFrom(r)
	if err != nil {
		return "", err
	}
	// Get the root element
	root := doc.Root()
	if root == nil {
		return "", fmt.Errorf("no root element found")
	}

	// Print root element name
	//fmt.Printf("Root: %s\n", root.Tag)

	// Iterate over child elements recursively
	m := map[string]any{}
	err = convertToContentTree(root, 0, m)
	if err != nil {
		return "", err
	}
	out := map[string]any{
		"type": "root",
		"body": map[string]any{
			"type":     "body",
			"version":  1,
			"children": m["children"],
		},
	}
	return convertToJSON(out)
}

func convertToContentTree(elem etree.Token, level int, m map[string]any) error {
	//indent := strings.Repeat("  ", level)
	switch t := elem.(type) {
	case *etree.Comment:
		// Ignoring
	case *etree.Element:
		tag := t.Tag
		if tag != "body" {
			if t.Tag == "content" || t.Tag == "ft-content" {
				for _, attr := range t.Attr {
					if attr.Key == "type" {
						tag = attr.Value
						break
					}
				}
			}
			transformer, ok := DefaultTransformers[tag]
			if !ok {
				//return { type: "__UNKNOWN__", data: xmlnode };
				//JS code likes to set the above object, we throw an error for now
				return fmt.Errorf("unknown transformer: %s", tag)
			}
			transformed := transformer(t)
			//to throw error for uncertainty inside div
			if transformed["type"] == "__UNKNOWN__" {
				return fmt.Errorf("unknown transformer: %s", transformed["class"])
			}
			if transformed["type"] != "__LIFT_CHILDREN__" {
				_, ok := transformed["children"]
				if !ok {
					transformed["children"] = []any{}
				}
				ch, ok := m["children"].([]map[string]any)
				if !ok {
					ch = []map[string]any{
						transformed,
					}
				} else {
					ch = append(ch, transformed)
				}
				m["children"] = ch
			}
			if transformed["type"] == "__LIFT_CHILDREN__" {
				/*
					if (ctnode.type === "__LIFT_CHILDREN__") {
					          // we don't want this node to stick around, but we want to keep its' children
					          return xmlnode.children.flatMap(walk);
					        }
				*/
				for _, child := range t.Child {
					err := convertToContentTree(child, level+1, m)
					if err != nil {
						return err
					}
				}
				return nil
			} else if transformed["children"] == DO_NOT_PROCESS_CHILDREN {
				/*
				 // this is how we indicate we shouldn't iterate, but this thing
				          // shouldn't have any children
				          delete ctnode.children;
				          return ctnode;
				*/
				delete(transformed, "children")
				return nil
			}
			/*
			 else if ("children" in ctnode && Array.isArray(ctnode.children)) {
			          return ctnode;
			        } else if ("children" in xmlnode) {
			          return {
			            ...ctnode,
			            // this is a flatmap because of <experimental/>
			            children: xmlnode.children.flatMap(walk),
			          };
			        }
			        return ctnode;
			*/

			for _, child := range t.Child {
				err := convertToContentTree(child, level+1, transformed)
				if err != nil {
					return err
				}
			}
			return nil
		}
		for _, child := range t.Child {
			err := convertToContentTree(child, level+1, m)
			if err != nil {
				return err
			}
		}

	case *etree.CharData:
		data := t.Data
		tx := map[string]any{
			"value": data,
			"type":  "text",
		}
		ch, ok := m["children"].([]map[string]any)
		if !ok {
			ch = []map[string]any{
				tx,
			}
		} else {
			ch = append(ch, tx)
		}
		m["children"] = ch
	}
	return nil
}

func convertToJSON(m map[string]any) (string, error) {
	marshal, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(marshal))
	return string(marshal), nil
}
