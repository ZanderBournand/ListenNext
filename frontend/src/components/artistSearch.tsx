'use client'

import { Badge, Card } from "flowbite-react";
import Image from "next/image";
import DefaultCover from "../../public/default_album.png"
import { Flame } from "lucide-react";
import ArtistPopularity from "./artistPopularity";
import { ReduceName } from "@/util/titles";

export default function ArtistSearch({artist}: any) {
    return (
        <Card className="w-5/6 h-40 xs:h-40 sm:h-32 md:h-40 lg:h-32 my-4 shadow-sm bg-gray-100/25">
            <div className="flex flex-row">
                <div className="flex h-full items-center">
                    <div className="rounded-3xl overflow-hidden w-24 h-24">
                        <Image
                            alt="artist profile"
                            src={artist?.image ? artist.image : DefaultCover}
                            height={120}
                            width={120}
                            className="object-cover object-center"
                        />
                    </div>
                </div>
                <div className="flex flex-col sm:flex-row md:flex-col lg:flex-row w-full pl-4">
                    <div className="flex flex-col xs:w-full">
                        <div className="text-lg font-bold tracking-tight text-gray-900 flex flex-row items-center">
                        {ReduceName(artist.name, 20)}
                        </div>
                        <div className="flex flex-row items-center pt-1 xs:pt-1 sm:pt-3 md:pt-1 lg:pt-4">
                        <span>Popularity:</span>
                        <ArtistPopularity popularity={artist.popularity} collapse={false}/>
                        </div>
                    </div>
                    <div className="flex flex-col xs:items-start xs:w-4/6 md:w-full lg:w-4/6 xs:pr-4 sm:items-end md:items-start lg:items-end">
                        <div className="text-md font-semibold text-gray-600 pt-2 sm:pt-0 md:pt-2 lg:pt-0">
                        New Releases
                        </div>
                        {artist.recent_releases_count > 0 && artist.upcoming_releases_count === 0 && (
                        <div className="pt-2 sm:pt-4 md:pt-2 lg:pt-4 w-max">
                            <Badge>
                            <span className="font-bold text-sm">{artist.recent_releases_count}</span> Recent
                            </Badge>
                        </div>
                        )}
                        {artist.upcoming_releases_count > 0 && artist.recent_releases_count === 0 && (
                        <div className="pt-2 sm:pt-4 md:pt-2 lg:pt-4 w-max">
                            <Badge>
                            <span className="font-bold text-sm">{artist.upcoming_releases_count}</span> Upcoming
                            </Badge>
                        </div>
                        )}
                        {artist.upcoming_releases_count > 0 && artist.recent_releases_count > 0 && (
                        <>
                            <div className="py-1 w-max">
                            <Badge>
                                <span className="font-bold text-sm">{artist.recent_releases_count}</span> Recent
                            </Badge>
                            </div>
                            <div className="py-1 w-max">
                            <Badge>
                                <span className="font-bold text-sm">{artist.upcoming_releases_count}</span> Upcoming
                            </Badge>
                            </div>
                        </>
                        )}
                        {artist.upcoming_releases_count === 0 && artist.upcoming_releases_count === 0 && (
                        <div className="pt-2 xs:pt-2 sm:pt-4 md:pt-2 lg:pt-4 w-max">
                            <Badge color="gray" size="sm">
                            N/A
                            </Badge>
                        </div>
                        )}
                    </div>
                </div>
            </div>
        </Card>
    )
}