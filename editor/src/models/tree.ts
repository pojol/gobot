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

function _getChildrenRelationInfo(
    parentChildren: Array<NodeNotifyInfo>,
    children: NodeNotifyInfo
) {
    let cinfo: NodeNotifyInfo = getDefaultNodeNotifyInfo();
    cinfo.id = children.id;

    parentChildren.push(cinfo);

    if (children.children && children.children.length) {
        children.children.forEach((cc) => {
            _getChildrenRelationInfo(cinfo.children, cc);
        });
    }
}

function _getRelationInfo(nod: NodeNotifyInfo) {
    var rinfo = getDefaultNodeNotifyInfo()
    rinfo.id = nod.id

    if (nod.children && nod.children.length) {
        nod.children.forEach((children) => {
            _getChildrenRelationInfo(rinfo.children, children);
        });
    }

    return rinfo;
}


function _syncMapInfo(nod: NodeNotifyInfo) {
    window.tree.set(nod.id, nod);

    if (nod.children && nod.children.length) {
        nod.children.forEach((children) => {
            _syncMapInfo(children);
        });
    }
}

function _findNode(id: string, parent: NodeNotifyInfo, callback: any) {
    if (parent.children && parent.children.length) {
        for (var i = 0; i < parent.children.length; i++) {
            if (parent.children[i].id === id) {
                callback(parent, parent.children[i], i);
                break;
            }

            _findNode(id, parent.children[i], callback);
        }
    }
};


function Add(state: TreeState, nod: NodeNotifyInfo, silent: boolean) {
    if (nod.ty == NodeTy.Root) {
        state.rootid = nod.id
    }

    let rinfo = _getRelationInfo(nod)
    _syncMapInfo(nod)

    console.info("add tar", JSON.stringify(rinfo))
    let olst = state.nodes
    olst.push(rinfo)

    let ohistory = state.history
    if (!silent) {
        let cmd = [{ cmd: Cmd.RMV, parm: nod.id }]
        ohistory.push(cmd)
    }

    state.nodes = olst
    state.history = ohistory
}


function Link(state: TreeState, parentid: string, childid: string, silent: boolean) {
    let children = getDefaultNodeNotifyInfo();
    let onods = state.nodes
    let ohistory = state.history

    let findSplice = (id: string, nod: NodeNotifyInfo) => {
        _findNode(
            id,
            nod,
            (
                parent: NodeNotifyInfo,
                innerChildren: NodeNotifyInfo,
                idx: number
            ) => {
                if (!silent) {
                    let cmd = [{ cmd: Cmd.Unlink, parm: [innerChildren.id] }];
                    console.info("history push", cmd);
                    ohistory.push(cmd);
                }

                parent.children.splice(idx, 1);
                children = innerChildren;
            }
        );
    };

    let findPush = (id: string, nod: NodeNotifyInfo) => {
        _findNode(id, nod, (_id: string, parent: NodeNotifyInfo) => {
            parent.children.push(children);
        });
    };

    for (let i = 0; i < onods.length; i++) {
        if (onods[i].id === childid) {
            children = onods[i];
            onods.splice(i, 1);
            break;
        }

        findSplice(childid, onods[i]);
    }

    if (children.id !== "") {
        for (let i = 0; i < onods.length; i++) {
            if (onods[i].id === parentid) {
                onods[i].children.push(children);
                break;
            }

            findPush(parentid, onods[i]);
        }
    }

    state.nodes = onods
};


function Unlink(state: TreeState, childid: string, silent: boolean) {
    let onods = state.nodes
    let ohistory = state.history
    let children;

    let find = (id: string, nod: NodeNotifyInfo) => {
        _findNode(
            id,
            nod,
            (
                innerParent: NodeNotifyInfo,
                innerChildren: NodeNotifyInfo,
                idx: number
            ) => {
                if (!silent) {
                    let cmd = [
                        { cmd: Cmd.Link, parm: [innerParent.id, innerChildren.id] },
                    ];
                    console.info("history push", cmd);
                    ohistory.push(cmd);
                }

                innerParent.children.splice(idx, 1);
                children = innerChildren;
            }
        );
    };

    for (var i = 0; i < onods.length; i++) {
        if (onods[i].id === childid) {
            children = onods[i];
            onods.splice(i, 1);
            break;
        }

        find(childid, onods[i]);
    }

    if (children) {
        onods.push(children);
    }

    state.nodes = onods
};

function _fillData(
    org: NodeNotifyInfo,
    info: NodeNotifyInfo,
    graph: boolean,
    edit: boolean
) {
    if (graph) {
        org.pos = info.pos;
    }

    if (edit) {
        switch (info.ty) {
            case NodeTy.Condition:
                org.code = info.code;
                break;
            case NodeTy.Loop:
                console.info(org.loop, "set loop", info.loop);
                org.loop = info.loop;
                break;
            case NodeTy.Wait:
                org.wait = info.wait;
                break;
            default:
                org.code = info.code;
                org.alias = info.alias;
        }
    }

    org.id = info.id
    org.ty = info.ty;
}


function _foreachRelation(parent: NodeNotifyInfo) {
    for (var i = 0; i < parent.children.length; i++) {
        if (window.tree.has(parent.children[i].id)) {
            _fillData(
                parent.children[i],
                window.tree.get(parent.children[i].id),
                true,
                true
            );
        }

        if (parent.children[i].children && parent.children[i].children.length) {
            _foreachRelation(parent.children[i]);
        }
    }
}

function _getTree(state: TreeState): NodeNotifyInfo {
    let root!: NodeNotifyInfo
    for (var i = 0; i < state.nodes.length; i++) {
        if (state.nodes[i].id === state.rootid) {
            root = state.nodes[i];
            break;
        }
    }

    console.info("root:", root, "tree", window.tree)
    if (root === undefined) {
        return getDefaultNodeNotifyInfo()
    }

    _fillData(root, window.tree.get(root.id), true, false);
    if (root && root.children.length) {
        _foreachRelation(root);
    }

    return root;
}

function CreateDebugBot(state: TreeState) {
    state.currentDebugTree = _getTree(state);
}

function UpdateEditInfo(state: TreeState, info: NodeNotifyInfo) {
    let tnode = window.tree.get(info.id);
    if (tnode === undefined) {
        tnode = getDefaultNodeNotifyInfo()
    }
    _fillData(tnode, info, false, true);

    if (info.notify) {
        message.success("apply info succ");
    }

    window.tree.set(info.id, tnode);
}

function SaveBehavior(state: TreeState, name: string) {
    let root = _getTree(state);
    var xmltree = {
        behavior: root,
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

const treeSlice = createSlice({
    name: "tree",
    initialState,
    reducers: {

        nodeLink(state, action: PayloadAction<NodeLinkInfo>) {
            let info = action.payload
            Link(state, info.parentid, info.childid, info.silent)
        },
        nodeUnlink(state, action: PayloadAction<NodeUnlinkInfo>) {
            let info = action.payload
            Unlink(state, info.targetid, info.silent)
        },

        save(state, action: PayloadAction<string>) {
            let behavirName = action.payload
            SaveBehavior(state, behavirName)
        },
        debug(state, action: PayloadAction<(tree: NodeNotifyInfo) => void>) {
            let callback = action.payload
            CreateDebugBot(state);
            callback(state.currentDebugTree)
        },
    },
});

export const { nodeLink, nodeUnlink, debug, save } = treeSlice.actions;
export default treeSlice;
