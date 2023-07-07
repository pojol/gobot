
import { createSlice, PayloadAction } from "@reduxjs/toolkit";


interface DebugInfoState {
    threadInfo: Array<ThreadInfo>
    metaInfo: string
    lock: boolean
}

const initialState: DebugInfoState = {
    metaInfo: "{}",
    threadInfo: new Array<ThreadInfo>(),
    lock: false
}

const debugInfoSlice = createSlice({
    name: "debuginfo",
    initialState,
    reducers: {
        setDebugInfo(state, action: PayloadAction<DebugInfoState>) {
            state.metaInfo = action.payload.metaInfo
            state.threadInfo = action.payload.threadInfo

            // 每次创建新的机器人时重制锁
            state.lock = action.payload.lock
        },
        setLock(state, action: PayloadAction<boolean>) {
            state.lock = action.payload
        }
    }
})

export const { setDebugInfo, setLock } = debugInfoSlice.actions;
export default debugInfoSlice;
