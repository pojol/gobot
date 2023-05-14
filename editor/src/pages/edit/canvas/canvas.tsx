import { NodeTy } from "@/constant/node_type";
import { Cell, Graph, Node, Shape } from "@antv/x6";

// 高亮
const magnetAvailabilityHighlighter = {
    name: "stroke",
    args: {
        attrs: {
            fill: "#fff",
            stroke: "#47C769",
        },
    },
};

export function CreateGraph(container: HTMLElement, wflex: number, hflex: number): Graph {
    const graph = new Graph({
        width: document.documentElement.clientWidth * wflex,
        height: document.documentElement.clientHeight,
        container: container,
        highlighting: {
            magnetAvailable: magnetAvailabilityHighlighter,
            magnetAdsorbed: {
                name: "stroke",
                args: {
                    attrs: {
                        fill: "#fff",
                        stroke: "#31d0c6",
                    },
                },
            },
        },
        snapline: {
            enabled: true,
            sharp: true,
        },
        connecting: {
            snap: true,
            allowBlank: false,
            allowLoop: false,
            allowPort: false,
            highlight: true,
            allowMulti: false,
            connector: "rounded",
            connectionPoint: "boundary",
            router: {
                name: "er",
                args: {
                    direction: "V",
                },
            },
            createEdge() {
                return new Shape.Edge({
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
                });
            },
        },
        keyboard: {
            enabled: true,
        },
        grid: {
            size: 10, // 网格大小 10px
            visible: true, // 绘制网格，默认绘制 dot 类型网格
        },
        history: true,
        selecting: {
            enabled: true,
            showNodeSelectionBox: true,
        },
        scroller: {
            enabled: true,
            pageVisible: false,
            pageBreak: false,
            pannable: true,
        },
        mousewheel: {
            enabled: true,
            modifiers: ["alt", "meta"],
        },
    });

    // 调整画布大小
    graph.resizeGraph(1024, 1024);

    return graph
}
