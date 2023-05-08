import React, { useEffect } from 'react';
import { Graph, Addon, Dom, Node } from "@antv/x6";
import { Select, Divider, Row, Col } from "antd";

import { useSelector, connect, ConnectedProps } from 'react-redux';

import "./side.css";
import PubSub from "pubsub-js";

import {
  IsActionNode,
  NodeTy,
} from "../../constant/node_type"
import { GetNode } from "./shape/shape";
import Topic from '@/constant/topic';
import { RootState } from '@/models/store';
import { IsPresetNode } from '../../constant/node_type';
import { GetNodInfo } from '@/models/node';

const { Dnd } = Addon

interface SideProps extends PropsFromRedux {
  graph: Graph
}

class EditSidePlane extends React.Component<SideProps> {
  private graph!: Graph;
  private dnd: any

  state = {
    prefabLst: [],
    tags: [],
    selectTags: [],
    theme: localStorage.theme,
  }

  componentDidUpdate(prevProps: SideProps) {
    this.graph = this.props.graph
    console.info("update side", this.props.graph)
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

  componentDidMount() {
    //const { graph } = this.props;
    
    PubSub.subscribe(Topic.PrefabUpdateAll, (topic: string, info: any) => {
      this.reloadPrefab()
    });

    PubSub.subscribe(Topic.ThemeChange, (topic: string, theme: string) => {
      this.setState({ theme: theme })
    });

    /*
        PubSub.subscribe(
          Topic.EditPanelEditCodeResize,
          (topic: string, flex: number) => {
            this.resizeSidePane()
          }
        );
    
        PubSub.subscribe(Topic.WindowResize, () => {
          this.resizeSidePane()
        });
        */
    this.resizeSidePane()
    this.reloadPrefab()
  }

  matchTags(tags: string[]): boolean {

    let selecttags = this.state.selectTags
    for (var i = 0; i < selecttags.length; i++) {
      for (var j = 0; j < tags.length; j++) {
        if (tags[j] === selecttags[i]) {
          console.info(tags[j], "match", selecttags[i])
          return true
        }
      }
    }

    return false
  }

  reloadPrefab() {
    const prefabMap = this.props.prefabMap

    this.setState({ prefabLst: [], tags: [] })
    let tmplst = new Array<string>()
    let taglst = new Array<PrefabTagInfo>()
    var tagSet = new Set<string>()

    prefabMap.forEach((value: PrefabInfo) => {

      if (this.state.selectTags.length !== 0) {
        if (this.matchTags(value.tags)) {
          tmplst.push(value.name)
        }
      } else {
        tmplst.push(value.name)
      }

      for (var i = 0; i < value.tags.length; i++) {
        tagSet.add(value.tags[i])
      }
    });

    tagSet.forEach(element => {
      taglst.push({ value: element })
    });

    tmplst.sort
    this.setState({ prefabLst: tmplst, tags: taglst })
  }

  resizeSidePane() {
    var div = document.getElementById("prefab-pane")
    if (div !== null) {
      var clienth = document.documentElement.clientHeight - 260
      div.style.height = clienth.toString() + "px"
      div.style.overflow = "auto"
      console.info("set", document.documentElement.clientHeight.toString())
    }
  }

  startDrag = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
    const target = e.currentTarget
    const name = target.getAttribute('data-type') as string
    let nod: Node

    if (IsPresetNode(name)) {
      nod = GetNode(NodeTy.Action, {})
    } else {
      nod = GetNode(name, {})
    }

    nod.setAttrs({ type: { name: name }, label: { text: name } })

    this.dnd.start(nod, e.nativeEvent as any)
  }

  dndContainerRef = (container: HTMLDivElement) => {
    //this.dndContainer = container
}

  onSelectChange = (value: string[]) => {
    this.setState({ selectTags: value }, () => {
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
          {this.state.prefabLst.map((item: any, index: number) =>
            <div
              key={index}
              data-type={item}
              className={"dnd-prefab-" + this.state.theme}
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


const mapStateToProps = (state: RootState) => ({
  prefabMap: state.prefabSlice.pmap
});

const connector = connect(mapStateToProps);
type PropsFromRedux = ConnectedProps<typeof connector>;

export default connector(EditSidePlane);
