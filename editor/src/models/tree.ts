import { NodeTy } from "@/constant/node_type";
import { createAction, createSlice, PayloadAction } from "@reduxjs/toolkit";
const { Post, PostBlob, PostGetBlob } = require("../utils/request");
import Cmd from "@/constant/cmd";
import { message } from "antd";
import OBJ2XML from "object-to-xml";
import Api from "@/constant/api";
import { getDefaultNodeNotifyInfo } from "./node";

interface TreeState {
    nodes: Array<NodeNotifyInfo>;
    history: Array<any>;
    rootid: string;
    currentTreeName: string;
    currentDebugTree: NodeNotifyInfo;
    currentDebugBot: string;
    currentClickNode: NodeClickInfo;
    currentLockedNode: NodeClickInfo;
    updatetick: number;
}


const initialState: TreeState = {
    nodes: new Array<NodeNotifyInfo>(),
    history: new Array<NodeNotifyInfo>(),
    rootid: "",
    currentTreeName: "",
    currentDebugBot: "",
    currentDebugTree: getDefaultNodeNotifyInfo(),
    currentClickNode: { id: "", type: "" },
    currentLockedNode: { id: "", type: "" },
    updatetick: 0,
};

function deepCopy<T>(obj: T): T {
    return JSON.parse(JSON.stringify(obj));
}

function add(state: TreeState, info: NodeAddInfo) {
    state.nodes.push(info.info)
}

function rmv(state: TreeState, id: string) {
    for (var i = 0; i < state.nodes.length; i++) {
        if (state.nodes[i].id === id) {
            state.nodes.splice(i, 1)
            break
        }

        _find(id, state.nodes[i], (parent: NodeNotifyInfo, target: NodeNotifyInfo, idx: number) => {
            console.info("remove node parent", parent.id, "target", target.id)
            parent.children.splice(idx, 1)
        })
    }
}

function _cut(state: TreeState, id: string): [NodeNotifyInfo, boolean] {
    let nod = getDefaultNodeNotifyInfo()
    let ok = false

    for (var i = 0; i < state.nodes.length; i++) {
        if (state.nodes[i].id === id) {
            nod = deepCopy(state.nodes[i])
            state.nodes.splice(i, 1)
            ok = true
            break
        }

        _find(id, state.nodes[i], (parent: NodeNotifyInfo, target: NodeNotifyInfo, idx: number) => {
            nod = deepCopy(target)
            parent.children.splice(idx, 1)
            ok = true
        })
    }
    return [nod, ok]
}

function _copy(state: TreeState, parentid: string, children: NodeNotifyInfo): void {
    for (var i = 0; i < state.nodes.length; i++) {
        if (state.nodes[i].id === parentid) {
            state.nodes[i].children.push(children)
            break
        }

        _find(parentid, state.nodes[i], (parent: NodeNotifyInfo, target: NodeNotifyInfo, idx: number) => {
            target.children.push(children)
        })
    }
}

function link(state: TreeState, parentid: string, childrenid: string) {
    let res = _cut(state, childrenid)
    if (res[1]) {
        _copy(state, parentid, res[0])
    } else {
        message.warning("link node err unknow children " + childrenid)
    }
}

// targetid 被断开连接的节点
function unlink(state: TreeState, targetid: string) {
    let res = _cut(state, targetid)
    if (res[1]) {
        state.nodes.push(res[0])
    }
}

export const UpdateType = {
    UpdateAll: "_update_all",
    UpdateAlias: "_update_alias",
    UpdateCode: "_update_code",
    UpdatePosition: "_update_position",
    UpdateLoop: "_update_loop",
    UpdateWait: "_update_wait",
}

function _update_all(cur: NodeNotifyInfo, up: NodeNotifyInfo): void {
    cur.code = up.code
    cur.alias = up.alias
    cur.pos = up.pos
    cur.loop = up.loop
    cur.wait = up.wait
}

