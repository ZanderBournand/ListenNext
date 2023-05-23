'use client'

import { useQuery } from "@apollo/client";
import { queryReleasesCount } from "@/util/queries";
import { Sidebar as FlowbiteSidebar } from "flowbite-react";
import '../animations/loading.css';
import classNames from "classnames";
import { useSidebarContext } from "@/context/sidebarContext";

export default function Sidebar({period, releaseType, setReleaseType}: any) {    
    const {data, loading} = useQuery<any>(queryReleasesCount, {
        fetchPolicy: "network-only"
    })
    const { isOpenOnSmallScreens: isSidebarOpenOnSmallScreens } =
    useSidebarContext();
    
    const handleFilterSwitch = (filerType: any) => {
        setReleaseType(filerType)
    }
    
    return (
        <div className={classNames(
                "flex flex-col w-auto bg-gray-100/100 mt-40 rounded-xl h-96 absolute top-24 z-10 md:sticky md:!block md:bg-gray-100/50 md:mt-0",
                {
                    hidden: !isSidebarOpenOnSmallScreens,
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
    )
}