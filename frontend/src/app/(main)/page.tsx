import RecommendationsAlert from "@/components/recommendationsAlert";
import ReleasesTabs from "@/components/releasesTabs";

export default async function Home() {      
  return (
    <main>
      <div className="bg-white min-h-screen py-10">
        <div className="max-w-7xl flex flex-col mx-auto px-6 pl-10 pb-8">
          <h1 className="w-full sm:w-2/6 text-c4 text-2xl font-semibold pb-6">
            New Releases 
          </h1>
          <div className="flex w-full">
            <RecommendationsAlert/>
          </div>
        </div>
        <ReleasesTabs/>
      </div>
    </main>
  );  
}