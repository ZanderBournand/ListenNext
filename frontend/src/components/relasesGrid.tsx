'use client'

import { gql } from "@apollo/client";
import ReleasePreview from "./preview";
import { useSuspenseQuery } from "@apollo/experimental-nextjs-app-support/ssr";
import { Tabs } from "flowbite-react";
import { useState } from "react";
import { queryAll } from '../util/queries'

export default function ReleasesGrid() {
    const {data} = useSuspenseQuery<any>(queryAll({ releaseType: 'album' }))
    const [pastReleases, setPastReleases] = useState<any>(data?.allTrendingReleases?.past)
    const [weekReleases, setWeekReleases] = useState<any>(data?.allTrendingReleases?.week)
    const [monthReleases, setMonthReleases] = useState<any>(data?.allTrendingReleases?.month)
    const [extendedReleases, setExtendedReleases] = useState<any>(data?.allTrendingReleases?.extended)


    return (
        <Tabs.Group>
            <Tabs.Item title="Recent">
                <div className="grid gap-16 grid-cols-fluid">
                {pastReleases?.map((release: any) => (
                    <ReleasePreview release={release}/>
                ))}
                </div>
            </Tabs.Item>
            <Tabs.Item title="This Week" active>
                <div className="grid gap-16 grid-cols-fluid">
                {weekReleases?.map((release: any) => (
                    <ReleasePreview release={release}/>
                ))}
                </div>
            </Tabs.Item>
            <Tabs.Item title="Next Month">
                <div className="grid gap-16 grid-cols-fluid">
                {monthReleases?.map((release: any) => (
                    <ReleasePreview release={release}/>
                ))}
                </div>
            </Tabs.Item>
            <Tabs.Item title="Too Long...">
                <div className="grid gap-16 grid-cols-fluid">
                {extendedReleases?.map((release: any) => (
                    <ReleasePreview release={release}/>
                ))}
                </div>
            </Tabs.Item>
        </Tabs.Group>
    )
}