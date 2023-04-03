import React, { useState, useEffect } from 'react';

import { Table, Tag, Tabs, message } from "antd";
import ApiChart from "./chart/chart_tree"
import Api from "@/constant/api";
const { Post } = require("../utils/request");

import PubSub from "pubsub-js";
import Topic from "../constant/topic";

const { TabPane } = Tabs;

export default function TestReport() {

  const columns = [
    {
      title: "Time",
      dataIndex: "time",
      key: "time",
      defaultSortOrder: undefined,
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
      render: (tags: any, record: any) => (
        <>
          {tags.map((tag: string) => {
            let color = 'green';
            if (tag === 'avg_request_time_ms') {
              color = 'volcano';
            } else if (tag === 'request_times') {
              color = 'geekblue';
            }
            return (
              <Tag color={color} key={tag} onClick={() => clickTag(tag, record)}>
                {tag}
              </Tag>

            );
          })}
        </>
      ),
    }

  ]

  const [reports, setReports] = useState([]);

  useEffect(() => {
    console.info("refresh reports")
    refresh()
  },[])

  const fillData = (info: any) => {

    var newdata = []
    for (var i = 0; i < info.length; i++) {

      var date = new Date(info[i].BeginTime * 1000);
      var convdataTime = date.getFullYear() + '-' + date.getMonth() + '-' + date.getDate() + ' ' + date.getHours() + ':' + date.getMinutes() + ':' + date.getSeconds();

      newdata.push({
        key: info[i].ID,
        time: convdataTime,
        tps: info[i].Tps,
        duration: info[i].Dura,
        botnum: info[i].BotNum,
        reqnum: info[i].ReqNum,
        errors: info[i].ErrNum,
        charts: ["avg_request_time_ms", "request_times"],
        apilst: info[i].ApiInfoLst,
      })
    }

    setReports(newdata)
  }

  const refresh = () => {
    setReports([])
    Post(localStorage.remoteAddr, Api.ReportInfo, {}).then((json: any) => {
      if (json.Code !== 200) {
        message.error("run fail:" + String(json.Code) + " msg: " + json.Msg);
      } else {
        if (json.Body.Info) {
          fillData(json.Body.Info)
        }
      }
    });
  }



  const clickTag = (e: any, record: any) => {

    if (record.apilst) {
      let lst = []

      if (e === "avg_request_time_ms") {
        for (var i = 0; i < record.apilst.length; i++) {
          lst.push({ "Api": record.apilst[i].Api, "Value": record.apilst[i].ConsumeNum })
        }
      } else if (e === "request_times") {
        for (i = 0; i < record.apilst.length; i++) {
          lst.push({ "Api": record.apilst[i].Api, "Value": record.apilst[i].ReqNum })
        }
      }


      PubSub.publish(Topic.ReportSelect, {
        Chart: e,
        ApiList: lst
      })

    }
  };

  return (
    <div>
      <Table columns={columns} dataSource={reports} />
      <Tabs defaultActiveKey="Tree">
        <TabPane tab="Tree" key="Tree">
          <ApiChart />
        </TabPane>
      </Tabs>
    </div>
  );
}
