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
    treeState: boolean;
    currentTreeName: string;
    currentDebugTree: NodeNotifyInfo;
    currentDebugBot: string;
    currentClickNode: NodeClickInfo;
}

function treeStateInit(): boolean {
    window.tree = new Map(); // 主要维护的是 editor 节点编辑后的数据
    return true
}

const initialState: TreeState = {
    nodes: new Array<NodeNotifyInfo>(),
    history: new Array<NodeNotifyInfo>(),
    rootid: "",
    currentTreeName: "",
    treeState: treeStateInit(),
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
        nodeFind(state, action: PayloadAction<(id: string, node: NodeNotifyInfo) => void>) {

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

            window.tree = new Map();

            state.currentTreeName = ""
            state.rootid = tree.id

            state.history.splice(0, state.history.length)
            state.nodes = [tree]
        },
        cleanTree(state, action: PayloadAction<void>) {
            console.info("clean tree")
            window.tree = new Map();

            state.currentTreeName = ""
            state.rootid = ""
            state.history.splice(0, state.history.length)
            state.nodes.splice(0, state.nodes.length)
        },
    },
});

export const { nodeAdd, nodeLink, nodeUnlink, cleanTree, nodeUpdate, nodeClick, initTree } = treeSlice.actions;
export default treeSlice;
