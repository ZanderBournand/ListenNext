'use client'

import { UserContext } from "@/context/userContext";
import { refreshLogin } from "@/util/mutations";
import { useContext, useEffect } from "react"

export default function LoginRefresher() {
    const { setUser, setLoadingUser } = useContext(UserContext)

    useEffect(() => {
        const fetchData = async () => {
            const cachedUser = localStorage.getItem('token');
            if (cachedUser !== null) {
                const { data }  = await fetch("http://localhost:8000/query", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${cachedUser}`,
                    },
                    cache: 'no-cache',
                    body: JSON.stringify({
                        query: refreshLogin,
                    }),
                }).then((res) => res.json());

                if (data !== null) {
                    console.log(data?.auth?.refreshLogin.is_streaming_auth)
                    setUser(data?.auth?.refreshLogin)
                }
            }
            setLoadingUser(false)
        };

        fetchData();
    }, []);

    return (<></>)
}