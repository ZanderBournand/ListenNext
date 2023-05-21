import ReleasesGrid from "@/components/relasesGrid";

export default async function Home() {
  return (
    <main>
      <div className="bg-white min-h-screen py-10">
        <h1 className="text-c4 text-2xl font-semibold pb-10 pl-6">
          New Releases 
          <span className="text-c3 text-xl">- Albums</span>
        </h1>
        <ReleasesGrid/>
      </div>
    </main>
  );  
}