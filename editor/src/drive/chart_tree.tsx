import React from "react";
import { Treemap } from "@ant-design/charts";
import PubSub from "pubsub-js";
import Topic from "../model/topic";


type Info  = {
  Api : string,
  ConsumeNum : number,
}

type ReportChartInfo = {
  Chart : string,
  ApiList : Array<Info>,
}

type ApiInfo  = {
  name : string,
  value : number,
}


type TreeData = {
  name : string,
  children : Array<ApiInfo>,
}

type State = {
  data : TreeData,
}

export default class ApiChart extends React.Component {
  state : State = {
    data : {
      name : "root",
      children : [],
    }
  }

  componentDidMount() {

    PubSub.subscribe(Topic.ReportSelect, (topic:string, info:ReportChartInfo) => {

      let newdata : TreeData = {
        name : "root",
        children:[],
      }

      for (var i of info.ApiList) {
        let newinfo : ApiInfo = {
          name : i.Api,
          value : i.ConsumeNum,
        }
        newdata.children.push(newinfo)
      }

      this.setState({data : newdata})
    })
  }

  render() {
    var config = {
      data: this.state.data,
      colorField: "name",
    };

    return (
      <div>
        <Treemap {...config} />
      </div>
    );
  }
}
