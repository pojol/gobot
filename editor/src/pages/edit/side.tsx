import React, { useEffect, useState, useLayoutEffect } from 'react';
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

interface SideProps {
  graph: Graph
  isGraphCreated: boolean
}


export const EditSidePlane: React.FC<SideProps> = ({ graph, isGraphCreated }) => {

  const { themeValue } = useSelector((state: RootState) => state.configSlice)
  const { pmap } = useSelector((state: RootState) => state.prefabSlice)
  const [dnd, setdnd] = useState<any>()
  const [selectTags, setSelectTags] = useState(new Array<string>())
  const [prefabLst, setPrefabLst] = useState(new Array<string>())
  const [tags, setTags] = useState(new Array<PrefabTagInfo>())

  useEffect(() => {
    if (isGraphCreated) {
      console.info("reload side plane")

      setdnd(new Dnd({
        target: graph,
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
      }))

      resizeSidePane()
      reloadPrefab()
    }

  }, [isGraphCreated])

  /*
    PubSub.subscribe(Topic.PrefabUpdateAll, (topic: string, info: any) => {
      this.reloadPrefab()
    });
  
  */
  const matchTags = (tags: string[]): boolean => {

    let selecttags = selectTags
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

  const reloadPrefab = () => {
    const prefabMap = pmap

    setSelectTags([])
    setTags([])

    let tmplst = new Array<string>()
    let taglst = new Array<PrefabTagInfo>()
    var tagSet = new Set<string>()

    prefabMap.forEach((value: PrefabInfo) => {

      if (selectTags.length !== 0) {
        if (matchTags(value.tags)) {
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
    setPrefabLst(tmplst)
    setTags(taglst)
  }

  const resizeSidePane = () => {
    var div = document.getElementById("prefab-pane")
    if (div !== null) {
      var clienth = document.documentElement.clientHeight - 260
      div.style.height = clienth.toString() + "px"
      div.style.overflow = "auto"
      console.info("set", document.documentElement.clientHeight.toString())
    }
  }

  const startDrag = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
    const target = e.currentTarget
    const name = target.getAttribute('data-type') as string
    let nod: Node

    if (IsPresetNode(name)) {
      nod = GetNode(NodeTy.Action, {})
    } else {
      nod = GetNode(name, {})
    }

    nod.setAttrs({ type: { name: name }, label: { text: name } })

    dnd.start(nod, e.nativeEvent as any)
  }

  const onSelectChange = (value: string[]) => {
    setSelectTags(value)
    reloadPrefab()
  };


  return (
    <div className="dnd-wrap">

      <Row justify="space-around" align="middle" gutter={[22, 12]}>
        <Col span={7}>
          <div
            data-type="SelectorNode"
            className="dnd-selector"
            onMouseDown={startDrag}
          >
            Selector
          </div>
        </Col>
        <Col span={7}>
          <div
            data-type="SequenceNode"
            className="dnd-sequence"
            onMouseDown={startDrag}
          >
            Sequence
          </div>
        </Col>
        <Col span={7}>
          <div
            data-type="ParallelNode"
            className="dnd-parallel"
            onMouseDown={startDrag}
          >
            Parallel
          </div>
        </Col>

        <Col span={7}>
          <div
            data-type="ConditionNode"
            className="dnd-condition"
            onMouseDown={startDrag}
          >
            Condition
          </div>
        </Col>
        <Col span={7}>
          <div
            data-type="LoopNode"
            className="dnd-loop"
            onMouseDown={startDrag}
          >
            Loop
          </div>
        </Col>
        <Col span={7}>
          <div
            data-type="WaitNode"
            className="dnd-wait"
            onMouseDown={startDrag}
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
        options={tags}
        onChange={onSelectChange}
      />

      <Divider>Prefab</Divider>

      <div id="prefab-pane" className="dnd-warp-prefab">
        {prefabLst.map((item: any, index: number) =>
          <div
            key={index}
            data-type={item}
            className={"dnd-prefab-" + themeValue}
            onMouseDown={startDrag}
          >
            {item}
          </div>
        )}
      </div>

    </div>
  );
}
