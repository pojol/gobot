import React, { useState } from 'react';

import { Link, Outlet, useAppData, useLocation } from 'umi';
import { ProLayout } from '@ant-design/pro-layout';
import type { ProSettings } from '@ant-design/pro-components';

import { Radio, ConfigProvider, theme, Modal, Tag, Input } from 'antd';
import type { RadioChangeEvent } from 'antd';
import ThemeType from '@/constant/constant';
import PubSub from "pubsub-js";

import { Provider, useSelector, useDispatch } from 'react-redux';
import store, { RootState } from '../models/store';

import {
  GithubFilled,
  ThunderboltTwoTone,
} from '@ant-design/icons';

import logo from '../assets/gobot.png';
import darklog from '../assets/gobot-dark.png'
import Topic from '@/constant/topic';

import { setThemeValue, setModalOpen } from '@/models/config';

var state = {
  theme: theme.defaultAlgorithm
}

function getState(): any {
  if (localStorage.theme === undefined || localStorage.theme === ThemeType.Light) {
    localStorage.theme = ThemeType.Light
    return theme.defaultAlgorithm
  } else if (localStorage.theme === ThemeType.Dark) {
    return theme.darkAlgorithm
  }
}

function getLog() : any {
  if (localStorage.theme === undefined || localStorage.theme === ThemeType.Light) {
    return logo
  } else if (localStorage.theme === ThemeType.Dark) {
    return darklog
  }
}

function resizeHandler() {
  let resizeTimer;
  if (!resizeTimer) {
    resizeTimer = setTimeout(() => {
      resizeTimer = null;
      PubSub.publish(Topic.WindowResize, {});
    }, 100);
  }
};

function Layout() {
  const { clientRoutes } = useAppData();
  const location = useLocation();
  const { heartColor, themeValue, modalOpen } = useSelector((state: RootState) => state.configSlice)
  const [address, setAddress] = useState("")
  const dispatch = useDispatch()

  const settings: ProSettings | undefined = {
    fixSiderbar: true,
    layout: 'top',
    splitMenus: true,
  };

  window.addEventListener("resize", resizeHandler, false);

  const themeChange = (e: RadioChangeEvent) => {
    dispatch(setThemeValue(e.target.value))

    if (e.target.value == ThemeType.Dark) {
      state.theme = theme.darkAlgorithm
      localStorage.theme = ThemeType.Dark
    } else if (e.target.value == ThemeType.Light) {
      state.theme = theme.defaultAlgorithm
      localStorage.theme = ThemeType.Light
    }

    console.info("change theme")
    PubSub.publish(Topic.ThemeChange, localStorage.theme)
  };

  const modalHandleOk = () => {
    console.info("set modal false")
    dispatch(setModalOpen(false))

    // 这里需要做一次检测
    if (address !== "" && address !== undefined) {
      localStorage.remoteAddr = address
    }

  }

  const modalHandleCancel = () => {
    dispatch(setModalOpen(false))
  }

  const modalConfigChange = (e: any) => {
    setAddress(e.target.value)
  }

  return (
    <ConfigProvider
      theme={
        {
          token: {
            colorPrimary: '#F0A04B',
          },
          algorithm: getState(),
        }
      }>
      <ProLayout
        route={clientRoutes[0]}
        location={location}
        logo={getLog()}
        title={'Gobot'}
        menuItemRender={(menuItemProps, defaultDom) => {
          if (menuItemProps.isUrl || menuItemProps.children) {
            return defaultDom;
          }
          if (menuItemProps.path && location.pathname !== menuItemProps.path) {
            return (
              <Link to={menuItemProps.path} target={menuItemProps.target}>
                {defaultDom}
              </Link>
            );
          }
          return defaultDom;
        }}{...settings}
        actionsRender={(props) => {

          var color = ""
          var desc = ""
          if (heartColor !== "") {
            color = "success"
            desc = heartColor
          } else {
            color = "#B6BBC4"
            desc = "Disconnected"
          }
 
          if (props.isMobile) return [];
          return [
            <Tag>v0.3.8</Tag>,
            <Radio.Group onChange={themeChange} value={themeValue} buttonStyle="solid" defaultValue={localStorage.theme} size={"small"}>
              <Radio.Button value={ThemeType.Dark}>Dark</Radio.Button>
              <Radio.Button value={ThemeType.Light}>Light</Radio.Button>
            </Radio.Group>,
            <Tag color={color}>{desc}</Tag>,
            <GithubFilled key="GithubFilled" twoToneColor='#eb2f96' onClick={function () { window.open("https://github.com/pojol/gobot"); }} />,
          ];

        }}
      >
        <Outlet />
      </ProLayout>

      <Modal
        open={modalOpen}
        onOk={modalHandleOk}
        onCancel={modalHandleCancel}
      >
        <Tag>e.g. http://178.128.113.58:30000 (Sample driver server address</Tag>
        <Input
          placeholder={"Input drive server address"}
          onChange={modalConfigChange}
        />
      </Modal>
    </ConfigProvider>
  );
}

export default () => {
  return (
    <Provider store={store}>
      <Layout />
    </Provider>
  )
}