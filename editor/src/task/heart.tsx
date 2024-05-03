import React, { useEffect, useState } from 'react';

import { TaskTimer } from 'tasktimer';
import axios from "axios";
import { useDispatch, useSelector } from 'react-redux';
import { setHeartColor } from "@/models/config";


const heart = async (): Promise<number> => {

    let ping: number = 0;
    const startTime = Date.now(); // 记录发送请求的时间戳

    try {
        const res = await axios.get(localStorage.remoteAddr + "/health", {});
        if (res.status === 200) {
            
            const endTime = Date.now(); // 记录收到响应的时间戳
            ping = endTime - startTime; // 计算时间差，即ping值

        }
    } catch (err) {
        console.info("health", err)
    }
    
    return ping

}

const HeartTask = () => {
    const dispatch = useDispatch();
    const [timer, setTimer] = useState<TaskTimer | null>(null);
  
    useEffect(() => {
      let mounted = true;
  
      const callback = async () => {
        let ping = await heart();
        if (mounted) {
          if (ping !== 0) {
            dispatch(setHeartColor(ping + " ms"));
          } else {
            dispatch(setHeartColor(""));
          }
        }
      };
  
      callback();
  
      const newTimer = new TaskTimer(2000);
      newTimer.on('tick', callback);
      newTimer.start();
  
      setTimer(newTimer);
  
      return () => {
        // 在组件卸载时清除定时器
        if (timer) {
          timer.stop();
        }
        setTimer(null);
        mounted = false;
      };
    }, [dispatch]);
  
    return <div></div>;
  }
  
  export default HeartTask;
  