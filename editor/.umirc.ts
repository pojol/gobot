import { defineConfig } from "umi";

export default defineConfig({
  routes: [
    {
      path: '/',
      redirect: '/editor',
    },
    { path: "/editor/*", component: "editor", name: "editor", key:"editor" },
    { path: "/bots", component: "bots", name: "bots", key :"bots" },
    { path: "/running", component: "running", name: "running" , key:"running"},
    { path: "/report", component: "report", name: "report" , key:"report"},
    { path: "/prefab", component: "prefab", name: "prefab", key:"prefab" },
    { path: "/config", component: "config", name: "config", key:"config" },
    { path: "/docs", component: "docs", name: "docs" , key:"docs"},
  ],
  plugins: [
    '@umijs/plugins/dist/react-query',
    '@umijs/plugins/dist/model',
  ],
  model: {
  },
  npmClient: 'pnpm',
});
