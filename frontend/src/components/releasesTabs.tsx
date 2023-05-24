'use client'

import { Button, Tabs } from "flowbite-react";
import ReleasesGrid from "./relasesGrid";
import { useEffect, useState } from "react";
import Sidebar from "./sidebar";
import { useSidebarContext } from "@/context/sidebarContext";
import { Settings2 } from "lucide-react";

const periods = ["past", "week", "month", "extended"];

export default function ReleasesTabs() {
    const [releaseType, setReleaseType] = useState<any>("album")
    const [period, setPeriod] = useState<any>("week")
    const [tabTitles, setTabTitles] = useState<any>(['Recent', 'This Week', 'Next Month', 'Too Long...']);

    const { isOpenOnSmallScreens, setOpenOnSmallScreens, closeTransition, openTransition, setCloseTransition, setOpenTransition } =
    useSidebarContext();

    useEffect(() => {
      setTabTitles(getTabTitles());
  
      const handleResize = () => {
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
      <div className="flex flex-col">
      <div className="pl-14 pb-6 md:hidden">
        <Button
          color="light"
          pill={true}
          onClick={() => {
            setOpenOnSmallScreens(!isOpenOnSmallScreens)
            setOpenTransition(true)
            setTimeout(() => {
                setOpenTransition(false)
            }, 500)
          }}
        >
          <Settings2 className="h-4 w-4"/>
          <span className="pl-2">Filter</span>
        </Button>
      </div>
      <div className="flex flex-row w-full max-w-7xl mx-auto px-0 md:px-6">
            <Sidebar period={period} releaseType={releaseType} setReleaseType={setReleaseType}/>
            <div className="w-full">
                <Tabs.Group className="w-full" onActiveTabChange={handleTabSwitch}>
                {tabTitles.map((title: any, index: any) => (
                    <Tabs.Item key={index} title={title} active={index === 1}>
                    <ReleasesGrid releaseType={releaseType} period={period} />
                    </Tabs.Item>
                ))}
                </Tabs.Group>
            </div>
        </div>
      </div>
    )
}