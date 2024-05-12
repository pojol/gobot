import { Cell, Graph, Node, Shape } from "@antv/x6";
import React, { useEffect, useState, useRef } from 'react';

/// <reference path="graph.d.ts" />
import OBJ2XML from "object-to-xml";

import { RootState } from "@/models/store";
import PubSub from "pubsub-js";
import Topic from "../../constant/topic";
import { TaskTimer } from 'tasktimer';


import { setDebugInfo, setLock } from "@/models/debuginfo";
import { history } from "umi";

import {
    AimOutlined,
    BugOutlined,
    StepForwardOutlined,
    ClearOutlined,
    FastForwardOutlined,
    CloudSyncOutlined,
    DeleteOutlined,
    UndoOutlined,
    ZoomInOutlined,
    ZoomOutOutlined,
    LockOutlined,
    UnlockOutlined,
    PushpinOutlined
} from "@ant-design/icons";
import { Button, Input, Modal, Tooltip, theme } from "antd";
import { IsActionNode, IsScriptNode, NodeTy } from "../../constant/node_type";

import { message } from "antd";
import { ConnectedProps, connect, useSelector } from "react-redux";
import { useLocation } from 'react-router-dom';
import "./graph.css";


import Api from "@/constant/api";
import ThemeType from "@/constant/constant";
import { GetNode } from "./shape/shape";
import { EditSidePlane } from "./side";
import { CreateGraph } from "./canvas/canvas";
import { GetNodInfo, getDefaultNodeNotifyInfo } from "@/models/node";
import {
    nodeRedraw, nodeRmv, cleanTree, nodeAdd, initTree, nodeClick, nodeLink,
    nodeUnlink,
    setCurrentDebugBot,
    nodeUpdate,
    UpdateType,
    nodeSave,
    unlockFocus,
    nodeUndo,
} from "@/models/tree";
import { SelectorDarkNode } from "./shape/selector";
import SequenceLightNode, { SequenceDarkNode } from "./shape/sequence";

const {
    LoadBehaviorWithBlob,
    LoadBehaviorWithFile,
} = require("../../utils/parse");
const { PostBlob, Post } = require("../../utils/request");


function iterate(nod: Node, callback: (nod: Node) => void) {
    if (nod !== null && nod !== undefined) {
        callback(nod);

        nod.eachChild((children, idx) => {
            iterate(children as Node, callback);
        });
    }
}

const mapStateToProps = (state: RootState) => ({
    prefabMap: state.prefabSlice.pmap,
    tree: state.treeSlice,
    debugInfo: state.debugInfoSlice,
});

const connector = connect(mapStateToProps);
type PropsFromRedux = ConnectedProps<typeof connector>;
interface GraphViewProps extends PropsFromRedux { }

