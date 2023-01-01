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

const { Dnd } = Addon

interface SideProps {
  graph: Graph
}


export default class EditSidePlane extends React.Component<SideProps> {
  private graph: Graph
  //private dndContainer: HTMLDivElement
  private dnd: any

  state = {
    prefabLst: []
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

  reloadPrefab() {
    let configmap = (window as any).config as Map<string, string>;

    this.setState({ prefabLst: [] })
    let tmplst = new Array<string>()

    configmap.forEach((value: string, key: string) => {
      if (key !== "system" && key !== "global" && key !== "") {
        console.info("reload prefab", key)
        tmplst.push(key)
      }
    });

    this.setState({ prefabLst: tmplst })
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

    PubSub.subscribe(Topic.ConfigUpdateAll, (topic: string, info: any) => {
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
    const type = target.getAttribute('data-type')
    var nod = new Node;
    
    switch (type) {
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
    }

    this.dnd.start(nod, e.nativeEvent as any)
  }

  dndContainerRef = (container: HTMLDivElement) => {
    //this.dndContainer = container
  }

  onSelectChange = (value: string[]) => {
    console.log(`selected ${value}`);
  };

  render() {
    const options = [{ value: 'gold' }, { value: 'lime' }, { value: 'green' }, { value: 'cyan' }];

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
          defaultValue={['gold', 'cyan']}
          style={{ width: '100%' }}
          options={options}
          onChange={this.onSelectChange}
        />

        <Divider>Prefab</Divider>

        <div id="prefab-pane" className="dnd-warp-prefab">
          {this.state.prefabLst.map((item: string) =>
            <div
              data-type="rect"
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
