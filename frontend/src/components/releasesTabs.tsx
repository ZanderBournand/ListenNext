'use client'

import ReleasePreview from "./preview";
import { Tabs } from "flowbite-react";
import ReleasesGrid from "./relasesGrid";
import Sidebar from "./sidebar";
import { useEffect, useState } from "react";

const periods = ["past", "week", "month", "extended"];

export default function ReleasesTabs() {
    const [releaseType, setReleaseType] = useState<any>("album")
    const [period, setPeriod] = useState<any>("week")

    const handleTabSwitch = (e: any) => {
        setPeriod(periods[e])
    }
    
    return (
        <div className="flex flex-row max-w-7xl mx-auto px-6">
            <Sidebar period={period} setReleaseType={setReleaseType}/>
            <Tabs.Group className="w-full" onActiveTabChange={handleTabSwitch}>
                <Tabs.Item title="Recent">
                    <ReleasesGrid releaseType={releaseType} period="past"/>
                </Tabs.Item>
                <Tabs.Item title="This Week" active>
                    <ReleasesGrid releaseType={releaseType} period="week"/>
                </Tabs.Item>
                <Tabs.Item title="Next Month">
                    <ReleasesGrid releaseType={releaseType} period="month"/>
                </Tabs.Item>
                <Tabs.Item title="Too Long...">
                    <ReleasesGrid releaseType={releaseType} period="extended"/>
                </Tabs.Item>
            </Tabs.Group>
        </div>
    )
}