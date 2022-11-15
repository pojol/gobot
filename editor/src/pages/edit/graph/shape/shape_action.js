import { Shape } from "@antv/x6";

class ActionNode extends Shape.Rect {}

ActionNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#f5f5f5",
      stroke: "#035397",
      strokeWidth: 1,
      borderRadius: 2,
    },
    type : "ActionNode"
  },
  width: 40,
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

export default ActionNode;
