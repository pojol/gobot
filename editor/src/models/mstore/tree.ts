import { NodeTy } from "@/constant/node_type";
import { createAction, createSlice, PayloadAction } from "@reduxjs/toolkit";
const { Post, PostBlob, PostGetBlob } = require("../../utils/request");
import Cmd from "@/constant/cmd";
import { message } from "antd";
import Api from "@/constant/api";
import Topic from "@/constant/topic";
import PubSub from "pubsub-js";
import { useDispatch } from "react-redux";

interface TreeState {
    nodes: Array<NodeNotifyInfo>;
    history: Array<any>;
    rootid: string;
    treeState: boolean;
    currentTreeName: string;
    currentDebugTree: NodeNotifyInfo;
    currentDebugBot: string;
}

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
    console.info("sync map info", nod.id, nod)

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

    _fillData(root, window.tree.get(root.id), true, false);
    if (root && root.children.length) {
        _foreachRelation(root);
    }

    return root;
}

function CreateDebugBot(state: TreeState) {
    state.currentDebugTree = _getTree(state);
}

function BotStep(state: TreeState) {
    if (state.currentDebugBot === "") {
        message.warning("have not created bot");
        return;
    }

    var botid = state.currentDebugBot

    Post(localStorage.remoteAddr, Api.DebugStep, { BotID: botid }).then(
        (json: any) => {

            if (json.Code !== 200) {
                if (json.Code === 1009) {
                    message.warning(json.Code.toString() + " " + json.Msg)
                    return;
                } else if (json.Code === 1007) {
                    message.success("the end");
                } else {
                    message.warning(json.Code.toString() + " " + json.Msg)
                }
            }

            PubSub.publish(Topic.DebugFocus, []);  // reset focus
            console.info("step", json.Code, json)

            // 推送 reponse 面板信息
            let threadinfo = JSON.parse(json.Body.ThreadInfo) as Array<ThreadInfo>
            PubSub.publish(Topic.DebugUpdateChange, threadinfo)

            // 推送当前节点信息
            let focusLst = new Array<string>
            threadinfo.forEach(element => {
                focusLst.push(element.curnod)
            });
            PubSub.publish(Topic.DebugFocus, focusLst)

            // 推送 meta 面板信息
            let metaStr = JSON.stringify(JSON.parse(json.Body.Blackboard))
            PubSub.publish(Topic.DebugUpdateBlackboard, metaStr);
        }
    );
}

const treeSlice = createSlice({
    name: "tree",
    initialState,
    reducers: {
        nodeAdd(state, action: PayloadAction<NodeAddInfo>) {
            let info = action.payload
            if (info.build) {
                Add(state, info.info, info.silent)
            }
        },
        nodeLink(state, action: PayloadAction<NodeLinkInfo>) {
            let info = action.payload
            Link(state, info.parentid, info.childid, info.silent)
        },
        nodeUnlink(state, action: PayloadAction<NodeUnlinkInfo>) {
            let info = action.payload
            Unlink(state, info.targetid, info.silent)
        },
        cleanTree(state, action: PayloadAction<void>) {
            window.tree = new Map();

            state.currentTreeName = ""
            state.rootid = ""
            state.history.splice(0, state.history.length)
            state.nodes.splice(0, state.nodes.length)
        },
        setCurrentDebugBot(state, action: PayloadAction<string>) {
            let botid = action.payload
            state.currentDebugBot = botid
        },
        debug(state, action: PayloadAction<(tree: NodeNotifyInfo) => void>) {
            let callback = action.payload
            CreateDebugBot(state);
            callback(state.currentDebugTree)
        },
        step(state, action: PayloadAction<void>) {
            BotStep(state)
        }
    },
});

export const { nodeAdd, nodeLink, nodeUnlink, cleanTree, debug, step, setCurrentDebugBot } = treeSlice.actions;
export default treeSlice;
