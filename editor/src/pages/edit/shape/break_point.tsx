import { Shape } from "@antv/x6";

class BPLightNode extends Shape.Ellipse { }

BPLightNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#EDBE8C",
      stroke: "#EDBE8C",
      strokeWidth: 1,
    },
    type: { name: "BreakPointNode" }
  },
  width: 10,
  height: 10,
});


export class BPDarkNode extends Shape.Ellipse { }

BPDarkNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#EDBE8C",
      stroke: "#EDBE8C",
      strokeWidth: 1,
    },
    type: { name: "BreakPointNode" }
  },
  width: 10,
  height: 10,
});

export default BPLightNode;
