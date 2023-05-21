'use client'

import ReleasePreview from "./preview";
import { useSuspenseQuery } from "@apollo/experimental-nextjs-app-support/ssr";
import { useEffect, useState } from "react";
import { query } from '../util/queries'
import { Button } from "flowbite-react";
import { useLazyQuery } from "@apollo/client";

export default function ReleasesGrid({releaseType, period}: any) {
    const {data: initialData} = useSuspenseQuery<any>(query, {variables: {
        releaseType: releaseType,
        direction: 'next',
        reference: 0,
        period: period
    }})
    const [releases, setReleases] = useState<any>(initialData?.trendingReleases?.releases)
    const [getMoreReleases, {data: moreData}] = useLazyQuery<any>(query);
    
    useEffect(() => {
        if (moreData) {
            setReleases([...releases, ...moreData?.trendingReleases?.releases])
        }
    }, [moreData])

    const handleShowMore = () => {
        getMoreReleases({variables: {
            releaseType: releaseType,
            direction: 'next',
            reference: releases.length,
            period: period
        }})
    }
    
    return (
        <>
        <div className="grid gap-16 grid-cols-fluid">
            {releases?.map((release: any) => (
                <ReleasePreview key={release._id} release={release}/>
            ))}
        </div>
        <div className="pt-16 flex justify-center items-center">
            <Button
                color="light"
                pill={true}
                onClick={handleShowMore}
            >
                Show More
            </Button>
        </div>
        </>
    )
}