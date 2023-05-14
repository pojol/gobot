import { Shape } from "@antv/x6";

class ActionLightNode extends Shape.Rect { }

ActionLightNode.config({
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
    label : {
      text: "ActionNode"
    },
    type: { name: "ActionNode" }
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


export class ActionDarkNode extends Shape.Rect { }

ActionDarkNode.config({
  attrs: {
    root: {
      magnet: true,
    },
    body: {
      fill: "#20262E",
      stroke: "#035397",
      strokeWidth: 1,
      borderRadius: 2,
    },
    label : {
      fill : "#fff",
      text: "ActionNode"
    },
    type: { name: "ActionNode" }
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

export default ActionLightNode;
