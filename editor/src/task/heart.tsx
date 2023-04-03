import * as React from "react";

import { useModel } from 'umi';

import { TaskTimer } from 'tasktimer';
import axios from "axios";


// offline
// heartColor = #DCDCDC
// online
// heartColor = #389e0d

const heart = async (): Promise<boolean> => {

    let heartStatus = false
    try {
        const res = await axios.get(localStorage.remoteAddr + "/health", {});
        if (res.status === 200) {
            heartStatus = true
        }
    } catch (err) {
        console.info("health", err)
    }

    return heartStatus
}

export default function HeartTask() {
    const { heartColor, setHeatColor } = useModel('heartColor')

    let callback = async () => {
        let res = await heart()
        if (res) {
            setHeatColor("#389e0d")
        } else {
            setHeatColor("#BDCDD6")
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