'use client'

import { Spinner } from "flowbite-react";

export default function Loading() {
    return (
    <div className="bg-white min-h-screen py-10">
            <div className="max-w-7xl mx-auto">
                <div className="flex flex-col sm:flex-col md:flex-row">
                    <div className="flex flex-col w-full ml-12 md:w-3/5">
                        <h1 className="text-c4 text-2xl font-semibold pb-6 pl-0">
                            Search Results
                        </h1>
                    </div>
                    <div className="flex flex-col w-full ml-12 pt-8 sm:w-full md:w-2/5 md:ml-0 md:pt-0">
                        <h1 className="text-c4 text-2xl font-semibold pb-6 pl-0">
                            Related Artists
                        </h1>
                    </div>
                </div>
                <div className="flex w-full items-center justify-center pt-24">
                    <Spinner aria-label="Extra small spinner example" size="lg" />
                </div>
            </div>
        </div>
    );
}