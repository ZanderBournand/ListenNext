import ArtistSearch from "@/components/artistSearch";
import RelatedSearch from "@/components/relatedSearch";
import { querySearchArtists } from "@/util/queries";

export default async function Search({params}: any) {      
    const { data } = await fetch("http://localhost:8000/query", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        cache: 'no-store',
        body: JSON.stringify({
            query: querySearchArtists,
              variables: {
                query: params.search ? decodeURIComponent(params.search) : ""
              }
        }),
    }).then((res) => res.json());
    
    return (
        <main>
        <div className="bg-white min-h-screen py-10">
            <div className="max-w-7xl mx-auto">
                <div className="flex flex-col sm:flex-col md:flex-row">
                    <div className="flex flex-col w-full ml-12 md:w-3/5">
                        <h1 className="text-c4 text-2xl font-semibold pb-6 pl-0">
                            Search Results <span className="text-c6 text-xl">({data?.searchArtists?.results.length})</span>
                        </h1>
                        {data?.searchArtists?.results.map((artist: any) => (
                            <ArtistSearch artist={artist}/>
                        ))}
                    </div>
                    <div className="flex flex-col w-full ml-12 pt-8 sm:w-full md:w-2/5 md:ml-0 md:pt-0">
                        <h1 className="text-c4 text-2xl font-semibold pb-6 pl-0">
                            Related Artists <span className="text-c6 text-xl">({data?.searchArtists?.related_artists.length})</span>
                        </h1>
                        {data?.searchArtists?.related_artists.map((artist: any) => (
                            <RelatedSearch artist={artist}/>
                        ))}
                    </div>
                </div>
            </div>
        </div>
        </main>
    );  
}