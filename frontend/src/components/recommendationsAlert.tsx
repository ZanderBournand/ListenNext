'use client'

import { UserContext } from "@/context/userContext";
import { SpotifyLoginUrl } from "@/util/queries";
import { Badge, Button } from "flowbite-react"
import { Info, Trash2, X } from "lucide-react"
import { useContext, useState } from "react"


export default function RecommendationsAlert() {
    const [dismissed, setDismissed] = useState(false)
    const { user, loadingUser } = useContext(UserContext)
    
    const onDismiss = () => {
        setDismissed(true)
    }

    const handleSpotifyLogin = async () => {
        const { data } = await fetch("http://localhost:8000/query", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            cache: 'no-store',
            body: JSON.stringify({
                query: SpotifyLoginUrl,
            }),
        }).then((res) => res.json());
  
        if (data !== null) {
          window.location.href = data?.spotifyUrl
        }
    }

    return (
        <>
        {!loadingUser && !user?.is_streaming_auth && 
            <div className="w-full pr-2 sm:pr-0 md:w-5/6" hidden={dismissed}>
                <div className="flex flex-col rounded-md bg-blue-200/25 py-2 px-4">
                    <div className="flex flex-col w-full">
                        <div className="flex flex-row items-center justify-between">
                            <div className="flex flex-row items-center">
                                <Info className="h-6 w-6" color="#355c7d"/>
                                <span className="pl-2 font-medium text-lg text-c3"><span className="hidden sm:inline">New</span> Recommendations <span className="hidden sm:inline">Feature</span> <span className="inline sm:hidden">Out Now</span>!</span>
                            </div>
                            <button className="rounded-md bg-blue-200/0 hover:bg-blue-200/25" onClick={onDismiss}>
                                <X className="h-6 w-6" color="#4476a1"/>
                            </button>
                        </div>
                        <div className="flex flex-row pt-3 pr-4 pl-1">
                            <span className="text-md font-light text-gray-700">
                                Log in with your streaming platform (
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        aria-label="Spotify"
                                        role="img"
                                        viewBox="0 0 512 512"
                                        width="23"
                                        height="23"
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
                                Spotify or 
                                <svg fill="#000000" 
                                    width="23px" 
                                    height="23px" 
                                    viewBox="0 0 24 24" 
                                    xmlns="http://www.w3.org/2000/svg"
                                    data-name="Layer 1"
                                    className="inline-block align-middle mx-1 mb-1"
                                >
                                    <path d="M14.94,5.19A4.38,4.38,0,0,0,16,2,4.44,4.44,0,0,0,13,3.52,4.17,4.17,0,0,0,12,6.61,3.69,3.69,0,0,0,14.94,5.19Zm2.52,7.44a4.51,4.51,0,0,1,2.16-3.81,4.66,4.66,0,0,0-3.66-2c-1.56-.16-3,.91-3.83.91s-2-.89-3.3-.87A4.92,4.92,0,0,0,4.69,9.39C2.93,12.45,4.24,17,6,19.47,6.8,20.68,7.8,22.05,9.12,22s1.75-.82,3.28-.82,2,.82,3.3.79,2.22-1.24,3.06-2.45a11,11,0,0,0,1.38-2.85A4.41,4.41,0,0,1,17.46,12.63Z"/>
                                </svg>
                                Apple Music) and get personalized recommendations for upcoming releases based on your listening history.
                            </span>
                        </div>
                        <div className="flex flex-col sm:flex-row items-center pt-6">
                            <Button color="light" className="relative w-5/6 sm:mt-0 mx-1 my-1 sm:w-max sm:mb-0" onClick={handleSpotifyLogin}>
                                <svg
                                xmlns="http://www.w3.org/2000/svg"
                                aria-label="Spotify"
                                role="img"
                                viewBox="0 0 512 512"
                                width="20"
                                height="20"
                                >
                                    <rect width="512" height="512" rx="15%" fill="#3bd75f" />
                                    <circle cx="256" cy="256" fill="#fff" r="192" />
                                    <g fill="none" stroke="#3bd75f" strokeLinecap="round">
                                        <path d="m141 195c75-20 164-15 238 24" strokeWidth="36" />
                                        <path d="m152 257c61-17 144-13 203 24" strokeWidth="31" />
                                        <path d="m156 315c54-12 116-17 178 20" strokeWidth="24" />
                                    </g>
                                </svg>
                                <span className="pl-2">Log in w/ Spotify</span>
                            </Button>
                            <Button color="light" className="relative w-5/6 sm:mt-0 mx-1 my-1 sm:w-max sm:mb-0" disabled>
                                <svg fill="#000000" 
                                    width="20px" 
                                    height="20px" 
                                    viewBox="0 0 24 24" 
                                    xmlns="http://www.w3.org/2000/svg"
                                    data-name="Layer 1"
                                >
                                    <path d="M14.94,5.19A4.38,4.38,0,0,0,16,2,4.44,4.44,0,0,0,13,3.52,4.17,4.17,0,0,0,12,6.61,3.69,3.69,0,0,0,14.94,5.19Zm2.52,7.44a4.51,4.51,0,0,1,2.16-3.81,4.66,4.66,0,0,0-3.66-2c-1.56-.16-3,.91-3.83.91s-2-.89-3.3-.87A4.92,4.92,0,0,0,4.69,9.39C2.93,12.45,4.24,17,6,19.47,6.8,20.68,7.8,22.05,9.12,22s1.75-.82,3.28-.82,2,.82,3.3.79,2.22-1.24,3.06-2.45a11,11,0,0,0,1.38-2.85A4.41,4.41,0,0,1,17.46,12.63Z"/>
                                </svg>
                                <span className="pl-2">Log in w/ Apple</span>
                                <Badge  className="bg-blue-200/75 absolute -top-3 -right-1">
                                    Coming Soon!
                                </Badge>
                            </Button>
                            <Button color="gray" className="bg-white text-black relative mx-2 my-1 w-5/6 sm:mt-0 sm:w-max sm:mb-0" onClick={onDismiss}>
                                <Trash2 className="h-4 w-4"/>
                                <span className="pl-1">Close</span>
                            </Button>
                        </div>
                    </div>
                </div>
            </div>
        }
        </>
    )
}