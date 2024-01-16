import { createAction, createSlice, PayloadAction } from "@reduxjs/toolkit";
import ThemeType from "@/constant/constant";

interface ConfigState {
    heartColor: string
    themeValue: string
    modalOpen: boolean
    runningTick: number
}

function initTheme(): string {
    if (localStorage.theme === undefined || localStorage.theme === ThemeType.Light) {
        return ThemeType.Light
    } else {
        return ThemeType.Dark
    }
}

function initModalOpen(): boolean {
    if (localStorage.remoteAddr === "" || localStorage.remoteAddr === undefined) {
        return true
    } else {
        return false
    }
}

const initialState: ConfigState = {
    heartColor: "",
    themeValue: initTheme(),
    modalOpen: initModalOpen(),
    runningTick: 2000,
}

const configSlice = createSlice({
    name: "debuginfo",
    initialState,
    reducers: {
        setHeartColor(state, action: PayloadAction<string>) {
            state.heartColor = action.payload
        },
        setThemeValue(state, action: PayloadAction<string>) {
            state.themeValue = action.payload
        },
        setModalOpen(state, action: PayloadAction<boolean>) {
            state.modalOpen = action.payload
        }
    }
})

export const { setHeartColor, setThemeValue, setModalOpen } = configSlice.actions;
export default configSlice;
