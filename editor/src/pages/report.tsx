import React, { useState, useEffect } from 'react';

import { Table, Tag, Tabs, message, Row } from "antd";
import ApiChart from "./chart/chart_tree"
import Api from "@/constant/api";
const { Post } = require("../utils/request");

import PubSub from "pubsub-js";
import Topic from "../constant/topic";

const { TabPane } = Tabs;

interface ReportApiInfo {
  Api: string,
  ConsumeNum: number,
  ErrNum: number,
  ReqNum: number,
  ReqSize: number,
  ResSize: number,
}

interface ReportViewInfo {
  name : string,
  apilst: Array<ReportApiInfo>,
  botnum: number,
  duration: number,
  errors: number,
  key: string,
  reqnum: number,
  time: string,
  tps: number,
}

export default function TestReport() {

  const columns = [
    {
      title: "Name",
      dataIndex: "name",
      key: "name",
    },
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
    }
  ]

  const [reports, setReports] = useState([]);
  const [row, setRow] = useState<ReportViewInfo>();

  useEffect(() => {
    console.info("refresh reports")
    refresh()
  }, [])

  useEffect(()=>{
    clickRow("avg_request_time_ms")
  }, [row])

  const fillData = (info: any) => {

    var newdata = []
    for (var i = 0; i < info.length; i++) {

      var date = new Date(info[i].BeginTime * 1000);
      var convdataTime = date.getFullYear() + '-' + (date.getMonth()+1).toString() + '-' + date.getDate() + ' ' + date.getHours() + ':' + date.getMinutes() + ':' + date.getSeconds();

      newdata.push({
        key: info[i].ID,
        name: info[i].Name,
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

  const clickRow = (ty: string) => {

    let lst = []
    let rows = row
    if (!rows) {
      return
    }

    if (ty === "avg_request_time_ms") {
      for (var i = 0; i < rows.apilst.length; i++) {
        lst.push({ "Api": rows.apilst[i].Api, "Value": rows.apilst[i].ConsumeNum })
      }
    } else if (ty === "request_times") {
      for (i = 0; i < rows.apilst.length; i++) {
        lst.push({ "Api": rows.apilst[i].Api, "Value": rows.apilst[i].ReqNum })
      }
    }

    PubSub.publish(Topic.ReportSelect, {
      Chart: ty,
      ApiList: lst
    })

  }

  const tableClick = (activeKey: string) => {
    console.info(activeKey)

    if (activeKey === "Latency") {
      clickRow("avg_request_time_ms")
    } else {
      clickRow("request_times")
    }
  }

  return (
    <div>
      <Table columns={columns} dataSource={reports}
        rowSelection={{
          type: "radio",
          ...{
            onChange: (selectedRowKeys: any, selectedRows: any) => {
              //console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
              let apis = new Array<ReportApiInfo>() 
              if (selectedRows[0].apilst !== null) {
                console.info("api length", selectedRows[0].apilst.length)
                for (var i = 0; i < selectedRows[0].apilst.length; i++) {
                  apis.push({
                    Api: selectedRows[0].apilst[i].Api,
                    ConsumeNum: selectedRows[0].apilst[i].ConsumeNum,
                    ErrNum: selectedRows[0].apilst[i].ErrNum,
                    ReqNum: selectedRows[0].apilst[i].ReqNum,
                    ReqSize: selectedRows[0].apilst[i].ReqSize,
                    ResSize: selectedRows[0].apilst[i].ResSize,
                  })
                }
              }

              setRow({
                apilst: apis,
                name: selectedRows[0].name,
                botnum: selectedRows[0].botnum,
                duration: selectedRows[0].duration,
                errors: selectedRows[0].errors,
                key: selectedRows[0].key,
                reqnum: selectedRows[0].reqnum,
                time: selectedRows[0].time,
                tps: selectedRows[0].tps,
              })
            }
          },
        }}
        onRow={(record) => {
          return {
            onMouseDown: (e) => {
              if (e.target.type !== "radio") {
                e.currentTarget.getElementsByClassName("ant-radio-wrapper")[0].click()
              }
            },// 点击行
          };
        }} />
      <Tabs defaultActiveKey="Latency" onTabClick={tableClick}>
        <TabPane tab="Latency" key="Latency">
          <ApiChart />
        </TabPane>
        <TabPane tab="Frequency" key="Frequency">
          <ApiChart />
        </TabPane>
      </Tabs>
    </div>
  );
}
