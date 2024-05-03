
declare module 'react-medium-editor'

interface PrefabInfo {
    name: string,
    tags: string[],
    code: string,
}

interface Window {
    tree: Map<string, any>,
}

interface BPNodeLinkInfo {
    parentid: string,
    bpid: string,
}

interface NodeUpdateInfo {
    type: string[],
    info: NodeNotifyInfo
}

interface NodeNotifyInfo {
    id: string,
    ty: string,
    code: string,
    loop: number,
    wait: number,
    alias: string,
    pos: {
        x: number,
        y: number,
    },
    children: Array<NodeNotifyInfo>,
    notify: boolean,
}

interface NodeClickInfo {
    id: string,
    type: string,
}

interface NodeAddInfo {
    info: NodeNotifyInfo,
    silent: boolean,
}


interface NodeLinkInfo {
    parentid: string,
    childid: string,
    silent: boolean,
}

interface NodeUnlinkInfo {
    targetid: string,
    silent: boolean,
}

interface NodeFindCallback {
    id: string,
    callback: (NodeNotifyInfo) => void
}

interface PrefabTagInfo {
    value: string
}


interface ThreadInfo {
    number: number
    errmsg: string
    curnod: string
    change: string
}