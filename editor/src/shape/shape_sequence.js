import { Shape } from "@antv/x6";

class SequenceNode extends Shape.Rect {}

SequenceNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#f5f5f5",
      stroke: "#2db7f5",
      strokeWidth: 1,
      borderRadius: 2,
    },
    label: {
      text: "Sequence",
    },
    type : "SequenceNode"
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

export default SequenceNode;
