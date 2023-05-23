'use client'

import ReleasePreview from "./preview";
import { useEffect, useState } from "react";
import { queryTrendingReleases } from '../util/queries'
import { Button, Spinner } from "flowbite-react";
import { useLazyQuery, useQuery } from "@apollo/client";

export default function ReleasesGrid({releaseType, period}: any) {    
    const {data: initialData, loading} = useQuery<any>(queryTrendingReleases, {
        variables: {
            releaseType: releaseType,
            direction: 'next',
            reference: 0,
            period: period
        },
    })
    const [releases, setReleases] = useState<any>(initialData?.trendingReleases?.releases)
    const [showMore, setShowMore] = useState<boolean>(initialData?.trendingReleases?.next)
    const [getMoreReleases, {data: moreData}] = useLazyQuery<any>(queryTrendingReleases);

    useEffect(() => {
        if (initialData) {
            setReleases([...initialData?.trendingReleases?.releases])
            setShowMore(initialData?.trendingReleases?.next)
        }
    }, [initialData])
    
    useEffect(() => {
        if (moreData) {
            setReleases([...releases, ...moreData?.trendingReleases?.releases])
            setShowMore(moreData?.trendingReleases?.next)
        }
    }, [moreData])

    const handleShowMore = () => {
        getMoreReleases({
            variables: {
                releaseType: releaseType,
                direction: 'next',
                reference: releases.length,
                period: period
            }
        })
    }
    
    return (
        <>
        {loading ?
            <div className="flex items-center justify-center pt-24">
                <Spinner
                    aria-label="Extra small spinner example"
                    size="lg"
                />
            </div>
            :
            <>
            <div className="grid gap-x-4 gap-y-12 grid-cols-fluid">
                {releases?.map((release: any) => (
                    <ReleasePreview key={release._id} release={release}/>
                ))}
            </div>
            {showMore && 
            <div className="pt-16 flex justify-center items-center">
                <Button
                    color="light"
                    pill={true}
                    onClick={handleShowMore}
                >
                    Show More
                </Button>
            </div>
            }
            </>
        }   
        </>
    )
}