'use client'

import ReleasePreview from "./preview";
import { Tabs } from "flowbite-react";
import ReleasesGrid from "./relasesGrid";

export default function ReleasesTabs({releaseType}: any) {
    console.log("RELEASE TYPE", releaseType)
    
    return (
        <Tabs.Group className="w-full">
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
    )
}