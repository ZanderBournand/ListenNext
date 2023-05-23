'use client'

import { Tabs } from "flowbite-react";
import ReleasesGrid from "./relasesGrid";
import Sidebar from "./sidebar";
import { useEffect, useState } from "react";
import { SidebarProvider, useSidebarContext } from "../context/sidebarContext";

const periods = ["past", "week", "month", "extended"];

export default function ReleasesTabs() {
    const [releaseType, setReleaseType] = useState<any>("album")
    const [period, setPeriod] = useState<any>("week")
    const [windowWidth, setWindowWidth] = useState<any>(0);
    const { isOpenOnSmallScreens, isPageWithSidebar, setOpenOnSmallScreens } =
    useSidebarContext();
    const [tabTitles, setTabTitles] = useState<any>(['Recent', 'This Week', 'Next Month', 'Too Long...']);

    useEffect(() => {
      setTabTitles(getTabTitles());
  
      const handleResize = () => {
        setWindowWidth(window.innerWidth);
        setTabTitles(getTabTitles());
      };
  
      window.addEventListener('resize', handleResize);
  
      return () => {
        window.removeEventListener('resize', handleResize);
      };
    }, []);

    const getTabTitles = () => {
        const currentWidth = typeof window !== 'undefined' ? window.innerWidth : 0;
    
        if (currentWidth < 600) {
          return ['Recent', 'Week', 'Month', 'Long']
        } else {
          return ['Recent', 'This Week', 'Next Month', 'Too Long...'];
        }
      };

    const handleTabSwitch = (e: any) => {
        setPeriod(periods[e])
    }
    
    return (
        <div className="flex flex-row max-w-7xl mx-auto px-6">
            {isPageWithSidebar && (
            <button
                aria-controls="sidebar"
                aria-expanded="true"
                className="h-8 absolute cursor-pointer rounded mt-2 text-gray-600 hover:bg-gray-100 hover:text-gray-900 focus:bg-gray-100 focus:ring-2 focus:ring-gray-100 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white dark:focus:bg-gray-700 dark:focus:ring-gray-700 md:hidden"
                onClick={() => setOpenOnSmallScreens(!isOpenOnSmallScreens)}
            >
                {isOpenOnSmallScreens ? (
                <svg
                    className="h-6 w-6"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <path
                    fillRule="evenodd"
                    d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                    clipRule="evenodd"
                    ></path>
                </svg>
                ) : (
                <svg
                    className="h-7 w-7"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <path
                    fillRule="evenodd"
                    d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h6a1 1 0 110 2H4a1 1 0 01-1-1zM3 15a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"
                    clipRule="evenodd"
                    ></path>
                </svg>
                )}
            </button>
            )}
            <Sidebar period={period} releaseType={releaseType} setReleaseType={setReleaseType}/>
            <div className="w-full ml-4 sm:ml-4 md:ml-0">
                <Tabs.Group className="w-full" onActiveTabChange={handleTabSwitch}>
                {tabTitles.map((title: any, index: any) => (
                    <Tabs.Item key={index} title={title} active={index === 1}>
                    <ReleasesGrid releaseType={releaseType} period={period} />
                    </Tabs.Item>
                ))}
                </Tabs.Group>
            </div>
        </div>
    )
}