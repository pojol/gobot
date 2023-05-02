import { Shape } from "@antv/x6";


class SelectorLightNode extends Shape.Rect { }

SelectorLightNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#f5f5f5",
      stroke: "#AAD2D2",
      strokeWidth: 1,
      borderRadius: 2,
    },
    label: {
      text: "Selector"
    },
    type: { name: "SelectorNode" }
  },
  width: 60,
  height: 15,
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



export class SelectorDarkNode extends Shape.Rect { }

SelectorDarkNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#20262E",
      stroke: "#AAD2D2",
      strokeWidth: 1,
      borderRadius: 2,
    },
    label: {
      fill : "#fff",
      text: "Selector"
    },
    type: { name: "SelectorNode" }
  },
  width: 60,
  height: 15,
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
            fill: "#20262E",
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

export default SelectorLightNode;
