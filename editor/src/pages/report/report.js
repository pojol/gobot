import { Table, Tag, Tabs,message} from "antd";
import * as React from "react";
import ApiChart from "./chart_tree";
import PubSub from "pubsub-js";
import Topic from "../../constant/topic";
import { Post } from "../../utils/request";
import Api from "../../constant/api";

import moment from 'moment';
import lanMap from "../../locales/lan";

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

      var date = new Date(info[i].BeginTime*1000);
      var convdataTime = date.getFullYear() + '-'+date.getMonth()+'-'+date.getDate()+' '+date.getHours() + ':' + date.getMinutes() + ':' + date.getSeconds();


      newdata.push({
        key: info[i].ID,
        time:convdataTime ,
        tps:info[i].Tps,
        duration: info[i].Dura,
        botnum: info[i].BotNum,
        reqnum: info[i].ReqNum,
        errors : info[i].ErrNum,
        charts: ["avg_request_time_ms", "request_times"],
        apilst : info[i].ApiInfoLst,
      })
    }

    this.setState({data:newdata})
  }

  refresh() {
    this.setState({data:[]})
    Post(localStorage.remoteAddr, Api.ReportInfo, {}).then((json) => {
      if (json.Code !== 200) {
        message.error("run fail:" + String(json.Code) + " msg: " + json.Msg);
      } else {
        if (json.Body.Info){
          this.fillData(json.Body.Info)
        }
      }
    });
  }

  refresh_lan() {
    var lan = moment.locale()
    var columns = this.state.columns
    for (var i = 0; i < columns.length; i++) {
      if (columns[i].key === "time") {
        columns[i].title = lanMap["app.report.time"][lan]
      }else if (columns[i].key === "duration") {
        columns[i].title = lanMap["app.report.duration"][lan]
      }else if (columns[i].key === "botnum") {
        columns[i].title = lanMap["app.report.botnum"][lan]
      }else if (columns[i].key === "reqnum") {
        columns[i].title = lanMap["app.report.reqnum"][lan]
      }else if (columns[i].key === "errors") {
        columns[i].title = lanMap["app.report.errors"][lan]
      }
    }

    this.setState({columns: columns})
  }

  componentDidMount() {
    PubSub.subscribe(Topic.ReportUpdate, (topic, info) => {
      this.refresh();
    });

    PubSub.subscribe(Topic.LanuageChange, ()=>{
      this.refresh_lan()
    })

    this.refresh();
    this.refresh_lan();
  }

  clickTag = (e,record) => {

    if (record.apilst) {
      let lst = []

      if (e === "avg_request_time_ms") {
        for (var i = 0; i < record.apilst.length; i++) {
          lst.push({"Api": record.apilst[i].Api, "Value": record.apilst[i].ConsumeNum})
        }
      } else if (e === "request_times") {
        for (i = 0; i < record.apilst.length; i++) {
          lst.push({"Api": record.apilst[i].Api, "Value": record.apilst[i].ReqNum})
        }
      }

      console.info("report", e, lst)

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
