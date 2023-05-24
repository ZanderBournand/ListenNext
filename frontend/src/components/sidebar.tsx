'use client'

import { useQuery } from "@apollo/client";
import { queryReleasesCount } from "@/util/queries";
import { Sidebar as FlowbiteSidebar } from "flowbite-react";
import '../animations/loading.css';
import classNames from "classnames";
import { useSidebarContext } from "@/context/sidebarContext";
import { useEffect, useRef } from "react";

export default function Sidebar({period, releaseType, setReleaseType}: any) {    
    const {data, loading} = useQuery<any>(queryReleasesCount, {
        fetchPolicy: "network-only"
    })
    const { isOpenOnSmallScreens, setOpenOnSmallScreens, closeTransition, openTransition, setCloseTransition } =
    useSidebarContext();
    const sidebarRef = useRef<HTMLDivElement>(null);
    
    const handleFilterSwitch = (filerType: any) => {
        setReleaseType(filerType)
    }

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
          if (sidebarRef.current && !sidebarRef.current.contains(event.target as Node)) {
            setOpenOnSmallScreens(!isOpenOnSmallScreens)
            setCloseTransition(true)
            setTimeout(() => {
                setCloseTransition(false)
            }, 500)
          }
        };
    
        if (isOpenOnSmallScreens) {
          document.addEventListener("mousedown", handleClickOutside);
        }
    
        return () => {
          document.removeEventListener("mousedown", handleClickOutside);
        };
      }, [isOpenOnSmallScreens]);
      
    
    return (
        <>
         {isOpenOnSmallScreens && (
            <div className="fixed top-0 left-0 right-0 bottom-0 bg-black opacity-50 z-40 md:hidden"></div>
        )}
        <div ref={sidebarRef} className={classNames(
            "flex flex-col w-auto bg-gray-100/100 fixed top-0 left-0 z-40 w- h-screen transition-all transform  md:h-96 md:top-24 md:z-10 md:rounded-xl md:sticky md:!block md:bg-gray-100/50 md:mt-0 md:translate-x-0 md:duration-0",
            {
                "duration-0 -translate-x-full block": !isOpenOnSmallScreens && !closeTransition,
                "duration-300 -translate-x-full block": !isOpenOnSmallScreens && closeTransition,
                "duration-300 translate-x-0 block": isOpenOnSmallScreens && openTransition,
                "duration-0 translate-x-0 block": isOpenOnSmallScreens && !openTransition
            },
        )}
        >
            <h3 className="pl-8 py-4 text-xl font-semibold">Filter</h3>
            <div className="border-b border-slate-400 mx-6"></div>
            <FlowbiteSidebar aria-label="Default sidebar example" className="pt-6">
                <FlowbiteSidebar.Items>
                <FlowbiteSidebar.ItemGroup>
                    <FlowbiteSidebar.Item 
                        label={(loading) ? 
                            <div className="dot-flashing"></div>
                            :
                            data?.allReleasesCount?.[period]?.all
                        }
                        onClick={() => handleFilterSwitch("all")}
                    >
                        <span className={releaseType === 'all' ? 'text-blue-500': ''}>All</span>
                    </FlowbiteSidebar.Item>
                    <FlowbiteSidebar.Item 
                        label={(loading) ? 
                            <div className="dot-flashing"></div>
                            :
                            data?.allReleasesCount?.[period]?.albums
                        }
                        onClick={() => handleFilterSwitch("album")}
                    >
                        <span className={releaseType === 'album' ? 'text-blue-500': ''}>Albums</span>
                    </FlowbiteSidebar.Item>
                    <FlowbiteSidebar.Item 
                        label={(loading) ? 
                            <div className="dot-flashing"></div>
                            :
                            data?.allReleasesCount?.[period]?.singles
                        }
                        onClick={() => handleFilterSwitch("single")}
                    >
                        <span className={releaseType === 'single' ? 'text-blue-500': ''}>Singles</span>
                    </FlowbiteSidebar.Item>
                </FlowbiteSidebar.ItemGroup>
                </FlowbiteSidebar.Items>
            </FlowbiteSidebar>
        </div>
        </>
    )
}