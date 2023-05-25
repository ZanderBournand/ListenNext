'use client'

import { Badge, Card } from "flowbite-react";
import Image from "next/image";
import DefaultCover from "../../public/default_album.png"
import { ReduceName } from "@/util/titles";

export default function RelatedSearch({artist}: any) {
    return (
        <Card className="w-5/6 lg:w-4/6 md:w-5/6 sm:w-3/6 xs:w-3/6 h-24 my-4 shadow-sm bg-gray-100/25">
            <div className="flex flex-row">
                <div className="flex h-full items-center">
                    <div className="rounded-3xl overflow-hidden w-12 h-12">
                        <Image
                            alt="artist profile"
                            src={artist?.image ? artist.image : DefaultCover}
                            height={120}
                            width={120}
                            className="object-cover object-center w-full h-full"
                        />
                    </div>
                </div>
                <div className="flex flex-row w-full pl-4">
                    <div className="flex flex-col w-full">
                        <div className="text-lg font-bold tracking-tight text-gray-900 flex flex-row items-center">
                            {ReduceName(artist.name, 12)}
                        </div>
                        <div className="flex flex-row items-center pt-1">
                            <div className="text-md font-semibold text-gray-600">
                                <span className="whitespace-nowrap">
                                    New
                                    <span className="sm:inline md:hidden lg:inline"> Releases</span>
                                    :
                                </span>
                            </div>
                            <div className="pt-1 pl-2">
                                {artist.recent_releases_count > 0 || artist.upcoming_releases_count > 0 && 
                                    <Badge size="xs" className="rounded-lg">
                                        <span className="font-bold text-sm">
                                            {artist.recent_releases_count + artist.upcoming_releases_count}
                                        </span>
                                        <span className="sm:hidden md:inline lg:hidden"> Release</span>
                                    </Badge>
                                }
                                {artist.upcoming_releases_count == 0 && artist.upcoming_releases_count == 0 &&
                                    <Badge color="gray" size="xs">
                                        N/A
                                    </Badge>
                                }
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </Card>
    )
}

// {artist.recent_releases_count > 0 && artist.upcoming_releases_count == 0 && 
//     <div className="pt-4">
//         <Badge size="xs">
//             <span className="font-bold text-sm">{artist.recent_releases_count}</span> Recent
//         </Badge>
//     </div>
// }
// {artist.upcoming_releases_count > 0 && artist.recent_releases_count == 0 &&
//     <div className="pt-4">
//         <Badge>
//         <span className="font-bold text-sm">{artist.upcoming_releases_count}</span> Upcoming
//         </Badge>
//     </div>
// }
// {artist.upcoming_releases_count > 0 && artist.recent_releases_count > 0 &&
//     <>
//         <div className="py-1">
//             <Badge size="xs">
//                 <span className="font-bold text-sm">{artist.recent_releases_count}</span> Recent
//             </Badge>
//         </div>
//         <div className="py-1">
//             <Badge>
//             <span className="font-bold text-sm">{artist.upcoming_releases_count}</span> Upcoming
//             </Badge>
//         </div>
//     </>
// }
// {artist.upcoming_releases_count == 0 && artist.upcoming_releases_count == 0 &&
//     <div className="pt-4">
//         <Badge color="gray" size="sm">
//         N/A
//         </Badge>
//     </div>
// }