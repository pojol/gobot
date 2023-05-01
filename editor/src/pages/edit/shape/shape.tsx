import { Node } from "@antv/x6";
import ThemeType from "@/constant/constant";
import { NodeTy } from "@/constant/node_type";
import RootNode from "./root";

import ParallelLightNode, { ParallelDarkNode } from "./parallel";
import ActionLightNode, { ActionDarkNode } from "./action";
import ConditionLightNode, { ConditionDarkNode } from "./condition";
import LoopLightNode, { LoopDarkNode } from "./loop";
import RootLightNode, { RootDarkNode } from "./root";
import SelectorLightNode, { SelectorDarkNode } from "./selector";
import SequenceLightNode, { SequenceDarkNode } from "./sequence";
import WaitLightNode, { WaitDarkNode } from "./wait";


export function GetNode(ty: string, parm: any): Node {

    console.info("get node", localStorage.theme, ty)

    if (localStorage.theme === ThemeType.Dark) {

        switch (ty) {
            case NodeTy.Parallel:
                return new ParallelDarkNode(parm)
            case NodeTy.Action:
                return new ActionDarkNode(parm)
            case NodeTy.Condition:
                return new ConditionDarkNode(parm)
            case NodeTy.Loop:
                return new LoopDarkNode(parm)
            case NodeTy.Root:
                return new RootDarkNode(parm)
            case NodeTy.Selector:
                return new SelectorDarkNode(parm)
            case NodeTy.Sequence:
                return new SequenceDarkNode(parm)
            case NodeTy.Wait:
                return new WaitDarkNode(parm)
            default:
                return new ActionDarkNode(parm)
        }

    } else if (localStorage.theme === ThemeType.Light) {

        switch (ty) {
            case NodeTy.Parallel:
                return new ParallelLightNode(parm)
            case NodeTy.Action:
                return new ActionLightNode(parm)
            case NodeTy.Condition:
                return new ConditionLightNode(parm)
            case NodeTy.Loop:
                return new LoopLightNode(parm)
            case NodeTy.Root:
                return new RootLightNode(parm)
            case NodeTy.Selector:
                return new SelectorLightNode(parm)
            case NodeTy.Sequence:
                return new SequenceLightNode(parm)
            case NodeTy.Wait:
                return new WaitLightNode(parm)
            default:
                return new ActionLightNode(parm)
        }

    }

    return new RootNode()
}