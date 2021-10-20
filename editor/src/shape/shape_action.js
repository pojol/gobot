import { Shape } from "@antv/x6";
import { NodeTy } from "../model/node_type";

class ActionNode extends Shape.Rect {}

ActionNode.config({
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
      text: "HTTP"
    },
    type : NodeTy.Action
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
