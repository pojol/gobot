import { Shape } from "@antv/x6";
import Constant from "../../../../constant/constant";

class RootNode extends Shape.Rect { }

RootNode.config({
  x: Constant.GraphWidth / 2,
  y: Constant.GraphHeight / 2,
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#f5f5f5",
      stroke: "#336666",
      strokeWidth: 2,
      borderRadius: 4,
    },
    label: {
      text: "root"
    },
    type: "RootNode"
  },
  width: 50,
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

export default RootNode;