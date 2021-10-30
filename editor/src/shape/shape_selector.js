import { Shape } from "@antv/x6";


class SelectorNode extends Shape.Rect {}

SelectorNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#f5f5f5",
      stroke: "#f50",
      strokeWidth: 1,
      borderRadius: 2,
    },
    label: {
      text: "Selector"
    },
    type : "SelectorNode"
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

export default SelectorNode;
