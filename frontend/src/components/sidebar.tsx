'use client'

import { Sidebar } from "flowbite-react";

export default function TestSidebar() {
    return (
        <div className="flex flex-col w-auto bg-gray-100/50 rounded-xl h-96 sticky top-24">
            <h3 className="pl-8 py-4 text-xl font-semibold">Filter</h3>
            <div className="border-b border-slate-400 mx-6"></div>
            <Sidebar aria-label="Default sidebar example" className="pt-6">
                <Sidebar.Items>
                <Sidebar.ItemGroup>
                    <Sidebar.Item href="#" label="171" abelColor="alternative">
                    All
                    </Sidebar.Item>
                    <Sidebar.Item href="#" label="67" abelColor="alternative">
                    Albums
                    </Sidebar.Item>
                    <Sidebar.Item href="#" label="104" abelColor="alternative">
                    Singles
                    </Sidebar.Item>
                </Sidebar.ItemGroup>
                </Sidebar.Items>
            </Sidebar>
        </div>
    )
}