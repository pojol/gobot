import { IsScriptNode, NodeTy } from "../model/node_type";



function getValueByElement(elem, tag) {
    for (var i = 0; i < elem.childNodes.length; i++) {
        if (elem.childNodes[i].nodeName === tag) {
            if (elem.childNodes[i].childNodes.length === 0) {
                return ""
            } else {
                return elem.childNodes[i].childNodes[0].nodeValue;
            }
        }
    }
    return undefined;
}

function parseChildren(xmlnode, children) {
    var nod = {};

    nod.id = xmlnode.getElementsByTagName("id")[0].childNodes[0].nodeValue;
    nod.ty = xmlnode.getElementsByTagName("ty")[0].childNodes[0].nodeValue;

    if (nod.ty === NodeTy.Loop) {
        nod.loop = getValueByElement(xmlnode, "loop")
    } else if (nod.ty === NodeTy.Wait) {
        nod.wait = getValueByElement(xmlnode, "wait")
    } else if (IsScriptNode(nod.ty)) {
        nod.code = getValueByElement(xmlnode, "code");
        nod.alias = getValueByElement(xmlnode, "alias");
    }

    nod.pos = {
        x: parseInt(
            xmlnode.getElementsByTagName("pos")[0].getElementsByTagName("x")[0]
                .childNodes[0].nodeValue
        ),
        y: parseInt(
            xmlnode.getElementsByTagName("pos")[0].getElementsByTagName("y")[0]
                .childNodes[0].nodeValue
        ),
    };

    nod.children = [];
    children.push(nod);

    for (var i = 0; i < xmlnode.childNodes.length; i++) {
        if (xmlnode.childNodes[i].nodeName === "children") {
            parseChildren(xmlnode.childNodes[i], nod.children);
        }
    }
}


function LoadBehaviorWithBlob(url, methon, name) {
    return new Promise(function (resolve, reject) {
        fetch(url + methon, {
            method: "POST",
            mode: "cors",
            headers: {
                "Content-Type": "application/x-www-form-urlencoded",
            },
            body: JSON.stringify({ Name: name }),
        })
            .then((response) => {
                if (response.ok) {
                    return response.blob();
                } else {
                    reject({ status: response.status });
                }
            })
            .then((response) => {
                resolve(response);
            })
            .catch((err) => {
                reject({ status: -1 });
            });
    });
}


function LoadBehaviorWithFile(name, blob) {
    let reader = new FileReader();
    let tree = {};

    reader.onload = function (ev) {
      var context = reader.result;
      try {
        let parser = new DOMParser();
        let xmlDoc = parser.parseFromString(context, "text/xml");
        
        var root = xmlDoc.getElementsByTagName("behavior")[0];
        if (root) {
          tree.id = root.getElementsByTagName("id")[0].childNodes[0].nodeValue;
          tree.ty = root.getElementsByTagName("ty")[0].childNodes[0].nodeValue;
          tree.pos = {
            x: parseInt(
              root.getElementsByTagName("pos")[0].getElementsByTagName("x")[0]
                .childNodes[0].nodeValue
            ),
            y: parseInt(
              root.getElementsByTagName("pos")[0].getElementsByTagName("y")[0]
                .childNodes[0].nodeValue
            ),
          };
          tree.children = [];
          if (root.getElementsByTagName("children")[0].hasChildNodes()) {
            parseChildren(
              root.getElementsByTagName("children")[0],
              tree.children
            );
          }
        }
  
      } catch (err) {
        console.info(err)
        return null
      }
    };
  
    reader.readAsText(blob);
    return tree
  }

  export {LoadBehaviorWithBlob, LoadBehaviorWithFile}