import RecommendationsAccess from "@/components/recommendationsAccess";
import RecommendationsAlert from "@/components/recommendationsAlert";
import ReleasesTabs from "@/components/releasesTabs";
import { LastScrapeTime } from "@/util/queries";
import { getTimeSince } from "@/util/random"

export default async function Home() {    
  const lastScrapeTime = await fetch("http://localhost:8000/query", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    cache: 'no-store',
    body: JSON.stringify({
      query: LastScrapeTime,
    }),
  }).then((res) => res.json())
    .then((res) => {
      if (res?.data?.lastScrapeTime) {
        const utcTimestamp = new Date(res?.data?.lastScrapeTime);
        return getTimeSince(utcTimestamp);
      }
      return "a long time ago...";
  });
      
  return (
    <main>
      <div className="bg-white min-h-screen py-10">
        <div className="max-w-7xl flex flex-col mx-auto px-6 pl-10 pb-8">
          <div className="flex flex-col md:flex-row">
            <h1 className="text-c4 text-2xl font-semibold pb-6">
              New Releases <span className="text-sm bg-red"> - updated {lastScrapeTime}</span>
            </h1>
            <RecommendationsAccess/>
          </div>
          <div className="flex w-full">
            <RecommendationsAlert/>
          </div>
        </div>
        <ReleasesTabs/>
      </div>
    </main>
  );  
}