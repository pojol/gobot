import { Table, Tag, Tabs,message} from "antd";
import * as React from "react";
import ApiChart from "./chart_tree";
import PubSub from "pubsub-js";
import Topic from "../model/topic";
import { Post } from "../model/request";
import Api from "../model/api";

const { TabPane } = Tabs;

export default class TestReport extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        {
          title: "Time",
          dataIndex: "time",
          key: "time",
        },
        {
          title: "TPS",
          dataIndex: "tps",
          key: "tps",
        },
        {
          title: "Duration",
          dataIndex: "duration",
          key: "duration",
        },
        {
          title: "BotNum",
          key: "botnum",
          dataIndex: "botnum",
        },
        {
          title: "ReqNum",
          key: "reqnum",
          dataIndex: "reqnum",
        },
        {
          title: "Errors",
          key: "errors",
          dataIndex: "errors",
        },
        {
          title: 'Charts',
          key: 'charts',
          dataIndex: 'charts',
          render: (tags,record) => (
              <>
                {tags.map(tag => {
                  let color = 'green';
                  if (tag === 'avg_request_time_ms') {
                    color = 'volcano';
                  } else if (tag === 'request_times') {
                      color = 'geekblue';
                  }
                  return (
                    <Tag color={color} key={tag} onClick={()=>this.clickTag(tag,record)}>
                    {tag}
                  </Tag>

                  );
                })}
              </>
            ),
      }

      ],
      data: [],
    };
  }

  fillData(info) {
    var newdata = [] 
    for (var i = 0; i < info.length; i++) {
      newdata.push({
        key: info[i].ID,
        time: info[i].BeginTime,
        tps:info[i].Tps,
        duration: info[i].Dura,
        botnum: info[i].BotNum,
        reqnum: info[i].ReqNum,
        errors : info[i].ErrNum,
        charts: ["avg_request_time_ms", "request_times"],
        apilst : info[i].Apilst,
      })
    }
    this.setState({data:newdata})
  }

  refresh() {
    this.setState({data:[]})
    Post(window.remote, Api.ReportInfo, {}).then((json) => {
      if (json.Code !== 200) {
        message.error("run fail:" + String(json.Code) + " msg: " + json.Msg);
      } else {
        console.info(json.Body.Info);
        if (json.Body.Info){
          this.fillData(json.Body.Info)
        }
      }
    });
  }

  componentDidMount() {
    PubSub.subscribe(Topic.ReportUpdate, (topic, info) => {
      this.refresh();
    });

    this.refresh();
  }

  clickTag = (e,record) => {
    if (record.apilst) {
      let lst = new Array()

      if (e === "avg_request_time_ms") {
        for (var i = 0; i < record.apilst.length; i++) {
          lst.push({"Api": record.apilst[i].Api, "Value": record.apilst[i].ConsumeNum})
        }
      } else if (e === "request_times") {
        for (var i = 0; i < record.apilst.length; i++) {
          lst.push({"Api": record.apilst[i].Api, "Value": record.apilst[i].ReqNum})
        }
      }

      PubSub.publish(Topic.ReportSelect, {
        Chart : e,
        ApiList : lst
      })
    }
  };

  render() {
    const data = this.state.data;
    return (
      <div>
        <Table columns={this.state.columns} dataSource={data} />
        <Tabs defaultActiveKey="Tree">
          <TabPane tab="Tree" key="Tree">
            <ApiChart />
          </TabPane>
        </Tabs>
      </div>
    );
  }
}
