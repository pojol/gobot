import * as React from "react";
import { Graph, Addon, Dom, Node } from "@antv/x6";
import { Select, Divider, Row, Col } from "antd";
import "./side.css";

import ConditionNode from "../graph/shape/shape_condition";
import SelectorNode from "../graph/shape/shape_selector";
import SequenceNode from "../graph/shape/shape_sequence";
import LoopNode from "../graph/shape/shape_loop";
import WaitNode from "../graph/shape/shape_wait";
import ParallelNode from "../graph/shape/shap_parallel";

import PubSub from "pubsub-js";
import Topic from "../../../constant/topic";

import {
  NodeTy,
} from "../../../constant/node_type";
import ActionNode from "../graph/shape/shape_action";

const { Dnd } = Addon

interface SideProps {
  graph: Graph
}

interface PrefabInfo {
  name: string,
  tags: string[],
  code: string,
}

interface PrefabTagInfo {
  value: string
}

export default class EditSidePlane extends React.Component<SideProps> {
  private graph: Graph
  //private dndContainer: HTMLDivElement
  private dnd: any

  state = {
    prefabLst: [],
    tags: [],
    selectTags: [],
  }

  componentWillReceiveProps(newProps: SideProps) {
    this.graph = newProps.graph
    this.dnd = new Dnd({
      target: this.graph,
      scaled: false,
      animation: true,
      //dndContainer: this.dndContainer,
      validateNode(droppingNode, options) {
        return droppingNode.shape === 'html'
          ? new Promise<boolean>((resolve) => {
            const { draggingNode, draggingGraph } = options
            const view = draggingGraph.findView(draggingNode)!
            const contentElem = view.findOne('foreignObject > body > div')
            Dom.addClass(contentElem, 'validating')
            setTimeout(() => {
              Dom.removeClass(contentElem, 'validating')
              resolve(true)
            }, 3000)
          })
          : true
      },
    })
  }

  matchTags(tags: string[]): boolean {

    let selecttags = this.state.selectTags
    for (var i = 0; i < selecttags.length; i++) {
      for (var j = 0; j < tags.length; j++) {
        if (tags[j] == selecttags[i]) {
          console.info(tags[j], "match", selecttags[i])
          return true
        }
      }
    }

    return false
  }

  reloadPrefab() {
    let configmap = (window as any).prefab as Map<string, PrefabInfo>;

    this.setState({ prefabLst: [], tags: [] })
    let tmplst = new Array<string>()
    let taglst = new Array<PrefabTagInfo>()
    var tagSet = new Set<string>()

    console.info("select tags", this.state.selectTags.length)

    configmap.forEach((value: PrefabInfo, key: string) => {

      if (this.state.selectTags.length !== 0) {
        console.info("need match")
        if (this.matchTags(value.tags)) {
          tmplst.push(key)
        }
      } else {
        tmplst.push(key)
      }

      for (var i = 0; i < value.tags.length; i++) {
        tagSet.add(value.tags[i])
      }
    });

    tagSet.forEach(element => {
      taglst.push({ value: element })
    });

    console.info("prefab lst", tmplst, "reload tags", taglst)
    this.setState({ prefabLst: tmplst, tags: taglst })
  }

  resizeSidePane() {
    console.info("resizeSidePane")
    var div = document.getElementById("prefab-pane")
    if (div !== null) {
      var clienth = document.body.clientHeight - 260
      div.style.height = clienth.toString() + "px"
      div.style.overflow = "auto"
      console.info("set", document.body.clientHeight.toString())
    }
  }

  componentWillMount() {

    PubSub.subscribe(Topic.PrefabUpdateAll, (topic: string, info: any) => {
      this.reloadPrefab()
    });

    PubSub.subscribe(
      Topic.EditPanelEditCodeResize,
      (topic: string, flex: number) => {
        this.resizeSidePane()
      }
    );

    PubSub.subscribe(Topic.WindowResize, () => {
      this.resizeSidePane()
    });
  }

  componentDidMount() {
    this.resizeSidePane()
  }

  startDrag = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
    const target = e.currentTarget
    const ty = target.getAttribute('data-type') as any
    var nod = new Node;

    switch (ty) {
      case NodeTy.Selector:
        nod = new SelectorNode();
        break;
      case NodeTy.Sequence:
        nod = new SequenceNode();
        break;
      case NodeTy.Condition:
        nod = new ConditionNode();
        break;
      case NodeTy.Loop:
        nod = new LoopNode();
        break;
      case NodeTy.Wait:
        nod = new WaitNode();
        break;
      case NodeTy.Parallel:
        nod = new ParallelNode()
        break;
      default:
        nod = new ActionNode()
        nod.setAttrs({ type: ty })
    }

    this.dnd.start(nod, e.nativeEvent as any)
  }

  dndContainerRef = (container: HTMLDivElement) => {
    //this.dndContainer = container
  }

  onSelectChange = (value: string[]) => {
    this.setState({ selectTags: value }, ()=>{
      this.reloadPrefab()
    })
  };

  render() {

    return (
      <div className="dnd-wrap" ref={this.dndContainerRef}>

        <Row justify="space-around" align="middle" gutter={[22, 12]}>
          <Col span={7}>
            <div
              data-type="SelectorNode"
              className="dnd-selector"
              onMouseDown={this.startDrag}
            >
              Selector
            </div>
          </Col>
          <Col span={7}>
            <div
              data-type="SequenceNode"
              className="dnd-sequence"
              onMouseDown={this.startDrag}
            >
              Sequence
            </div>
          </Col>
          <Col span={7}>
            <div
              data-type="ParallelNode"
              className="dnd-parallel"
              onMouseDown={this.startDrag}
            >
              Parallel
            </div>
          </Col>

          <Col span={7}>
            <div
              data-type="ConditionNode"
              className="dnd-condition"
              onMouseDown={this.startDrag}
            >
              Condition
            </div>
          </Col>
          <Col span={7}>
            <div
              data-type="LoopNode"
              className="dnd-loop"
              onMouseDown={this.startDrag}
            >
              Loop
            </div>
          </Col>
          <Col span={7}>
            <div
              data-type="WaitNode"
              className="dnd-wait"
              onMouseDown={this.startDrag}
            >
              Wait
            </div>
          </Col>
        </Row>


        <Divider>Filter</Divider>

        <Select
          mode="multiple"
          showArrow
          style={{ width: '100%' }}
          options={this.state.tags}
          onChange={this.onSelectChange}
        />

        <Divider>Prefab</Divider>

        <div id="prefab-pane" className="dnd-warp-prefab">
          {this.state.prefabLst.map((item: string) =>
            <div
              data-type={item}
              className="dnd-prefab"
              onMouseDown={this.startDrag}
            >
              {item}
            </div>
          )}
        </div>


      </div>

    );
  }
}
