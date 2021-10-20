import { Shape } from "@antv/x6";


class LoopNode extends Shape.Ellipse {}

LoopNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#f5f5f5",
      stroke: "#f50",
      strokeWidth: 1,
    },
    label: {
      text: "Loop"
    },
    type : "LoopNode"
  },
  width: 60,
  height: 30,
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

export default LoopNode;
