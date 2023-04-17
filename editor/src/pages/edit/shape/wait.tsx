import { Shape } from "@antv/x6";


class WaitLightNode extends Shape.Rect { }

WaitLightNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#f5f5f5",
      stroke: "#ADA9D2",
      strokeWidth: 1,
      borderRadius: 2,
      rx: 3,
      ry: 3,
    },
    label: {
      text: "Wait"
    },
    type: { name: "WaitNode" }
  },
  width: 20,
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

export class WaitDarkNode extends Shape.Rect { }

WaitDarkNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#20262E",
      stroke: "#ADA9D2",
      strokeWidth: 1,
      borderRadius: 2,
      rx: 3,
      ry: 3,
    },
    label: {
      fill : "#fff",
      text: "Wait"
    },
    type: { name: "WaitNode" }
  },
  width: 20,
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


export default WaitLightNode;
