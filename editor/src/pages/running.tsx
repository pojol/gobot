import {
  Table,
  message,
} from "antd";
import React, { useState, useEffect } from 'react';
const { Post } = require("../utils/request");
import Api from "../constant/api";
import { TaskTimer } from 'tasktimer';
import { useSelector } from 'react-redux';
import { RootState } from "@/models/store";


export default function Running() {

  const [runs, setRuns] = React.useState({});
  const [botLst, setBotLst] = React.useState([]);
  const { runningTick } = useSelector((state:RootState) => state.configSlice)

  const cloumns = [
    {
      title: "ID",
      dataIndex: "id",
      key: "id",
    },
    {
      title: "Name",
      dataIndex: "name",
      key: "name",
    },
    {
      title: "Current",
      dataIndex: "cur",
      key: "cur",
    },
    {
      title: "Target",
      dataIndex: "max",
      key: "max",
    },
    {
      title: "Errors",
      dataIndex: "errors",
      key: "errors",
    },
  ]

  useEffect(() => {
    refreshBotList();

    let callback = async () => {
      refreshBotList();
    }

    const timer = new TaskTimer(runningTick);
    timer.on('tick', () => {
      callback()
  });
  timer.start();

  return () => {
    timer.stop();
  }
  }, []);

  const fillBotList = (lst: any) => {
    if (lst) {
      var botlist = [];
      for (var i = 0; i < lst.length; i++) {
        botlist.push({
          id: lst[i].ID,
          key: lst[i].ID,
          name: lst[i].Name,
          cur: lst[i].Cur,
          max: lst[i].Max,
          errors: lst[i].Errors,
        });
      }
      setBotLst(botlist);
    }
  }


  const refreshBotList = () => {
    setBotLst([]);

    Post(localStorage.remoteAddr, Api.BotList, {}).then((json: any) => {
      if (json.Code !== 200) {
        message.error("run fail:" + String(json.Code) + " msg: " + json.Msg);
      } else {
        fillBotList(json.Body.Lst);
      }
    });
  }

  return (
    <div>
      <Table columns={cloumns} dataSource={botLst} />
    </div>
  )

}