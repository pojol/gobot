import { Shape } from "@antv/x6";

class ParallelNode extends Shape.Rect {}

ParallelNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#f5f5f5",
      stroke: "#B3CDAE",
      strokeWidth: 1,
    },
    label: {
      text: "Parallel",
    },
    type : "ParallelNode"
  },
  width: 30,
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

export default ParallelNode;
