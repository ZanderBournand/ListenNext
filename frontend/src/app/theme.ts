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
        default: "flex flex-row items-center justify-center max-w-6xl mx-10 pr-16 pb-4 w-full"
      },
      tabitem: {
        base: "w-full py-2",
        styles: {
          default: {
            active: {
              on: "text-c3 border-x border-t rounded-t-lg font-mediun sm:text-md md:text-lg",
              off: "border-b hover:bg-gray-100/50 rounded-t-lg text-md sm:text-md md:text-lg"
            }
          }
        }
      }
    }
  },
  footer: {
    brand: {
      span: "text-2xl font-semibold text-c1"
    },
    root: {
      base: "",
    }
  },
  sidebar: {
    root: {
      base: "md:w-44 lg:w-64",
      inner: "",
    },
  },
  badge: {
    root: {
      color: {
        ["light"]: "bg-blue-200/50 text-c3",
        ["medium"]: "bg-blue-500/25 text-c1",
        ["dark"]: "bg-blue-700/25 text-blue-950",
      }
    }
  },
  alert: {
    closeButton: {
      color: {
        info: "bg-blue-200/0 hover:bg-blue-200/25"
      },
      icon: "h-6 w-6 fill-c3"
    }
  }
};
