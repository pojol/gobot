import { IsScriptNode, NodeTy } from "../model/node_type";

function parseChildren(childrenNodes, children) {
    var nod = {};
    var childrens = []

    childrenNodes.forEach(children => {
        if (children.nodeName === "id") {
            nod.id = children.childNodes[0].nodeValue
        } else if (children.nodeName === "ty") {
            nod.ty = children.childNodes[0].nodeValue
        } else if (children.nodeName === "pos") {
            nod.pos = {}
            children.childNodes.forEach(pos => {
                if (pos.nodeName === "x") {
                    nod.pos.x = parseInt(pos.childNodes[0].nodeValue)
                } else if(pos.nodeName === "y") {
                    nod.pos.y = parseInt(pos.childNodes[0].nodeValue)
                }
            })
        } else if (children.nodeName === "children") {
            childrens.push(children.childNodes)
        } else if (children.nodeName === "loop") {
            nod.loop = parseInt(children.childNodes[0].nodeValue)
        } else if (children.nodeName === "wait") {
            nod.wait = parseInt(children.childNodes[0].nodeValue)
        } else if (children.nodeName === "code") {
            nod.code = children.childNodes[0].nodeValue
        } else if (children.nodeName === "alias") {
            if (children.childNodes && children.childNodes.length) {
                nod.alias = children.childNodes[0].nodeValue
            } 
        }
    })

    nod.children = [];
    children.push(nod);

    console.info(nod.id, nod.ty, childrens)
    childrens.forEach(c => {
        parseChildren(c, nod.children)
    })
}


function LoadBehaviorWithBlob(url, methon, name) {
    return new Promise(function (resolve, reject) {
        fetch(url +"/"+ methon, {
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
                resolve({
                    name: name,
                    blob: response
                });
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
            var children = undefined
            if (root) {

                root.childNodes.forEach(nod => {
                    if (nod.nodeName === "id") {
                        tree.id = nod.childNodes[0].nodeValue
                    } else if (nod.nodeName === "ty") {
                        tree.ty = nod.childNodes[0].nodeValue
                    } else if (nod.nodeName === "pos") {
                        tree.pos = {}
                        nod.childNodes.forEach(pos => {
                            if (pos.nodeName === "x") {
                                tree.pos.x = parseInt(pos.childNodes[0].nodeValue)
                            } else if(pos.nodeName === "y") {
                                tree.pos.y = parseInt(pos.childNodes[0].nodeValue)
                            }
                        })
                    } else if (nod.nodeName === "children") {
                        children = nod.childNodes
                    }
                })
                tree.children = [];

                console.info("root", tree)

                if (children !== undefined) {
                    parseChildren(
                        children,
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

export { LoadBehaviorWithBlob, LoadBehaviorWithFile }