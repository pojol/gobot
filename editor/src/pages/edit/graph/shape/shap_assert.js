import { Shape } from "@antv/x6";

class AssertNode extends Shape.Polygon {}

AssertNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#f5f5f5",
      stroke: "#E04D01",
      strokeWidth: 1,
      refPoints: "0,5 2.5,0 5,5",
    },
    label: {
      text: "Assert",
    },
    type: "AssertNode",
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

export default AssertNode;
