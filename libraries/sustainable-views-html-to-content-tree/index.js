import {fromHtml as hastFromString} from "hast-util-from-html"
import {toString as hastToString} from "hast-util-to-string"

export default function convert(article) {
  let tree = {
    type: "root",
    data: {
      title: article.title,
      uuid: article.uuid,
      standfirst: ""
    },
    body: {
      type: "body",
      children: []
    }
  }

  let hast = hastFromString(article.body)
  // we love to wrap an impure function in a closure so nobody knows we're
  // impure, don't we folks?
  function walk(node) {
    if (node.type == "element") {
      if (node.tagName == "body") {
	let [first] = node.children
	if (first && first.tagName == "h2") {
	  // hell yeah i'm mutating do you wnna fight about it
	  node.children.shift()
	  // todo is this ever more than a string? emphasis?
	  tree.data.standfirst = hastToString(first)
	}
      }
    }
    if ("children" in node) {
      for (let child of node.children) {
	tree.body.children.push(walk(child))
      }
    }
  }
  walk(hast)
  return tree
}

function hastToContentTree(hnode) {
  let tree = {
    type: "root",
    data: {
      title: "",
      uuid: "",
      standfirst: ""
    }
  }
  let cnode = {
    type: [node.type]
  }

  if (node.type == "root") {
    return
  }
}
