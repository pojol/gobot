import * as React from "react";

import { TaskTimer } from 'tasktimer';
import axios from "axios";
import { useDispatch, useSelector } from 'react-redux';
import { setHeartColor } from "@/models/config";

// offline
// heartColor = #DCDCDC
// online
// heartColor = #389e0d

const heart = async (): Promise<boolean> => {

    let heartStatus = false
    let ping: number | null = null;
    const startTime = Date.now(); // 记录发送请求的时间戳

    try {
        const res = await axios.get(localStorage.remoteAddr + "/health", {});
        if (res.status === 200) {
            heartStatus = true
            const endTime = Date.now(); // 记录收到响应的时间戳
            ping = endTime - startTime; // 计算时间差，即ping值
            console.info("health", ping, "ms")
        }
    } catch (err) {
        console.info("health", err)
    }

    return heartStatus
}

export default function HeartTask() {
    const dispatch = useDispatch()
    
    let callback = async () => {
        let res = await heart()
        if (res) {
            dispatch(setHeartColor("#389e0d"))
        } else {
            dispatch(setHeartColor("#BDCDD6"))
        }
    }

    callback()

    const timer = new TaskTimer(5000);
    timer.on('tick', () => {
        callback()
    });
    timer.start();


    return (<div></div>);

}