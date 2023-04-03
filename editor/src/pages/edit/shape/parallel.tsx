import { Shape } from "@antv/x6";

class ParallelLightNode extends Shape.Rect { }

ParallelLightNode.config({
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
    type: { name: "ParallelNode" }
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


export class ParallelDarkNode extends Shape.Rect { }

ParallelDarkNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#20262E",
      stroke: "#B3CDAE",
      strokeWidth: 1,
    },
    label: {
      fill : "#fff",
      text: "Parallel",
    },
    type: { name: "ParallelNode" }
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

export default ParallelLightNode;
