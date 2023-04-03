import { Shape } from "@antv/x6";


class LoopLightNode extends Shape.Ellipse { }

LoopLightNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#f5f5f5",
      stroke: "#EDBE8C",
      strokeWidth: 1,
    },
    label: {
      text: "Loop",
    },
    type: { name: "LoopNode" }
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


export class LoopDarkNode extends Shape.Ellipse { }

LoopDarkNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#20262E",
      stroke: "#EDBE8C",
      strokeWidth: 1,
    },
    label: {
      fill : "#fff",
      text: "Loop",
    },
    type: { name: "LoopNode" }
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

export default LoopLightNode;
