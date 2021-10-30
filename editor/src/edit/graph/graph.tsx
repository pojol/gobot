import * as React from "react";
import { Graph, Addon, Shape, Cell, Node } from "@antv/x6";
import ActionNode from "../../shape/shape_action";
import ConditionNode from "../../shape/shape_condition";
import SelectorNode from "../../shape/shape_selector";
import SequenceNode from "../../shape/shape_sequence";
import RootNode from "../../shape/shape_root";
import LoopNode from "../../shape/shape_loop";
import WaitNode from "../../shape/shape_wait";
import AssertNode from "../../shape/shap_assert";
import { NodeTy, IsScriptNode } from "../../model/node_type";

import "./graph.css";
import { message } from "antd";
import PubSub from "pubsub-js";
import Topic from "../../model/topic";

const { Dnd, Stencil } = Addon;

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

// 在节点上绑定数据
// 打印出树状结构
// 将信息传递给不同的view （ ant -> graph node -> node metadata

export default class GraphView extends React.Component {
  graph: Graph;
  container: HTMLElement;
  dnd: any;
  stencilContainer: HTMLDivElement;

  componentDidMount() {
    // 新建画布
    const graph = new Graph({
      container: this.container,
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
    });

    var root = new RootNode();
    graph.addNode(root);

    PubSub.publish(Topic.NodeAdd, {
      parent: root.id,
      child: this.getNodInfo(root),
    });

    const stencil = new Stencil({
      title: "Components",
      search(nod, keyword) {
        var attr = nod.getAttrs();
        var label = attr.label.text as String;
        if (label !== null) {
          return label.toLowerCase().indexOf(keyword.toLowerCase()) !== -1;
        }

        return false;
      },
      placeholder: "Search by shape name",
      notFoundText: "Not Found",
      target: graph,
      collapsable: true,
      stencilGraphWidth: 180,
      stencilGraphHeight: 100,
      groups: [
        {
          name: "group1",
          title: "Control",
        },
        {
          name: "group2",
          title: "Condition",
        },
        {
          name: "group3",
          title: "Script",
        },
        {
          name: "group4",
          title: "Decorator",
        },
      ],
    });
    this.stencilContainer.appendChild(stencil.container);

    stencil.load([new SelectorNode(), new SequenceNode()], "group1");
    stencil.load([new ConditionNode(), new AssertNode()], "group2");
    stencil.load([new ActionNode()], "group3");
    stencil.load([new LoopNode(), new WaitNode()], "group4");

    graph.bindKey("del", () => {
      const cells = this.graph.getSelectedCells();

      if (cells.length) {
        for (var i = 0; i < cells.length; i++) {

          if (cells[i].getAttrs().type.toString() !== NodeTy.Root) {

            if (cells[i].getParent() == null) {
              graph.removeCell(cells[i])
            } else {
              PubSub.publish(Topic.NodeRmv, cells[i].id);
              cells[i].getParent()?.removeChild(cells[i]);
            }
          }
        }
      }
      return false;
    });

    graph.on("edge:removed", ({ edge, options }) => {
      if (!options.ui) {
        return;
      }

      this.findNode(edge.getTargetCellId(), (child) => {
        //var ts = child.removeFromParent( { deep : false } );  // options 没用？
        PubSub.publish(Topic.LinkRmv, child.id);
        child.getParent()?.removeChild(edge);
        //var ts = child.removeFromParent({ deep: false });
        //this.graph.addCell(ts);
      });

      //graph.removeEdge(edge.id);
    });

    graph.on("edge:connected", ({ isNew, edge }) => {
      const source = edge.getSourceNode();
      const target = edge.getTargetNode();

      if (isNew) {
        if (source !== null && target !== null) {
          edge.setZIndex(0)
          source.addChild(target);
          this.connect(source, target);
        }
      }
    });

    graph.on("node:click", ({ node }) => {
      PubSub.publish(Topic.NodeClick, {
        id: node.id,
        type: node.getAttrs().type,
      });
    });

    graph.on("node:added", ({ node, index, options }) => {
      var ty = node.getAttrs().type.toString();
      if (ty === NodeTy.Selector || ty === NodeTy.Sequence || ty === NodeTy.Root) {
        return;
      }

      node.setAttrs({
        label: {
          text: "",
        },
      });
    });

    graph.on("node:removed", ({ node, index, options }) => { });
    graph.on("node:moved", ({ e, x, y, node, view: NodeView }) => {
      this.findNode(node.id, (nod) => {
        PubSub.publish(Topic.UpdateGraphParm, this.getNodInfo(node));
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
          },
        },
      ]);
    });

    graph.on("edge:mouseleave", ({ edge }) => {
      edge.removeTools();
    });

    graph.centerContent();
    this.dnd = new Dnd({
      target: graph,
      scaled: false,
      animation: true,
    });
    this.graph = graph;

    PubSub.subscribe(Topic.UpdateNodeParm, (topic: string, info: any) => {
      if (info.parm.ty === NodeTy.Action) {
        this.findNode(info.parm.id, (nod) => {
          nod.setAttrs({
            label: { text: info.parm.alias },
          });
        });
      } else if (info.parm.ty === NodeTy.Loop) {
        this.findNode(info.parm.id, (nod) => {
          console.info("Topic.UpdateNodeParm", nod.id)
          nod.setAttrs({
            label: { text: this.getLoopLabel(info.parm.loop) },
          });
        });
      } else if (info.parm.ty === NodeTy.Wait) {
        this.findNode(info.parm.id, (nod) => {
          nod.setAttrs({
            label: { text: info.parm.wait.toString() + " ms" },
          });
        });
      }
    });

    PubSub.subscribe(Topic.FileLoadGraph, (topic: string, jsontree: any) => {
      // 通过 json 重绘画布
      this.graph.clearCells();

      if (jsontree.id !== "") {
        this.redraw(jsontree);
      }
    });

    PubSub.subscribe(Topic.Focus, (topic: string, info: any) => {
      if (info.Cur !== "") {
        this.findNode(info.Cur, (nod) => {
          nod.setAttrs({
            body: {
              strokeWidth: 3,
            },
          });
        });
      }
      if (info.Prev !== "") {
        this.findNode(info.Prev, (nod) => {
          nod.setAttrs({
            body: {
              strokeWidth: 1,
            },
          });
        });
      }
    });

    PubSub.subscribe(Topic.Create, (topic: string, info: any) => {
      this.refreshNodes((nod)=>{  // 
        nod.setAttrs({
          body: {
            strokeWidth: 1,
          },
        });
      })
    })
  }

  connect(source: Node, target: Node) {
    var nodinfo = {
      parent: source.id,
      child: this.getNodInfo(target),
    };
    PubSub.publish(Topic.NodeAdd, nodinfo);
  }

  getLoopLabel(val: Number) {
    var tlab = "";
    if (val !== 0) {
      tlab = val.toString() + " times";
    } else {
      tlab = "endless";
    }

    return tlab;
  }

  redrawChild(parent: Node, child: any) {
    for (var i = 0; i < child.length; i++) {
      var nod: Node;
      if (child[i].ty === NodeTy.Selector) {
        nod = new SelectorNode({ id: child[i].id });
      } else if (child[i].ty === NodeTy.Sequence) {
        nod = new SequenceNode({ id: child[i].id });
      } else if (child[i].ty === NodeTy.Condition) {
        nod = new ConditionNode({ id: child[i].id });
      } else if (child[i].ty === NodeTy.Action) {
        nod = new ActionNode({ id: child[i].id });
      } else if (child[i].ty === NodeTy.Loop) {
        nod = new LoopNode({ id: child[i].id });
      } else if (child[i].ty === NodeTy.Assert) {
        nod = new AssertNode({ id: child[i].id });
      } else if (child[i].ty === NodeTy.Wait) {
        nod = new WaitNode({ id: child[i].id });
      } else {
        message.warn("未知的节点类型" + child[i].ty);
        break;
      }
      nod.setPosition({
        x: child[i].pos.x,
        y: child[i].pos.y,
      });
      this.graph.addNode(nod);
      this.graph.addEdge(
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
      this.connect(parent, nod);

      if (IsScriptNode(child[i].ty)) {
        nod.setAttrs({ label: { text: child[i].alias } })
        PubSub.publish(Topic.UpdateNodeParm, {
          parm: {
            id: nod.id,
            ty: child[i].ty,
            code: child[i].code,
            alias: child[i].alias,
          },
          notify: false,
        });
      } else if (child[i].ty === NodeTy.Loop) {
        nod.setAttrs({ label: { text: this.getLoopLabel(child[i].loop) } });
        PubSub.publish(Topic.UpdateNodeParm, {
          parm: {
            id: nod.id,
            ty: child[i].ty,
            loop: child[i].loop,
          },
          notify: false,
        });
      } else if (child[i].ty === NodeTy.Wait) {
        nod.setAttrs({ label: { text: child[i].wait.toString() + " ms" } });
        PubSub.publish(Topic.UpdateNodeParm, {
          parm: {
            id: nod.id,
            ty: child[i].ty,
            wait: child[i].wait,
          },
          notify: false,
        });
      }

      if (child[i].children && child[i].children.length) {
        this.redrawChild(nod, child[i].children);
      }
    }
  }

  redraw(jsontree: any) {
    var root = new RootNode();
    root.setPosition({
      x: jsontree.pos.x,
      y: jsontree.pos.y,
    });
    this.graph.addNode(root);
    PubSub.publish(Topic.NodeAdd, {
      parent: root.id,
      child: this.getNodInfo(root),
    });

    if (jsontree.children && jsontree.children.length) {
      this.redrawChild(root, jsontree.children);
    }
  }

  setLabel(id: String, name: String) {
    var flag = false;
    this.findNode(id, (nod) => {
      flag = true;
    });

    if (!flag) {
      message.warning("没有在树中查找到该节点 " + id);
    }
  }

  fillChildInfo(child: Node, info: any) {
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
        this.fillChildInfo(cchild as Node, childInfo);
      }
    });
  }

  getNodInfo(nod: Node) {
    var info = {
      id: nod.id,
      ty: nod.getAttrs().type.toString(),
      pos: {
        x: nod.position().x,
        y: nod.position().y,
      },
      children: [],
    };

    nod.eachChild((child, idx) => {
      if (child instanceof Node) {
        this.fillChildInfo(child as Node, info);
      }
    });

    return info;
  }

  refStencil = (container: HTMLDivElement) => {
    this.stencilContainer = container;
  };

  refContainer = (container: HTMLDivElement) => {
    this.container = container;
  };

  findChild = (parent: Cell, id: String, callback: (nod: Cell) => void) => {
    if (parent.id === id) {
      callback(parent);
      return;
    } else {
      parent.eachChild((child, idx) => {
        this.findChild(child, id, callback);
      });
    }
  };

  findNode = (id: String, callback: (nod: Cell) => void) => {
    var nods = this.graph.getRootNodes();
    if (nods.length >= 0) {
      if (nods[0].id === id) {
        callback(nods[0]);
      } else {
        nods[0].eachChild((child, idx) => {
          this.findChild(child, id, callback);
        });
      }
    }
  };

  refreshNode = (parent: Cell, callback: (nod: Cell) => void) => {
    callback(parent)
    parent.eachChild((child, idx) => {
      this.refreshNode(child, callback)
    })
  }

  refreshNodes = (callback: (nod: Cell) => void) => {
    var nods = this.graph.getRootNodes();
    if (nods.length >= 0) {
      callback(nods[0]);
      nods[0].eachChild((child, idx) => {
        this.refreshNode(child, callback)
      })
    }
  }

  debug = () => { };

  render() {
    return (
      <div className="app">
        <div className="app-stencil" ref={this.refStencil} />
        <div className="app-content" ref={this.refContainer} />
      </div>
    );
  }
}
