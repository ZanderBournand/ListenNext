import ReleasesTabs from "@/components/releasesTabs";
import TestSidebar from "@/components/sidebar";
import { Button } from "flowbite-react";
import { useState } from "react";

export default async function Home() {  
  return (
    <main>
      <div className="bg-white min-h-screen py-10">
        <div className="max-w-7xl mx-auto px-6">
          <h1 className="text-c4 text-2xl font-semibold pb-10 pl-10">
            New Releases 
            <span className="text-c3 text-xl">- Albums</span>
          </h1>
        </div>
        <div className="flex flex-row max-w-7xl mx-auto px-6">
          <TestSidebar/>
          <ReleasesTabs releaseType="album"/>
        </div>
      </div>
    </main>
  );  
}