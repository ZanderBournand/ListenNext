'use client'

import { UserContext } from "@/context/userContext"
import { useContext } from "react"
import { ArrowUpRight } from "lucide-react"
import { Button } from "flowbite-react"
import { useRouter } from "next/navigation"

export default function RecommendationsAccess() {
    const { user } = useContext(UserContext)
    const router = useRouter();

    return (
        <>
        {user?.is_streaming_auth &&
            <Button color="green" className="ml-0 md:ml-8 bg-gray-100 flex flex-row px-0 py-0 items-center text-black relative mx-2 my-1 w-5/6 sm:mt-0 sm:w-max sm:mb-0 hover:bg-green-50 focus:border-0" onClick={() => {router.push('/recommendations')}}>
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    aria-label="Spotify"
                    role="img"
                    viewBox="0 0 512 512"
                    width="18"
                    height="18"
                    className="inline-block align-middle mx-1"
                >
                    <rect width="512" height="512" rx="15%" fill="#3bd75f" />
                    <circle cx="256" cy="256" fill="#fff" r="192" />
                    <g fill="none" stroke="#3bd75f" strokeLinecap="round">
                        <path d="m141 195c75-20 164-15 238 24" strokeWidth="36" />
                        <path d="m152 257c61-17 144-13 203 24" strokeWidth="31" />
                        <path d="m156 315c54-12 116-17 178 20" strokeWidth="24" />
                    </g>
                </svg>
                <div className="text-base pl-1">
                    See Your Recommendations! 
                </div>
                <ArrowUpRight className="h-4 w-4 ml-2"/>
            </Button>
        }
        </>
    )
}