const GraphView = (props: GraphViewProps) => {

    const graphRef = React.useRef<Graph>();
    const containerRef = React.useRef<HTMLDivElement>(null);

    const [BpNodes, setBpNodes] = useState<BPNodeLinkInfo[]>();

    const [modalVisible, setModalVisible] = useState<boolean>(false);
    const [behaviorName, setBehaviorName] = useState<string>("");
    const [wflex, setWflex] = useState<number>(0.6);

    const { graphFlex } = useSelector((state: RootState) => state.resizeSlice)
    const { lock } = useSelector((state: RootState) => state.debugInfoSlice)
    const { currentTreeName, nodes, updatetick, currentClickNode, currentDebugBot } = useSelector((state: RootState) => state.treeSlice)
    const [isGraphCreated, setIsGraphCreated] = useState(false);

    const [timer, setTimer] = useState<TaskTimer | null>(null);
    // 使用 useRef 创建 ref 对象
    const timerRef = useRef<TaskTimer | null>(null);

    useEffect(() => {
        console.info("create graph")
        const graph = CreateGraph(containerRef.current, wflex, graphFlex)

        graph.bindKey("del", () => {
            ClickDel();
            return false;
        });

        graph.bindKey("ctrl+z", () => {
            ClickUndo();
        });

        graph.bindKey(["f10", "command+f10", "ctrl+f10"], () => {
            ClickStep(1);
        });

        graph.bindKey(["f11", "command+f11", "ctrl+f11"], () => {
            ClickCreateDebug(1);
        });

        graph.on("edge:removed", ({ edge, options }) => {
            if (!options.ui) {
                return;
            }

            _findCell(edge.getTargetCellId(), (child) => {
                props.dispatch(nodeUnlink({ targetid: child.id, silent: false }))
                props.dispatch(nodeRedraw())
                //child.getParent()?.removeChild(edge);
            });
        });

        graph.on("edge:connected", ({ isNew, edge }) => {
            const source = edge.getSourceNode();
            const target = edge.getTargetNode();

            if (isNew) {
                if (source !== null && target !== null) {
                    const typename = source.getAttrs().type.name?.toString()
                    if (typename === undefined) {
                        message.warning("Cannot get node name");
                        return
                    }

                    if (target.getAttrs().type.name === NodeTy.Root) {
                        message.warning("Cannot connect to root node");
                        graph.removeEdge(edge.id, { disconnectEdges: true });
                        return;
                    }

                    if (IsScriptNode(typename) && source.getChildCount() > 0) {
                        message.warning("Action node can only mount a single node");
                        graph.removeEdge(edge.id, { disconnectEdges: true });
                        return;
                    }

                    if (target.parent !== undefined && target.parent != null) {
                        message.warning("Cannot connect to a node that has a parent node");
                        graph.removeEdge(edge.id, { disconnectEdges: true });
                        return;
                    }

                    edge.setZIndex(0);
                    source.addChild(target);
                    props.dispatch(nodeLink({ parentid: source.id, childid: target.id, silent: false }))
                }
            }
        });

        graph.on("node:click", ({ node }) => {
            props.dispatch(nodeClick({ id: node.id, type: node.getAttrs().type.name as string }))
        });

        graph.on("blank:click", () => {
            props.dispatch(nodeClick({ id: "", type: "" }))

            var nods = graph.getRootNodes();
            if (nods.length > 0) {
                iterate(nods[0], (nod) => {
                    if (nod.getAttrs().type !== undefined) {
                        nod.setPortProp(
                            nod.getPorts()[0].id as string,
                            "attrs/portBody/r",
                            4
                        );
                    }
                });
            }
        })

        graph.on("node:added", ({ node, index, options }) => {
            let build = false
            if (options.others !== undefined) {
                build = options.others.build;
            }

            if (!build) {
                props.dispatch(nodeAdd({
                    info: GetNodInfo(props.prefabMap, node, "", ""),
                    silent: false,
                }))
            }
        });

        graph.on("node:mouseenter", ({ node }) => {

            var ty = node.getAttrs().type.name as string
            if (ty !== NodeTy.BreakPoint) {
                node.setPortProp(node.getPorts()[0].id as string, "attrs/portBody/r", 8);
            }

        });

        graph.on("node:moved", ({ e, x, y, node, view: NodeView }) => {
            iterate(node, (nod) => {
                if (nod.isNode()) {
                    var newinfo = getDefaultNodeNotifyInfo()
                    newinfo.id = nod.id
                    newinfo.pos = {
                        x: nod.position().x,
                        y: nod.position().y,
                    }

                    props.dispatch(nodeUpdate({
                        info: newinfo,
                        type: [UpdateType.UpdatePosition]
                    }))
                }
            });
        });

        graph.on("edge:mouseenter", ({ edge }) => {
            edge.addTools([
                "source-arrowhead",
                "target-arrowhead",
                {
                    name: "button-remove",
                    args: {
                        distance: -30,
                        onClick({ e, cell, view }: any) {
                            //var sourcenod = cell.getSourceNode();
                            var targetnod = cell.getTargetNode();
                            //
                            //this.graph.removeEdge(cell.id, { disconnectEdges: true });
                            props.dispatch(
                                nodeUnlink({
                                    targetid: targetnod.id,
                                    silent: false,
                                })
                            );
                            props.dispatch(nodeRedraw())
                            //sourcenod.unembed(targetnod);
                        },
                    },
                },
            ]);
        });

        graph.on("edge:mouseleave", ({ edge }) => {
            edge.removeTools();
        });

        PubSub.subscribe(Topic.WindowResize, () => {
            resizeViewpoint(wflex);
        });

        PubSub.subscribe(Topic.ThemeChange, (topic: string, theme: string) => {
            var nods = graph.getRootNodes();
            if (nods.length > 0) {
                iterate(nods[0], (nod) => {
                    let bodyfill = "#f5f5f5";
                    let labelfill = "#20262E";
                    let portfill = "#fff"

                    if (theme === ThemeType.Dark) {
                        bodyfill = "#20262E";
                        labelfill = "#fff";
                        portfill = "#20262E"
                    }

                    nod.setAttrs({
                        body: {
                            fill: bodyfill,
                        },
                        label: {
                            fill: labelfill,
                        },
                    });

                    if (nod.isEdge() && (nod.parent instanceof SequenceDarkNode || nod.parent instanceof SequenceLightNode)) {
                        nod.setLabelAt(0, {
                            attrs: {
                                text: {
                                    text: nod.getLabelAt(0).attrs.text.text,
                                },
                                body: {
                                    fill: bodyfill,
                                },
                                label: {
                                    fill: labelfill,
                                },
                            },
                        })
                    }

                    if (nod.isNode()) {
                        nod.setPortProp(nod.getPortAt(0).id as string, "attrs/portBody/fill", portfill)
                    }
                });
            }
        });

        graphRef.current = graph;
        console.info("graph init done", graphRef.current)

        if (history.location.pathname.length > 8) { // "/editor"
            let botname = history.location.pathname.slice(8);

            if (botname != null && botname != "") {
                let chineseChar = decodeURIComponent(botname);
                console.info("load bot", botname, " => ", chineseChar);
                LoadBehaviorWithBlob(
                    localStorage.remoteAddr,
                    Api.FileGet,
                    chineseChar
                ).then((file: any) => {
                    //props.dispatch(cleanTree());
                    LoadBehaviorWithFile(chineseChar, file.blob, (tree: any) => {
                        graph.clearCells();
                        props.dispatch(initTree(tree))
                        redraw(tree, true);
                    });
                });
            }
        }

        setIsGraphCreated(true);

        return () => {
            setBpNodes(prevBpNodes => { return [] })
            graph.dispose()
        }
    }, []);

    useEffect(() => {

        redrawTree()

    }, [updatetick])

    useEffect(() => {

        resizeViewpoint(graphFlex)

    }, [graphFlex])

    useEffect(() => {
        const newTimer = new TaskTimer(200);
        newTimer.on('tick', () => {
            step(currentDebugBot, BpNodes);
        });
        setTimer(newTimer);

        // 将新创建的 timer 赋值给 ref 对象
        timerRef.current = newTimer;

        return () => {
            // 清理定时器
            if (timer) {
                timer.stop();
            }
            setTimer(null);
            // 在组件卸载时将 ref 对象置为 null
            timerRef.current = null;
        }
    }, [currentDebugBot, BpNodes]);

    useEffect(() => {

        if (graphRef.current) {
            graphRef.current.bindKey(["up"], () => {
                if (currentClickNode.id !== "") {
                    _findCell(currentClickNode.id, (cell) => {
                        var nod = (cell as Node<Node.Properties>)
                        nod.setPosition({
                            x: nod.position().x,
                            y: nod.position().y - 1,
                        })
                    });
                }
            });
            graphRef.current.bindKey(["down"], () => {
                if (currentClickNode.id !== "") {
                    _findCell(currentClickNode.id, (cell) => {
                        var nod = (cell as Node<Node.Properties>)
                        nod.setPosition({
                            x: nod.position().x,
                            y: nod.position().y + 1,
                        })
                    });
                }
            });
            graphRef.current.bindKey(["left"], () => {
                if (currentClickNode.id !== "") {
                    _findCell(currentClickNode.id, (cell) => {
                        var nod = (cell as Node<Node.Properties>)
                        nod.setPosition({
                            x: nod.position().x - 1,
                            y: nod.position().y,
                        })
                    });
                }
            });
            graphRef.current.bindKey(["right"], () => {
                if (currentClickNode.id !== "") {
                    _findCell(currentClickNode.id, (cell) => {
                        var nod = (cell as Node<Node.Properties>)
                        nod.setPosition({
                            x: nod.position().x + 1,
                            y: nod.position().y,
                        })
                    });
                }
            });

            graphRef.current.bindKey(["f9", "command+f9", "ctrl+f9"], () => {
                if (currentClickNode.id !== "") {
                    _findCell(currentClickNode.id, (cell) => {
                        var node = (cell as Node<Node.Properties>)
                        ClickBreakpointImpl(node)
                    })
                }
            })
        }

        return () => {
            if (graphRef.current) {
                graphRef.current.unbindKey(["up"])
                graphRef.current.unbindKey(["down"])
                graphRef.current.unbindKey(["left"])
                graphRef.current.unbindKey(["right"])
                graphRef.current.unbindKey(["f9", "command+f9", "ctrl+f9"])
            }
        }

    }, [currentClickNode])

    const redrawTree = () => {

        if (graphRef.current) {
            graphRef.current.clearCells();

            console.info("redraw tree nods", nodes)

            if (nodes.length > 0) {
                for (var i = 0; i < nodes.length; i++) {
                    redraw(nodes[i], true)
                }
            } else {
                var root = GetNode(NodeTy.Root, {});
                root.setPosition(
                    graphRef.current.getGraphArea().width / 2,
                    graphRef.current.getGraphArea().height / 2 - 200
                );
                graphRef.current.addNode(root);
                props.dispatch(initTree(
                    GetNodInfo(props.prefabMap, root, "", "")
                ))
            }

        }
    }

    const getLoopLabel = (val: Number) => {
        var tlab = "";
        if (val !== 0) {
            tlab = val.toString() + " times";
        } else {
            tlab = "endless";
        }

        return tlab;
    }

    // 重绘视口
    const resizeViewpoint = (graphFlex: number) => {
        var width = document.documentElement.clientWidth * graphFlex;

        console.info("resize panel", graphFlex, document.documentElement.clientWidth);

        // 设置视口大小
        if (graphRef.current) {
            graphRef.current.resize(width, document.documentElement.clientHeight - 62);
            /*
                        var root = GetNode(NodeTy.Root, {});
                        if (root !== null && root !== undefined) {
                            console.info("root set position", graphRef.current.getGraphArea().width / 2, graphRef.current.getGraphArea().height / 2 - 200);
                            root.setPosition(
                                graphRef.current.getGraphArea().width / 2,
                                graphRef.current.getGraphArea().height / 2 - 200
                            );
                        }
                        */
        }
    }

    const _findCellChild = (parent: Cell, id: String, callback: (nod: Cell) => void) => {
        if (parent.id === id) {
            callback(parent);
            return;
        } else {
            parent.eachChild((child, idx) => {
                _findCellChild(child, id, callback);
            });
        }
    };

    const _findCell = (id: String, callback: (nod: Cell) => void) => {
        let nods: Node<Node.Properties>[]

        if (graphRef.current) {
            nods = graphRef.current.getRootNodes();
            if (nods.length >= 0) {
                if (nods[0].id === id) {
                    callback(nods[0]);
                } else {
                    nods[0].eachChild((child, idx) => {
                        _findCellChild(child, id, callback);
                    });
                }
            }
        }
    };

    const redrawChild = (parent: any, parentty: string, child: any, build: boolean, idx: number) => {
        var nod: Node;

        if (graphRef.current == null) {
            return
        }

        let others = { build: build, silent: true, code: "", alias: "" }

        switch (child.ty) {
            case NodeTy.Selector:
                nod = GetNode(child.ty, { id: child.id });
                nod.setAttrs({ label: { text: "sel" } })
                break;
            case NodeTy.Sequence:
                nod = GetNode(child.ty, { id: child.id });
                nod.setAttrs({ label: { text: "seq" } })
                break;
            case NodeTy.Parallel:
                nod = GetNode(child.ty, { id: child.id });
                nod.setAttrs({ label: { text: "par" } })
                break;
            case NodeTy.Loop:
                nod = GetNode(child.ty, { id: child.id });
                nod.setAttrs({ label: { text: child.loop.toString() + " times" } })
                break;
            case NodeTy.Wait:
                nod = GetNode(child.ty, { id: child.id });
                nod.setAttrs({ label: { text: child.wait.toString() + " ms" } })
                break;
            case NodeTy.Condition:
                nod = GetNode(child.ty, { id: child.id });
                others.code = child.code
                break;
            default:
                nod = GetNode(child.ty, { id: child.id });
                others.code = child.code
                if (child.alias === "") {
                    others.alias = child.ty
                } else {
                    others.alias = child.alias
                }
                nod.setAttrs({ type: { name: child.ty } });
        }

        nod.setPosition({
            x: child.pos.x,
            y: child.pos.y,
        });

        graphRef.current.addNode(nod, { others: others });

        if (parent) {
            var edge = new Shape.Edge({
                attrs: {
                    line: {
                        stroke: "#a0a0a0",
                        strokeWidth: 1,
                        targetMarker: {
                            name: "classic",
                            size: 3,
                        },
                    },
                },
                zIndex: 0,
                source: parent,
                target: nod,
            })

            let bodyfill = "#20262E"
            let labelfill = "#fff"
            if (localStorage.theme === ThemeType.Light) {
                bodyfill = "#f5f5f5"
                labelfill = "#20262E"
            }

            if (parentty === NodeTy.Sequence) {
                edge.appendLabel({
                    attrs: {
                        text: {
                            text: (idx + 1).toString(),
                        },
                        body: {
                            fill: bodyfill,
                        },
                        label: {
                            fill: labelfill,
                        },
                    },
                })
            }

            graphRef.current.addEdge(edge);
            parent.addChild(nod);
        }

        if (IsScriptNode(child.ty)) {
            nod.setAttrs({ label: { text: child.alias } });
        } else if (child.ty === NodeTy.Loop) {
            nod.setAttrs({ label: { text: getLoopLabel(child.loop) } });
        } else if (child.ty === NodeTy.Wait) {
            nod.setAttrs({ label: { text: child.wait.toString() + " ms" } });
        } else if (child.ty === NodeTy.Sequence) {
            nod.setAttrs({ label: { text: "seq" } });
        } else if (child.ty === NodeTy.Selector) {
            nod.setAttrs({ label: { text: "sel" } });
        } else if (child.ty === NodeTy.Parallel) {
            nod.setAttrs({ label: { text: "par" } });
        }

        if (child.children && child.children.length) {
            for (var i = 0; i < child.children.length; i++) {
                redrawChild(nod, child.ty, child.children[i], build, i);
            }
        }
    }

    const redraw = (jsontree: NodeNotifyInfo, build: boolean) => {

        if (jsontree.ty === NodeTy.Root) {
            var root = GetNode(NodeTy.Root, { id: jsontree.id });
            root.setPosition({
                x: jsontree.pos.x,
                y: jsontree.pos.y,
            });

            if (graphRef.current) {
                graphRef.current.addNode(root, { others: { build: build, silent: true } });
            }

            if (jsontree.children && jsontree.children.length) {
                for (var i = 0; i < jsontree.children.length; i++) {
                    redrawChild(root, "NodeTy.Root", jsontree.children[i], build, i);
                }
            }
        } else {
            redrawChild(null, "", jsontree, build, 0);
        }
    }

    // 加亮显示当前运行到的某个 Node
    const debugFocus = (info: Array<string>) => {
        console.info("debug focus", info)

        // clean
        cleanStepInfo();

        info.forEach(element => {
            _findCell(element, (nod) => {
                nod.transition(
                    "attrs/body/strokeWidth", "4px", {
                    //interp: Interp.unit,
                    timing: 'bounce', // Timing.bounce
                },
                )()
            });
        });

    }

    const removeCell = (cell: Cell) => {
        if (cell.getParent() == null) {
            if (graphRef.current) {
                graphRef.current.removeCell(cell);
            }
        } else {
            //PubSub.publish(Topic.NodeRmv, cell.id);
            cell.getParent()?.removeChild(cell);
        }
    }

    const ClickDel = () => {
        let cells: Cell<Cell.Properties>[]

        if (graphRef.current) {
            cells = graphRef.current.getSelectedCells()
            if (cells.length) {

                for (var i = 0; i < cells.length; i++) {
                    if (cells[i].isNode()) {
                        props.dispatch(nodeRmv(cells[i].id))
                    }
                    /*
                                        if (cells[i].getAttrs().type.toString() !== NodeTy.Root) {
                                            removeCell(cells[i]);
                                        }
                                        */
                }

                props.dispatch(nodeRedraw())
            }
        }

    };

    const CleanTree = () => {
        props.dispatch(cleanTree())
    }

    const ClickUpload = () => {

        if (currentTreeName === "") {
            setModalVisible(true);
        } else {
            props.dispatch(nodeSave(""))
        }

    };

    // 基于某个模版，创建一个可运行的 Bot
    const ClickCreateDebug = (e: any) => {
        cleanStepInfo();
        timer?.stop()

        props.dispatch(setDebugInfo({ metaInfo: "{}", threadInfo: [], lock: false }))
        for (var nod of nodes) {
            if (nod.ty !== NodeTy.Root) {
                continue
            }

            var xmltree = {
                behavior: nod,
            };

            var blob = new Blob([OBJ2XML(xmltree)], {
                type: "application/json",
            });

            PostBlob(localStorage.remoteAddr, Api.DebugCreate, name, blob).then(
                (json: any) => {
                    if (json.Code !== 200) {
                        message.error(
                            "create fail:" + String(json.Code) + " msg: " + json.Msg
                        );
                    } else {
                        props.dispatch(setCurrentDebugBot(json.Body.BotID))
                        message.success("create debug bot succ");
                    }
                }
            );
        }
    };

    const ClickZoomIn = () => {
        if (graphRef.current) {
            graphRef.current.zoomTo(graphRef.current.zoom() * 1.2);
        }
    };

    const ClickBreakpointImpl = (node: Node<Node.Properties> | null) => {

        // 只对脚本节点进行操作
        if (node !== null) {
            if (IsScriptNode(node.getAttrs().type.name as string)) {
                setBpNodes(prevBpNodes => {
                    const tmp = prevBpNodes || []; // 保证 tmp 不为 undefined
                    const ln = tmp.find((element) => element.parentid === node?.id);

                    console.info(prevBpNodes, "node dblclick", node?.id, ln);
                    if (ln !== undefined) { // 删除
                        graphRef.current?.removeNode(ln.bpid);
                        const updatedNodes = tmp.filter((element) => element.parentid !== node?.id);
                        console.info("remove bp node", updatedNodes);
                        return updatedNodes;
                    } else {
                        var bpnode = GetNode(NodeTy.BreakPoint, {});
                        bpnode.setPosition(node.position().x - (bpnode.getSize().width + 7),
                            node.position().y + (bpnode.getSize().height / 2));
                        graphRef.current?.addNode(bpnode, { others: { build: true } });  // 不发送addnode事件

                        const updatedNodes = [...tmp, { parentid: node?.id, bpid: bpnode.id }];
                        console.info("add bp node", updatedNodes);
                        return updatedNodes;
                    }
                });
            }
        }

    }

    const ClickBreakpoint = () => {

        if (currentClickNode.id !== "") {
            _findCell(currentClickNode.id, (cell) => {
                var node = (cell as Node<Node.Properties>)
                ClickBreakpointImpl(node)
            })
        }
        
    }

    const ClickZoomOut = () => {
        if (graphRef.current) {
            graphRef.current.zoomTo(graphRef.current.zoom() * 0.8);
        }
    };

    const ClickZoomReset = () => {
        if (graphRef.current) {
            graphRef.current.zoomTo(1);
        }
    };

    const ClickStep = (e: any) => {
        if (props.tree.currentDebugBot === "") {
            message.warning("have not created bot");
            return;
        }

        props.dispatch(setLock(true))
        step(props.tree.currentDebugBot, BpNodes);
    };

    const step = (botid: string, bpnodes: BPNodeLinkInfo[] | undefined) => {
        Post(localStorage.remoteAddr, Api.DebugStep, { BotID: botid }).then(
            (json: any) => {
                if (json.Code !== 200) {
                    if (timerRef.current?.state === "running") {
                        timerRef.current?.stop();
                    }

                    if (json.Code === 1009) {
                        message.warning(json.Code.toString() + " " + json.Msg)
                        return;
                    } else if (json.Code === 1007) {
                        message.success("the end");
                    } else {
                        message.warning(json.Code.toString() + " " + json.Msg)
                    }
                }

                debugFocus([]); // clean

                // 推送 reponse 面板信息
                let threadinfo = JSON.parse(json.Body.ThreadInfo) as Array<ThreadInfo>

                // 推送当前节点信息
                let focusLst = new Array<string>
                threadinfo.forEach(element => {
                    const ln = bpnodes?.find((bp) => bp.parentid === element.curnod);
                    if (ln !== undefined) { // match bp node
                        if (timerRef.current?.state === "running") {
                            timerRef.current?.stop();
                        }
                    }

                    focusLst.push(element.curnod)
                });
                debugFocus(focusLst)

                // 推送 meta 面板信息
                let metaStr = JSON.stringify(JSON.parse(json.Body.Blackboard))
                props.dispatch(setDebugInfo({
                    metaInfo: metaStr,
                    threadInfo: threadinfo,
                    lock: false,
                }))


            }
        );

    }

    const ClickRunning = (e: any) => {
        if (props.tree.currentDebugBot === "") {
            message.warning("have not created bot");
            return;
        }

        if (timer?.state !== "running") {
            console.info("running this timer", timer)
            timer?.start()
        }
    }

    const behaviorNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setBehaviorName(e.target.value);
    };

    const modalHandleOk = () => {
        setModalVisible(false);
        if (behaviorName !== "") {
            props.dispatch(nodeSave(behaviorName))
            setBehaviorName("");
        } else {
            message.warning("please enter the file name of the behavior tree");
        }
    };

    const modalHandleCancel = () => {
        setModalVisible(false);
    };

    const cleanStepInfo = () => {
        // clean
        let nods: Node<Node.Properties>[]

        if (graphRef.current) {
            nods = graphRef.current.getRootNodes();
            if (nods.length > 0) {
                iterate(nods[0], (nod) => {
                    nod.setAttrs({
                        body: {
                            strokeWidth: 1,
                        },
                    });
                });
            }
        }
    };

    const ClickUndo = () => {
        props.dispatch(nodeUndo())
    };

    return (
        <div className='app'>
            <EditSidePlane
                graph={graphRef.current}
                isGraphCreated={isGraphCreated}
            />

            <div className="app-content" ref={containerRef} />

            <div
                className={"app-zoom-win"}
                style={{ marginLeft: 7, whiteSpace: "nowrap" }}
            >
                <Tooltip placement="topLeft" title="Add a breakpoint to the selected node [F9]">
                    <Button icon={<PushpinOutlined />} onClick={ClickBreakpoint} />
                </Tooltip>
                <Tooltip placement="topLeft" title="ZoomIn">
                    <Button icon={<ZoomInOutlined />} onClick={ClickZoomIn} />
                </Tooltip>
                <Tooltip placement="topLeft" title="Reset">
                    <Button icon={<AimOutlined />} onClick={ClickZoomReset} />
                </Tooltip>
                <Tooltip placement="topLeft" title="ZoomOut">
                    <Button icon={<ZoomOutOutlined />} onClick={ClickZoomOut} />
                </Tooltip>
                <Tooltip placement="topLeft" title="Undo [ ctrl+z ]">
                    <Button icon={<UndoOutlined />} onClick={ClickUndo} />
                </Tooltip>
                <Tooltip placement="topLeft" title="Delete Node [ del ]">
                    <Button icon={<DeleteOutlined />} onClick={ClickDel} />
                </Tooltip>
                <Tooltip placement="topLeft" title="Clean">
                    <Button icon={<ClearOutlined />} onClick={CleanTree} />
                </Tooltip>
            </div>

            <div className={"app-debug-reset"}>
                <Tooltip
                    placement="topRight"
                    title={"Create or reset to starting point [F11]"}
                >
                    <Button
                        icon={<BugOutlined />}
                        style={{ width: 50 }}
                        onClick={ClickCreateDebug}
                    >
                        {" "}
                    </Button>
                </Tooltip>
            </div>
            <div className={"app-debug-run"}>
                <Tooltip placement="topRight" title={"Run to the end node [F10]"}>
                    <Button
                        style={{ width: 50 }}
                        icon={<FastForwardOutlined />}
                        onClick={ClickRunning}
                    >
                        { }
                    </Button>
                </Tooltip>
            </div>
            <div className={"app-debug-step"}>
                <Tooltip placement="topRight" title={"Run to the next node [F10]"}>
                    <Button
                        style={{ width: 50 }}
                        icon={<StepForwardOutlined />}
                        onClick={ClickStep}
                    >
                        { }
                    </Button>
                </Tooltip>
            </div>
            <div className={"app-debug-upload"}>
                <Tooltip placement="topRight" title={"Upload the bot to the server"}>
                    <Button
                        icon={<CloudSyncOutlined />}
                        style={{ width: 50 }}
                        onClick={ClickUpload}
                    ></Button>
                </Tooltip>
            </div>

            <Modal
                open={modalVisible}
                onOk={modalHandleOk}
                onCancel={modalHandleCancel}
            >
                <Input
                    placeholder="input behavior file name"
                    onChange={behaviorNameChange}
                />
            </Modal>
        </div>
    )
}

export default connector(GraphView);
