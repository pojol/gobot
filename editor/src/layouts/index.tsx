import React, { useState } from 'react';

import { Link, Outlet, useAppData, useLocation, useModel, history } from 'umi';
import { ProLayout } from '@ant-design/pro-layout';
import type { ProSettings } from '@ant-design/pro-components';

import { Radio, ConfigProvider, theme, Modal, Tag, Input } from 'antd';
import type { RadioChangeEvent } from 'antd';
import ThemeType from '@/constant/constant';
import PubSub from "pubsub-js";

import { Provider } from 'react-redux';
import store from '../models/store';

import {
  GithubFilled,
  ThunderboltTwoTone,
} from '@ant-design/icons';

import logo from '../assets/gobot.ico';
import Topic from '@/constant/topic';

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

export default function Layout() {
  const { clientRoutes } = useAppData();
  const location = useLocation();
  const { themeValue, setThemeValue } = useModel('theme')
  const { heartColor } = useModel('heartColor')
  const { open, setOpen } = useModel('modalConfig')
  const [address, setAddress] = useState("")

  const settings: ProSettings | undefined = {
    fixSiderbar: true,
    layout: 'top',
    splitMenus: true,
  };

  const themeChange = (e: RadioChangeEvent) => {
    setThemeValue(e.target.value);

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
    setOpen(false)

    // 这里需要做一次检测
    if (address !== "" && address !== undefined) {
      localStorage.remoteAddr = address
    }

  }

  const modalHandleCancel = () => {
    setOpen(false)
  }

  const modalConfigChange = (e: any) => {
    setAddress(e.target.value)
  }

  return (
    <Provider store={store}>
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
          logo={logo}
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

            if (props.isMobile) return [];
            return [
              <Radio.Group onChange={themeChange} value={themeValue} buttonStyle="solid" defaultValue={localStorage.theme} size={"small"}>
                <Radio.Button value={ThemeType.Dark}>Dark</Radio.Button>
                <Radio.Button value={ThemeType.Light}>Light</Radio.Button>
              </Radio.Group>,
              <ThunderboltTwoTone key="ThunderboltTwoTone" twoToneColor={heartColor} />,
              <GithubFilled key="GithubFilled" twoToneColor='#eb2f96' onClick={function () { window.open("https://github.com/pojol/gobot"); }} />,
            ];
          }}
        >
          <Outlet />
        </ProLayout>

        <Modal
          open={open}
          onOk={modalHandleOk}
          onCancel={modalHandleCancel}
        >
          <Tag>e.g. http://123.60.17.61:8888 (Sample driver server address</Tag>
          <Input
            placeholder={"Input drive server address"}
            onChange={modalConfigChange}
          />
        </Modal>
      </ConfigProvider>

    </Provider>

  );
}