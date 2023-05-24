import Image from "next/image";
import DefaultCover from "../../public/default_album.png"
import { Badge } from "flowbite-react";
import { CalendarCheck } from "lucide-react";

export default function ReleasePreview({release}: any) {
    const artistNames = release.artists.map((artist: any) => artist.name).join(", ");
    
    const date = new Date(release?.release_date);
    const month = date.toLocaleString('en-US', { month: 'long' });
    const day = date.getDate();
    const formattedDate = month + ' ' + day;

    const currentDate = new Date();
    const hasPassed = date <= currentDate
    
    return (
        <div className="flex flex-col items-center justify-start">
        <div className="flex flex-col items-center relative">
            <div className="pb-1 flex flex-row">{formattedDate} {hasPassed && 
                    <Badge className="pb-2 ml-2" color="success" size="sm">
                        <div className="flex flex-row items-center">
                            <CalendarCheck className="h-5 w-5"/>
                            <span className="pl-1">Out Now!</span>
                        </div>
                    </Badge>
                }</div>
            <div className="mb-2">
                <Image alt="album image" src={release?.cover?.length === 0 ? DefaultCover : release.cover} width={200} height={200}  className="rounded-md"/>
            </div>
            <h2 className="text-md text-c1 text-center" style={{maxWidth: '275px'}}>{artistNames}</h2>
            <h1 className="text-md font-medium text-center px-4" style={{maxWidth: '275px'}}>{release.title}</h1>
        </div>
        </div>
    );
}