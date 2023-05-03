
import { IsActionNode, NodeTy } from "@/constant/node_type";
import { Graph, Addon, Dom, Node } from "@antv/x6";



export function getDefaultNodeNotifyInfo(): NodeNotifyInfo {
    return {
        id: "",
        ty: "",
        code: "",
        loop: 1,
        wait: 1,
        pos: {
            x: 0,
            y: 0,
        },
        children: [],
        notify: false,
        alias: "",
    };
}

function fillChildInfo(child: Node, info: any) {
    var childInfo = {
      id: child.id,
      ty: child.getAttrs().type.toString(),
      pos: {
        x: child.position().x,
        y: child.position().y,
      },
      children: [],
    };
    info.children.push(childInfo);
  
    child.eachChild((cchild, idx) => {
      if (cchild instanceof Node) {
        fillChildInfo(cchild as Node, childInfo);
      }
    });
  }
  
  export function GetNodInfo(prefab: Array<PrefabInfo>, nod: Node, code: string, alias: string): NodeNotifyInfo {
    var info = getDefaultNodeNotifyInfo();
  
    info.id = nod.id;
    info.ty = nod.getAttrs().type.name as string;
    info.pos = {
      x: nod.position().x,
      y: nod.position().y,
    };
  
    if (info.ty === NodeTy.Action || IsActionNode(info.ty)) {
      if (code !== "") {
        info.code = code
      } else {
        prefab.forEach((p) => {
          if (p.name === info.ty) {
            info.code = p.code;
          }
  
          if (alias === "") {
            info.alias = info.ty
          }
        })
  
        if (info.code === "") {
  
        }
      }
  
      if (alias !== "") {
        info.alias = alias
      }
    }
  
    if (info.ty === NodeTy.Condition) {
      if (code !== "") {
        info.code = code
      } else {
        info.code = `
  -- Write expression to return true or false
  function execute()
  
  end
        `;
      }
    }
  
    nod.eachChild((child, idx) => {
      if (child instanceof Node) {
        fillChildInfo(child as Node, info);
      }
    });
  
    return info;
  }