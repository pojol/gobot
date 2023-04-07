import HeartTask from "@/task/heart";
import * as React from "react";

// You will need to import the styles separately
// You probably want to do this just once during the bootstrapping phase of your application.
import "react-reflex/styles.css";

// then you can import the components
import { ReflexContainer, ReflexSplitter, ReflexElement, HandlerProps } from "react-reflex";

/// <reference path="@/edit/node.d.ts" />

import "./editor.css"
import GraphView from "./edit/graph"
import Blackboard from "./edit/blackboard";
import Nodes from "./edit/node/tab"

import { NodeTy } from "@/constant/node_type";

export default class Editor extends React.Component {

  componentDidMount() {
    if (localStorage.codeboxTheme === undefined || localStorage.codeboxTheme === "") {
      localStorage.codeboxTheme = "default"
    }
  }

  onResizeEditPane(domElement: HandlerProps) {
    //PubSub.publish(Topic.EditPanelCodeMetaResize, domElement.component.props.flex)
    console.info("resize", domElement.component.props)
  }

  onResizeGraphPane(domElement: HandlerProps) {
    //PubSub.publish(Topic.EditPanelEditCodeResize, domElement.component.props.flex)
    console.info("resize", domElement.component.props)
  }

  onResizeChangePane(domElement: HandlerProps) {
    //PubSub.publish(Topic.EditPanelEditChangeResize, domElement.component.props.flex)
    console.info("resize", domElement.component.props)
  }

  render() {
    return (
      <div>
        <HeartTask />
        <div className="container">
          <ReflexContainer orientation="vertical">
            <ReflexElement className="left-pane" flex={0.6} minSize={200} onStopResize={this.onResizeGraphPane}>
              <ReflexContainer orientation="horizontal">
                <ReflexElement className="left-pane" minSize={300} flex={1} >
                  <GraphView />
                </ReflexElement>
              </ReflexContainer>
            </ReflexElement>

            <ReflexSplitter propagate={true} />

            <ReflexElement className="right-pane" flex={0.4} minSize={100}>
              <ReflexContainer orientation="horizontal">
                <ReflexElement className="left-pane" minSize={100} propagateDimensions={true} onStopResize={this.onResizeEditPane}>
                  <Nodes {...{ nodety: NodeTy.Action, dimensions: { width: 0, height: 0 } }}/>
                </ReflexElement>

                <ReflexSplitter />

                <ReflexElement className="left-pane" minSize={100}>
                  <Blackboard />
                </ReflexElement>
              </ReflexContainer>
            </ReflexElement>
          </ReflexContainer>
        </div>
      </div>
    );
  }

}