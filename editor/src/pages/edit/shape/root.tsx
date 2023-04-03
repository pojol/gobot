import { Shape } from "@antv/x6";

class RootLightNode extends Shape.Rect { }

RootLightNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#f5f5f5",
      strokeWidth: 2,
      borderRadius: 4,
    },
    label: {
      text: "root",
    },
    type: { name: "RootNode" }
  },
  width: 50,
  height: 20,
  ports: {
    items: [
      {
        group: "out",
      },
    ],
    groups: {
      out: {
        position: {
          name: "bottom",
        },
        attrs: {
          portBody: {
            magnet: true,
            r: 5,
            fill: "#fff",
            stroke: "#3199FF",
            strokeWidth: 1,
          },
        },
      },
    },
  },
  portMarkup: [
    {
      tagName: "circle",
      selector: "portBody",
    },
  ],
});

export class RootDarkNode extends Shape.Rect { }

RootDarkNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#20262E",
      strokeWidth: 2,
      borderRadius: 4,
    },
    label: {
      fill : "#fff",
      text: "root",
    },
    type: { name: "RootNode" }
  },
  width: 50,
  height: 20,
  ports: {
    items: [
      {
        group: "out",
      },
    ],
    groups: {
      out: {
        position: {
          name: "bottom",
        },
        attrs: {
          portBody: {
            magnet: true,
            r: 5,
            fill: "#fff",
            stroke: "#3199FF",
            strokeWidth: 1,
          },
        },
      },
    },
  },
  portMarkup: [
    {
      tagName: "circle",
      selector: "portBody",
    },
  ],
});

export default RootLightNode;