import { CustomFlowbiteTheme } from "flowbite-react";

export const flowbiteTheme: CustomFlowbiteTheme = {
  navbar: {
    link: {
        base: "text-lg",
        active: {
            on: "text-c3",
            off: "text-c4 hover:text-c3"
        }
    },
  },
  tab: {
    tablist: {
      styles: {
        default: "flex flex-row max-w-6xl pl-2 pb-4"
      }
    }
  }
};
