import { Cell, Graph, Node, Shape } from "@antv/x6";
import React, { useEffect, useState } from 'react';

/// <reference path="graph.d.ts" />
import OBJ2XML from "object-to-xml";

import { RootState } from "@/models/store";
import PubSub from "pubsub-js";
import Topic from "../../constant/topic";
import { TaskTimer } from "tasktimer/lib/TaskTimer";


import { setDebugInfo, setLock } from "@/models/debuginfo";
import { history } from "umi";

import {
    AimOutlined,
    BugOutlined,
    CaretRightOutlined,
    ClearOutlined,
    CloudUploadOutlined,
    DeleteOutlined,
    UndoOutlined,
    ZoomInOutlined,
    ZoomOutOutlined,
    LockOutlined,
    UnlockOutlined 
} from "@ant-design/icons";
import { Button, Input, Modal, Tooltip } from "antd";
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
} from "@/models/tree";

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

    const [modalVisible, setModalVisible] = useState<boolean>(false);
    const [behaviorName, setBehaviorName] = useState<string>("");
    const [wflex, setWflex] = useState<number>(0.6);

    const { graphFlex } = useSelector((state: RootState) => state.resizeSlice)
    const { lock } = useSelector((state: RootState) => state.debugInfoSlice)
    const { nodes, updatetick, currentClickNode,currentLockedNode } = useSelector((state: RootState) => state.treeSlice)
    const [isGraphCreated, setIsGraphCreated] = useState(false);

    useEffect(() => {
        console.info("create graph")
        const graph = CreateGraph(containerRef.current, wflex, graphFlex)

        graph.bindKey("del", () => {
            ClickDel();
            return false;
        });

        graph.bindKey("ctrl+z", () => {
            //PubSub.publish(Topic.Undo, {});
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

            findNode(edge.getTargetCellId(), (child) => {
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
                    console.info("connect to", typename)
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
            node.setPortProp(node.getPorts()[0].id as string, "attrs/portBody/r", 8);
        });

        // node:mouseleave 消息容易获取不到，先每次获取到这个消息将所有节点都设置一下
        graph.on("node:mouseleave", ({ node }) => {
            node.setPortProp(node.getPorts()[0].id as string, "attrs/portBody/r", 4);

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

                    if (nod.isNode()) {
                        nod.setPortProp(nod.getPortAt(0).id as string, "attrs/portBody/fill", portfill)
                    }
                });
            }
        });

        graphRef.current = graph;
        console.info("graph init done", graphRef.current)

        console.info("load path", history.location.pathname)
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
        //containerRef.current = graph.container;
        //containerRef.current.appendChild(graph.container);
        //stencilContainerRef.current.appendChild(graph.container);
    }, []);

    useEffect(() => {

        redrawTree()

    }, [updatetick])

    useEffect(() => {
        resizeViewpoint(graphFlex)

        // redraw

    }, [graphFlex])

    const redrawTree = () => {

        if (graphRef.current) {
            graphRef.current.clearCells();

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
        }
    }


    const findChild = (parent: Cell, id: String, callback: (nod: Cell) => void) => {
        if (parent.id === id) {
            callback(parent);
            return;
        } else {
            parent.eachChild((child, idx) => {
                findChild(child, id, callback);
            });
        }
    };

    const findNode = (id: String, callback: (nod: Cell) => void) => {
        let nods: Node<Node.Properties>[]

        if (graphRef.current) {
            nods = graphRef.current.getRootNodes();
            if (nods.length >= 0) {
                if (nods[0].id === id) {
                    callback(nods[0]);
                } else {
                    nods[0].eachChild((child, idx) => {
                        findChild(child, id, callback);
                    });
                }
            }
        }
    };


    const redrawChild = (parent: any, child: any, build: boolean) => {
        var nod: Node;

        if (graphRef.current == null) {
            return
        }

        let others = { build: build, silent: true, code: "", alias: "" }
        console.info(child)

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
            graphRef.current.addEdge(
                new Shape.Edge({
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
            );

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
                redrawChild(nod, child.children[i], build);
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
                    redrawChild(root, jsontree.children[i], build);
                }
            }
        } else {
            redrawChild(null, jsontree, build);
        }
    }

    // 加亮显示当前运行到的某个 Node
    const debugFocus = (info: Array<string>) => {
        // clean
        cleanStepInfo();

        info.forEach(element => {
            findNode(element, (nod) => {
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

    const UnlockFocus = () => {
        if (currentClickNode.id == "")  {
            message.warning("please select a node")
            return
        }

        if (currentClickNode.type == NodeTy.Root)  {
            message.warning("root node can't be locked")
            return
        }

        if (currentLockedNode.id == "")  {  // locked

            props.dispatch(unlockFocus({
                id: currentClickNode.id,
                type: currentClickNode.type,
            }))
        } else {
            props.dispatch(unlockFocus({
                id: "",
                type: "",
            }))
        }

    }

    const ClickUpload = () => {
        setModalVisible(true);
    };

    // 基于某个模版，创建一个可运行的 Bot
    const ClickCreateDebug = (e: any) => {
        cleanStepInfo();

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

    const _step_loop = (botid :  string, done : any) => {
        Post(localStorage.remoteAddr, Api.DebugStep, { BotID: botid }).then(
            (json: any) => {

                if (json.Code !== 200) {
                    if (json.Code === 1009) {
                        message.warning(json.Code.toString() + " " + json.Msg)
                        done()
                        return;
                    } else if (json.Code === 1007) {
                        message.success("the end");
                        done()
                        return;
                    } else {
                        message.warning(json.Code.toString() + " " + json.Msg)
                    }
                }

                debugFocus([]); // clean

                // 推送 reponse 面板信息
                let threadinfo = JSON.parse(json.Body.ThreadInfo) as Array<ThreadInfo>

                // 推送当前节点信息
                let focusLst = new Array<string>
                try {
                    threadinfo.forEach(element => {
                        focusLst.push(element.curnod)
    
                        if (element.curnod === currentLockedNode.id) {
                            throw new Error("find")                            
                        }
                    });
                } catch (error) {
                    done()
                    return
                }
                
                debugFocus(focusLst)

                // 推送 meta 面板信息
                let metaStr = JSON.stringify(JSON.parse(json.Body.Blackboard))
                props.dispatch(setDebugInfo({
                    metaInfo: metaStr,
                    threadInfo: threadinfo,
                    lock: true,
                }))
            }
        );
    
    }

    const ClickStep = (e: any) => {
        if (props.tree.currentDebugBot === "") {
            message.warning("have not created bot");
            return;
        }

        var botid = props.tree.currentDebugBot

        if (currentLockedNode.id !== "") {
            const timer = new TaskTimer(100);

            timer.on('tick', () => {

                const donecallback = () =>{
                    console.info("step loop done")
                    timer.stop();
                    props.dispatch(setLock(false))
                }

                _step_loop(botid, donecallback)

            });
            timer.start();

        } else {
            props.dispatch(setLock(true))
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
    
                    debugFocus([]); // clean
    
                    // 推送 reponse 面板信息
                    let threadinfo = JSON.parse(json.Body.ThreadInfo) as Array<ThreadInfo>
    
                    // 推送当前节点信息
                    let focusLst = new Array<string>
                    threadinfo.forEach(element => {
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
    };

    const behaviorNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setBehaviorName(e.target.value);
    };

    const modalHandleOk = () => {
        setModalVisible(false);
        if (behaviorName !== "") {
            props.dispatch(nodeSave(behaviorName))
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
        //PubSub.publish(Topic.Undo, {});
    };

    const getLockedState = () => {
        if (currentLockedNode.id == "") {
            return <UnlockOutlined />
        } else {
            return <LockOutlined />
        }
    }

    return (
        <div className='app'>
            <EditSidePlane
                graph={graphRef.current}
                isGraphCreated={isGraphCreated} 
            />
            
            <div className="app-content" ref={containerRef} />

            <div
                className={"app-zoom-win"}
                style={{ marginLeft: 2, whiteSpace: "nowrap" }}
            >
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
                <Tooltip placement="topLeft" title="lock / unlock focus">
                    <Button icon={getLockedState()} onClick={UnlockFocus} />
                </Tooltip>
            </div>

            <div className={"app-step-win"}>
                <Tooltip placement="topRight" title={"Run to the next node [F10]"}>
                    <Button
                        type="primary"
                        style={{ width: 70 }}
                        icon={<CaretRightOutlined />}
                        disabled={lock}
                        onClick={ClickStep}
                    >
                        { }
                    </Button>
                </Tooltip>
            </div>
            <div className={"app-reset-win"}>
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
            <div className={"app-upload-win"}>
                <Tooltip placement="topRight" title={"Upload the bot to the server"}>
                    <Button
                        icon={<CloudUploadOutlined />}
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
