import Image from "next/image";
import DefaultCover from "../../public/default_album.png"

export default function ReleasePreview({release}: any) {
    const artistNames = release.artists.map((artist: any) => artist.name).join(", ");
    
    return (
        <div className="flex flex-col items-center justify-start">
        <div className="flex flex-col items-center">
            <div className="mb-2">
                <Image alt="album image" src={release.cover.length === 0 ? DefaultCover : release.cover} width="200" height="200" className="rounded-md"/>
            </div>
            <h2 className="text-md text-c1">{artistNames}</h2>
            <h1 className="text-md font-medium text-center px-4">{release.title}</h1>
        </div>
        </div>
    );
}