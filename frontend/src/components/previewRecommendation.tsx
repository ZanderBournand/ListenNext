import Image from "next/image";
import DefaultCover from "../../public/default_album.png"
import { Badge } from "flowbite-react";
import { Library, Disc } from "lucide-react";
import { DateTime } from "luxon";
import { ReduceName } from "@/util/titles";

export default function RecommendationPreview({release}: any) {
    const artistNames = release.artists.map((artist: any) => artist.name).join(", ");
    
    const date = DateTime.fromISO(release?.release_date).toUTC();
    const formattedDate = date.toFormat('MMMM d');

    return (
        <div className="flex flex-col items-center justify-start">
        <div className="flex flex-col items-center relative">
            <div className="pb-1 flex flex-row">{formattedDate} 
                {release?.type == 'single' ?
                    <Badge className="mb-1 ml-2" color="blue" size="sm">
                        <div className="flex flex-row items-center">
                            <Disc className="h-4 w-4"/>
                            <span className="pl-1">Single</span>
                        </div>
                    </Badge>
                    :
                    <Badge className="mb-1 ml-2" color="purple" size="sm">
                        <div className="flex flex-row items-center">
                            <Library className="h-4 w-4"/>
                            <span className="pl-1">Album</span>
                        </div>
                    </Badge>
                }
            </div>
            <div className="mb-2">
                <Image alt="album image" src={release?.cover?.length === 0 ? DefaultCover : release.cover} width={200} height={200}  className="rounded-md"/>
            </div>
            <h2 className="text-md text-c1 text-center font-semibold" style={{maxWidth: '275px'}}>{artistNames}</h2>
            <h1 className="text-md font-medium text-center px-4" style={{maxWidth: '275px'}}>{ReduceName(release.title, 50)}</h1>
        </div>
        </div>
    );
}