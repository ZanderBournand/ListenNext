"use client";

import type { PropsWithChildren } from "react";
import { createContext, useContext, useEffect, useState } from "react";

interface SidebarContextProps {
  isOpenOnSmallScreens: boolean;
  isPageWithSidebar: boolean;
  closeTransition: boolean,
  openTransition: boolean,
  // eslint-disable-next-line no-unused-vars
  setOpenOnSmallScreens: (isOpen: boolean) => void;
  setCloseTransition: (isCloseTransition: boolean) => void
  setOpenTransition: (isOpenTransition: boolean) => void
}

const SidebarContext = createContext<SidebarContextProps>(undefined!);

export function SidebarProvider({
  children,
}: PropsWithChildren<Record<string, unknown>>) {
  const location = isBrowser() ? window.location.pathname : "/";
  const [isOpen, setOpen] = useState(false);
  const [closeTransition, setCloseTransition] = useState(false)
  const [openTransition, setOpenTransition] = useState(false)

  // Close Sidebar on page change on mobile
  useEffect(() => {
    if (isSmallScreen()) {
      setOpen(false);
    }
  }, [location]);

  // Close Sidebar on mobile tap inside main content
  useEffect(() => {
    function handleMobileTapInsideMain(event: MouseEvent) {
      const main = document.querySelector("main");
      const isClickInsideMain = main?.contains(event.target as Node);

      if (isSmallScreen() && isClickInsideMain) {
        setOpen(false);
      }
    }

    document.addEventListener("mousedown", handleMobileTapInsideMain);
    return () => {
      document.removeEventListener("mousedown", handleMobileTapInsideMain);
    };
  }, []);

  return (
    <SidebarContext.Provider
      value={{
        isOpenOnSmallScreens: isOpen,
        isPageWithSidebar: true,
        setOpenOnSmallScreens: setOpen,
        closeTransition: closeTransition,
        openTransition: openTransition,
        setCloseTransition: setCloseTransition,
        setOpenTransition: setOpenTransition
      }}
    >
      {children}
    </SidebarContext.Provider>
  );
}

function isBrowser(): boolean {
  return typeof window !== "undefined";
}

function isSmallScreen(): boolean {
  return isBrowser() && window.innerWidth < 300;
}

export function useSidebarContext(): SidebarContextProps {
  const context = useContext(SidebarContext) as SidebarContextProps | undefined;

  if (!context) {
    throw new Error(
      "useSidebarContext should be used within the SidebarContext provider!"
    );
  }

  return context;
}
