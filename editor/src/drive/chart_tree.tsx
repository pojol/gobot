import React from "react";
import { Treemap } from "@ant-design/charts";
import PubSub from "pubsub-js";
import Topic from "../model/topic";


type Info  = {
  Api : string,
  Value : number,
}

type ReportChartInfo = {
  Chart : string,
  ApiList : Array<Info>,
}

type ApiInfo  = {
  ty : string,
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
          ty : info.Chart,
          name : i.Api,
          value : i.Value,
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
      tooltip: {
        formatter: (datum: any) => {

          if (datum.path[0].data.ty === "avg_request_time_ms") {
            return { name: datum.name, value: datum.value + ' ms' };
          } else if (datum.path[0].data.ty === "request_times") {
            return { name: datum.name, value: datum.value + " times" };
          }

          return { name: datum.name, value: datum.value };
        },
      }
    };

    return (
      <div>
        <Treemap {...config} />
      </div>
    );
  }
}