function update(state: TreeState, info: NodeUpdateInfo): void {

    const _update = (action: string[], cur: NodeNotifyInfo, up: NodeNotifyInfo): void => {
        let posapply = false
        for (var ty of action) {
            switch (ty) {
                case UpdateType.UpdateAll:
                    _update_all(cur, up)
                    break
                case UpdateType.UpdateAlias:
                    cur.alias = up.alias
                    break
                case UpdateType.UpdateCode:
                    cur.code = up.code
                    break
                case UpdateType.UpdatePosition:
                    cur.pos = up.pos
                    posapply = true
                    break
                case UpdateType.UpdateLoop:
                    cur.loop = up.loop
                    break
                case UpdateType.UpdateWait:
                    cur.wait = up.wait
                    break
            }
        }

        if (!posapply) {
            message.success("update node " + cur.id + " success")
        }
    }

    for (var node of state.nodes) {
        if (node.id === info.info.id) {
            _update(info.type, node, info.info)
        }

        _find(info.info.id, node, (parent: NodeNotifyInfo, target: NodeNotifyInfo) => {
            _update(info.type, target, info.info)
        })
    }

}

function _find(id: string, parent: NodeNotifyInfo, callback: (parent: NodeNotifyInfo, target: NodeNotifyInfo, idx: number) => void) {

    if (parent.children && parent.children.length) {
        for (var i = 0; i < parent.children.length; i++) {
            if (parent.children[i].id === id) {
                callback(parent, parent.children[i], i)
                break
            }

            _find(id, parent.children[i], callback)
        }
    }

}

export function find(nodes: NodeNotifyInfo[], id: string): NodeNotifyInfo {
    let nod = getDefaultNodeNotifyInfo()

    for (var i = 0; i < nodes.length; i++) {
        if (nodes[i].id === id) {
            return nodes[i]
        }

        _find(id, nodes[i], (parent: NodeNotifyInfo, target: NodeNotifyInfo) => {
            nod = target
        })
    }

    return nod
}

function save(state: TreeState, name: string) {

    for (var nod of state.nodes) {
        if (nod.ty !== NodeTy.Root) {
            continue
        }

        var xmltree = {
            behavior: nod,
        };

        var blob = new Blob([OBJ2XML(xmltree)], {
            type: "application/json",
        });

        PostBlob(
            localStorage.remoteAddr,
            Api.FileBlobUpload,
            name,
            blob
        ).then((json: any) => {
            if (json.Code !== 200) {
                message.error(
                    "upload fail:" + String(json.Code) + " msg: " + json.Msg
                );
            } else {
                console.info(json.Body)
                message.success("upload succ ");
            }
        });
    }
}

const treeSlice = createSlice({
    name: "tree",
    initialState,
    reducers: {
        nodeAdd(state, action: PayloadAction<NodeAddInfo>) {
            let info = action.payload
            add(state, info)
        },
        nodeRmv(state, action: PayloadAction<string>) {
            rmv(state, action.payload)
        },
        nodeLink(state, action: PayloadAction<NodeLinkInfo>) {
            let info = action.payload
            link(state, info.parentid, info.childid)
        },
        nodeUnlink(state, action: PayloadAction<NodeUnlinkInfo>) {
            let info = action.payload
            unlink(state, info.targetid)
        },
        nodeUpdate(state, action: PayloadAction<NodeUpdateInfo>) {
            let info = action.payload
            update(state, info)
        },
        nodeClick(state, action: PayloadAction<NodeClickInfo>) {
            state.currentClickNode = action.payload
        },
        nodeRedraw(state, action: PayloadAction<void>) {
            state.updatetick++
        },
        setCurrentDebugBot(state, action: PayloadAction<string>) {
            let botid = action.payload
            state.currentDebugBot = botid
        },
        initTree(state, action: PayloadAction<NodeNotifyInfo>) {
            let tree = action.payload
            if (tree === null || tree === undefined) {
                return
            }
            if (tree.ty !== NodeTy.Root) {
                console.warn("tree parent node is not root")
                return
            }

            state.currentTreeName = ""
            state.rootid = tree.id

            state.history.splice(0, state.history.length)
            state.nodes = [tree]

            console.info("init tree", state.nodes)
        },
        cleanTree(state, action: PayloadAction<void>) {
            console.info("clean tree")

            state.currentTreeName = ""
            state.rootid = ""
            state.history.splice(0, state.history.length)
            state.nodes.splice(0, state.nodes.length)

            state.updatetick++
        },
        nodeSave(state, action: PayloadAction<string>) {
            let behavirName = action.payload
            save(state, behavirName)
        },
        unlockFocus(state, action: PayloadAction<NodeClickInfo>) {
            let info = action.payload
            state.currentLockedNode = info
        }
    },
});

export const { nodeAdd, nodeRmv, nodeLink, nodeUnlink, cleanTree, nodeUpdate, nodeClick, nodeRedraw, initTree, setCurrentDebugBot, nodeSave,unlockFocus } = treeSlice.actions;
export default treeSlice;
