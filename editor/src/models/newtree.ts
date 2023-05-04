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
}


const initialState: TreeState = {
    nodes: new Array<NodeNotifyInfo>(),
    history: new Array<NodeNotifyInfo>(),
    rootid: "",
    currentTreeName: "",
    currentDebugBot: "",
    currentDebugTree: getDefaultNodeNotifyInfo(),
    currentClickNode: { id: "", type: "" },
};

function add(state: TreeState, info: NodeAddInfo) {

}

function rmv(state: TreeState, id: string) {

}

function link(state: TreeState, parentid: string, childrenid: string) {

}

// targetid 被断开连接的节点
function unlink(state: TreeState, targetid: string) {

}

function update(state: TreeState, info: NodeNotifyInfo) {

    // 这里只会更新节点的属性
    let _update = (cur: NodeNotifyInfo, up: NodeNotifyInfo): void => {
        cur.code = up.code
        cur.alias = up.alias
        cur.pos = up.pos
        cur.loop = up.loop
        cur.wait = up.wait
    }

    for (var i = 0; i < state.nodes.length; i++) {
        if (state.nodes[i].id === info.id) {
            _update(state.nodes[i], info)
        }

        _find(info.id, state.nodes[i], (parent: NodeNotifyInfo, target: NodeNotifyInfo) => {
            _update(target, info)
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

const treeSlice = createSlice({
    name: "tree",
    initialState,
    reducers: {
        nodeAdd(state, action: PayloadAction<NodeAddInfo>) {
            console.info("node add", action.payload.info.id)
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
        nodeUpdate(state, action: PayloadAction<NodeNotifyInfo>) {
            let info = action.payload
            update(state, info)
        },
        nodeClick(state, action: PayloadAction<NodeClickInfo>) {
            state.currentClickNode = action.payload
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
        },
    },
});

export const { nodeAdd, nodeLink, nodeUnlink, cleanTree, nodeUpdate, nodeClick, initTree } = treeSlice.actions;
export default treeSlice;
