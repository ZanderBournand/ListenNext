'use client'

import { queryReleasesCount } from "@/util/queries";
import { useSuspenseQuery } from "@apollo/experimental-nextjs-app-support/ssr";
import { Sidebar as FlowbiteSidebar } from "flowbite-react";

export default function Sidebar({period, setReleaseType}: any) {
    const {data} = useSuspenseQuery<any>(queryReleasesCount)
    
    const handleFilterSwitch = (filerType: any) => {
        setReleaseType(filerType)
    }
    
    return (
        <div className="flex flex-col w-auto bg-gray-100/50 rounded-xl h-96 sticky top-24">
            <h3 className="pl-8 py-4 text-xl font-semibold">Filter</h3>
            <div className="border-b border-slate-400 mx-6"></div>
            <FlowbiteSidebar aria-label="Default sidebar example" className="pt-6">
                <FlowbiteSidebar.Items>
                <FlowbiteSidebar.ItemGroup>
                    <FlowbiteSidebar.Item 
                        href="#" 
                        label={data?.allReleasesCount?.[period]?.all}
                        abelColor="alternative" 
                        onClick={() => handleFilterSwitch("all")}
                    >
                    All
                    </FlowbiteSidebar.Item>
                    <FlowbiteSidebar.Item 
                        href="#" 
                        label={data?.allReleasesCount?.[period]?.albums}
                        abelColor="alternative" 
                        onClick={() => handleFilterSwitch("album")}
                    >
                    Albums
                    </FlowbiteSidebar.Item>
                    <FlowbiteSidebar.Item 
                        href="#" 
                        label={data?.allReleasesCount?.[period]?.singles}
                        abelColor="alternative" 
                        onClick={() => handleFilterSwitch("single")}
                    >
                    Singles
                    </FlowbiteSidebar.Item>
                </FlowbiteSidebar.ItemGroup>
                </FlowbiteSidebar.Items>
            </FlowbiteSidebar>
        </div>
    )
}