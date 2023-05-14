import { Shape } from "@antv/x6";

class SequenceLightNode extends Shape.Rect { }

SequenceLightNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#f5f5f5",
      stroke: "#A7C1D5",
      strokeWidth: 1,
      borderRadius: 2,
      rx: 3, // 圆角矩形
      ry: 3,
    },
    label: {
      text: "Sequence",
    },
    type: { name: "SequenceNode" }
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


export class SequenceDarkNode extends Shape.Rect { }

SequenceDarkNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#20262E",
      stroke: "#A7C1D5",
      strokeWidth: 1,
      borderRadius: 2,
      rx: 3, // 圆角矩形
      ry: 3,
    },
    label: {
      fill : "#fff",
      text: "Sequence",
    },
    type: { name: "SequenceNode" }
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

export default SequenceLightNode;
