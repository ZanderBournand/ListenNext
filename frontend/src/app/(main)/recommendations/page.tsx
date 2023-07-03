'use client'

import CarouselRelease from "@/components/carouselRelease";
import { UserContext } from "@/context/userContext"
import { queryAllRecommendations } from "@/util/queries";
import classNames from "classnames";
import { Button } from "flowbite-react";
import { Disc, Library } from "lucide-react";
import { useContext, useEffect, useState } from "react"

export default function Recommendations() {
    const { user } = useContext(UserContext)
    const [loadingRecommendations, setLoadingRecommendations] = useState(true)
    const [recommendations, setRecommendations] = useState<any>(null)
    const [filteredRecommendations, setFilteredRecommendations] = useState<any>(null)
    const [filterType, setFilterType] = useState<any>("all")

    useEffect(() => {
        const fetchRecommendations = async () => {
            const cachedUser = localStorage.getItem('token');
            if (cachedUser !== null) {
                const { data } = await fetch("http://localhost:8000/query", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${cachedUser}`,
                    },
                    cache: 'no-store',
                    body: JSON.stringify({
                        query: queryAllRecommendations,
                    }),
                }).then((res) => res.json())
                setLoadingRecommendations(false)
                setRecommendations(data?.allRecommendations)
                setFilteredRecommendations(data?.allRecommendations)
            }
        }

        fetchRecommendations()
    }, [])

    const handleAllReleasesButton = () => {
        setFilteredRecommendations(recommendations)
        setFilterType("all")
    }

    const handleAlbumReleasesButton = () => {
        if (recommendations) {
            const filteredAlbums = {
              past: recommendations.past.filter((release: any) => release.type !== 'single'),
              week: recommendations.week.filter((release: any) => release.type !== 'single'),
              month: recommendations.month.filter((release: any) => release.type !== 'single'),
              extended: recommendations.extended.filter((release: any) => release.type !== 'single'),
            };
            setFilteredRecommendations(filteredAlbums);
        } 
        setFilterType("albums")
    }

    const handleSingleReleasesButton = () => {
        if (recommendations) {
            const filteredSingles = {
              past: recommendations.past.filter((release: any) => release.type === 'single'),
              week: recommendations.week.filter((release: any) => release.type === 'sinle'),
              month: recommendations.month.filter((release: any) => release.type === 'single'),
              extended: recommendations.extended.filter((release: any) => release.type === 'single'),
            };
            setFilteredRecommendations(filteredSingles);
        }
        setFilterType("singles")
    }

    return (
        <>
        {user?.is_streaming_auth && !loadingRecommendations ?
            <div className="bg-white min-h-screen py-10">
               <div className="max-w-7xl mx-auto">
                <div className="flex flex-col sm:flex-col md:flex-row">
                        <div className="flex flex-col w-full ml-12 md:w-4/6">
                            <div className="flex flex-col lg:flex-row">
                                <h1 className="text-c4 text-2xl font-semibold pb-6 pl-0">
                                    Your Recommendations
                                </h1>
                                <Button.Group className="ml-0 mb-4 lg:ml-8 lg:mb-0">
                                    <Button 
                                        color="gray"
                                        className={classNames("hover:text-c3",{
                                            "!bg-blue-50 font-semibold !text-blue-600 focus:text-c3": filterType === 'all',
                                        })} 
                                        onClick={handleAllReleasesButton}
                                    >
                                        <p>
                                        All
                                        </p>
                                    </Button>
                                    <Button 
                                        color="gray"  
                                        className={classNames("hover:text-c3",{
                                            "!bg-blue-50 font-semibold !text-blue-600 focus:text-c3": filterType === 'albums',
                                        })}  
                                        onClick={handleAlbumReleasesButton}
                                    >
                                        <Library className="h-4 w-4 mr-1"/>
                                        <p>
                                        Albums
                                        </p>
                                    </Button>
                                    <Button 
                                        color="gray" 
                                        className={classNames("hover:text-c3",{
                                            "!bg-blue-100 font-semibold !text-blue-600 focus:text-c3": filterType === 'singles',
                                        })} 
                                        onClick={handleSingleReleasesButton}
                                    >
                                        <Disc className="h-4 w-4 mr-1"/>
                                        <p>
                                        Singles
                                        </p>
                                    </Button>
                                </Button.Group>
                            </div>
                            <h2 className="text-gray-700 text-xl font-semibold pt-4">Out Now</h2>
                            <CarouselRelease releases={filteredRecommendations?.past}/>
                            <h2 className="text-gray-700 text-xl font-semibold pt-6">This Week</h2>
                            <CarouselRelease releases={filteredRecommendations?.week}/>
                            <h2 className="text-gray-700 text-xl font-semibold pt-6">This Month</h2>
                            <CarouselRelease releases={filteredRecommendations?.month}/>
                            <h2 className="text-gray-700 text-xl font-semibold pt-6">Near Future</h2>
                            <CarouselRelease releases={filteredRecommendations?.extended}/>
                        </div>
                        <div className="flex flex-col w-full ml-12 pt-8 sm:w-full md:w-2/6 md:ml-0 md:pt-0">
                            <h1 className="text-c4 text-2xl font-semibold pb-6 pl-0">
                                Relevant Artists
                            </h1>
                        </div>
                    </div>
               </div>
            </div>
            :
            <>
            {user?.is_streaming_auth && loadingRecommendations ?
                <div>loading</div>
                :
                <></>
            }
            </>
        }
        </>
    )
}