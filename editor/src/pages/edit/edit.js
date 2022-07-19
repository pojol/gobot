import * as React from "react";
import PubSub from "pubsub-js";

import GraphView from "./graph/graph";
import Edit from "./node/edit_tab";
import Blackboard from "./meta/meta";

import Topic from "../../constant/topic";

// You will need to import the styles separately
// You probably want to do this just once during the bootstrapping phase of your application.
import "react-reflex/styles.css";

// then you can import the components
import { ReflexContainer, ReflexSplitter, ReflexElement } from "react-reflex";

import "./edit.css";

export default class EditPlane extends React.Component {
  constructor(props) {
    super(props);
    this.state = {};
  }

  componentDidMount() {
    
  }

  onResizeEditPane(domElement, component) {
    console.info("resize", domElement.component.props)

    PubSub.publish(Topic.EditPanelCodeMetaResize, domElement.component.props.flex)
  }

  onResizeGraphPane(domElement) {
    console.info("resize", domElement.component.props)
    PubSub.publish(Topic.EditPanelEditCodeResize, domElement.component.props.flex)
  }

  onResizeChangePane(domElement) {
    console.info("resize", domElement.component.props)
    PubSub.publish(Topic.EditPanelEditChangeResize, domElement.component.props.flex)
  }

  render() {

    return (
      <div className="container">
        <ReflexContainer orientation="vertical">
          <ReflexElement className="left-pane" flex={0.6} minSize="200" onStopResize={this.onResizeGraphPane}>
            <ReflexContainer orientation="horizontal">
              <ReflexElement className="left-pane" minSize="300" flex={1} >
                <GraphView />
              </ReflexElement>
            </ReflexContainer>
          </ReflexElement>

          <ReflexSplitter propagate={true} />

          <ReflexElement className="right-pane" flex={0.4} minSize="100">
            <ReflexContainer orientation="horizontal">
              <ReflexElement className="left-pane" minSize="100" propagateDimensions={true} onStopResize={this.onResizeEditPane}>
                <Edit />
              </ReflexElement>

              <ReflexSplitter />

              <ReflexElement className="left-pane" minSize="100">
                <Blackboard />
              </ReflexElement>
            </ReflexContainer>
          </ReflexElement>
        </ReflexContainer>
      </div>
    );
  }
}
