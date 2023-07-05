'use client'

import { Button, Card } from "flowbite-react";
import Image from "next/image";
import DefaultCover from "../../public/default_album.png"
import { ReduceName } from "@/util/titles";
import { ArrowRightFromLine } from "lucide-react";

export default function ArtistRecommendation({artist}: any) {    
    return (
        <Card className="w-5/6 lg:w-4/5 md:w-4/5 sm:w-3/6 xs:w-4/6 h-24 md:h-36 lg:h-24 my-4 shadow-sm bg-gray-100/25">
            <div className="flex flex-row md:flex-col lg:flex-row items-center">
                <div className="flex h-full items-center">
                    <div className="rounded-3xl overflow-hidden w-20 h-20 lg:w-20 lg:h-20">
                        <Image
                            alt="artist profile"
                            src={artist?.image ? artist.image : DefaultCover}
                            height={120}
                            width={120}
                            className="object-cover object-center"
                        />
                    </div>
                </div>
                <div className="flex flex-row w-full pl-4">
                    <div className="flex flex-col w-full pt-0 md:pt-2 pl-2 md:pl-0">
                        <div className="text-lg font-bold tracking-tight justify-start md:justify-center lg:justify-start text-gray-900 flex flex-row items-center">
                            <span className="text-center">{ReduceName(artist.name, 12)}</span>
                        </div>
                        <Button color="gray" size="xs" className="bg-white text-black relative w-5/6  mt-2 md:w-max sm:mb-0 block md:hidden lg:block">
                                <span className="pl-1">See More</span>
                                <ArrowRightFromLine className="h-4 w-4 ml-2"/>
                        </Button>
                    </div>
                </div>
            </div>
        </Card>
    )
}