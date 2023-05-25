'use client'

import { UserContext } from "@/context/userContext";
import { loginUser, registerUser } from "@/util/mutations";
import { Badge, Button, Card, Checkbox, Label, TextInput } from "flowbite-react";
import { useRouter } from 'next/navigation';
import React, { useContext, useState } from "react";

export default function Login() {
    const router = useRouter();  
  
    const [login, setLogin] = useState<any>(true)
    const [displayName, setDisplayName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");

    const [emailTaken, setEmailTaken] = useState(false)
    const [wrongLoginInfo, setWrongLoginInfo] = useState(false)
    const [loading, setLoading] = useState(false)

    const { setUser } = useContext(UserContext);
    
    const handleFormSubmit = async (event: any) => {
        event.preventDefault();
        setLoading(true)

        if (login) {
            const { data } = await fetch("http://localhost:8000/query", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                cache: 'no-store',
                body: JSON.stringify({
                    query: loginUser,
                    variables: {
                        email: email,
                        password: password,
                    }
                }),
            }).then((res) => res.json());
            setLoading(false)
            
            if (data == null) {
                setWrongLoginInfo(true)
                return
            }

            const token = data?.auth?.login?.token;
            localStorage.setItem('token', token);

            setUser(data?.auth?.login?.user)
            router.push("/");
        } else {
            const { data } = await fetch("http://localhost:8000/query", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                cache: 'no-store',
                body: JSON.stringify({
                    query: registerUser,
                    variables: {
                        display_name: displayName,
                        email: email,
                        password: password,
                    }
                }),
            }).then((res) => res.json());
            setLoading(false)
            
            if (data == null) {
                setEmailTaken(true)
                return
            }

            const token = data?.auth?.register?.token;
            localStorage.setItem('token', token);

            setUser(data?.auth?.login?.user)
            router.push("/");
        }
      };
    
    return (
        <div className="flex w-full flex-col items-center pt-6">
            <Card className="w-full sm:w-max">
            <h1 className="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white">
                {login ? "Sign in to your account" : "Create a new account"}
              </h1>
            <div className="flex justify-evenly items-center pt-0 flex-col sm:flex-row sm:pt-1">
                <Button color="light" className="w-full py-0 mx-2 mt-2 mb-4 sm:mt-0 sm:w-max sm:mb-0">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        aria-label="Spotify"
                        role="img"
                        viewBox="0 0 512 512"
                        width="26"
                        height="26"
                    >
                        <rect width="512" height="512" rx="15%" fill="#3bd75f" />
                        <circle cx="256" cy="256" fill="#fff" r="192" />
                        <g fill="none" stroke="#3bd75f" strokeLinecap="round">
                            <path d="m141 195c75-20 164-15 238 24" strokeWidth="36" />
                            <path d="m152 257c61-17 144-13 203 24" strokeWidth="31" />
                            <path d="m156 315c54-12 116-17 178 20" strokeWidth="24" />
                        </g>
                    </svg>
                    <span className="pl-2">{login ? "Log in" : "Sign up"} with Spotify</span>
                </Button>
                <Button color="light" className="relative w-full py-0 mx-2 mt-2 mb-2 sm:mt-0 sm:w-max sm:mb-0">
                    <svg fill="#000000" 
                        width="26px" 
                        height="26px" 
                        viewBox="0 0 24 24" 
                        xmlns="http://www.w3.org/2000/svg"
                        data-name="Layer 1"
                    >
                        <path d="M14.94,5.19A4.38,4.38,0,0,0,16,2,4.44,4.44,0,0,0,13,3.52,4.17,4.17,0,0,0,12,6.61,3.69,3.69,0,0,0,14.94,5.19Zm2.52,7.44a4.51,4.51,0,0,1,2.16-3.81,4.66,4.66,0,0,0-3.66-2c-1.56-.16-3,.91-3.83.91s-2-.89-3.3-.87A4.92,4.92,0,0,0,4.69,9.39C2.93,12.45,4.24,17,6,19.47,6.8,20.68,7.8,22.05,9.12,22s1.75-.82,3.28-.82,2,.82,3.3.79,2.22-1.24,3.06-2.45a11,11,0,0,0,1.38-2.85A4.41,4.41,0,0,1,17.46,12.63Z"/>
                    </svg>
                    <span className="pl-2">{login ? "Log in" : "Sign up"} with Apple</span>
                    <Badge color="info" className="absolute -top-3 -right-1">
                        Coming Soon!
                    </Badge>
                </Button>
            </div>
            <div className="flex flex-row items-center">
                <div className="border-t-2 border-gray-400/50 rounded-lg w-full"></div>
                <span className="text-md text-gray-600 px-2">or</span>
                <div className="border-t-2 border-gray-400/50 rounded-lg w-full"></div>
            </div>
            <form className="flex flex-col gap-4" onSubmit={handleFormSubmit}>
              {!login &&
                <div>
                    <div className="mb-2 block">
                        <Label
                        htmlFor="name1"
                        value="Display name"
                        />
                    </div>
                    <TextInput
                        id="name1"
                        type="text"
                        placeholder="MusicEnjoyer3000"
                        required={true}
                        value={displayName}
                        onChange={(e) => setDisplayName(e.target.value)}
                    />
                </div>
              }
              <div>
                <div className="mb-2 block">
                  <Label
                    htmlFor="email1"
                    value="Your email"
                  />
                </div>
                <TextInput
                  id="email1"
                  type="email"
                  placeholder="name@gmail.com"
                  required={true}
                  value={email}
                  color={emailTaken || wrongLoginInfo ? "failure": "gray"}
                  onChange={(e) => {
                    setEmail(e.target.value)
                    setEmailTaken(false)
                    setWrongLoginInfo(false)
                  }}
                  helperText={emailTaken ? <React.Fragment><span className="font-medium text-red-500 pl-1">Oops!</span>{' '}<span className="text-red-500">Username already taken!</span></React.Fragment>: wrongLoginInfo ? <React.Fragment><span className="font-medium text-red-500 pl-1">Oops!</span>{' '}<span className="text-red-500">Email and/or password are incorrect!</span></React.Fragment> : ""}
                />
              </div>
              <div>
                <div className="mb-2 block">
                  <Label
                    htmlFor="password1"
                    value="Your password"
                  />
                </div>
                <TextInput
                  id="password1"
                  type="password"
                  placeholder="••••••••"
                  required={true}
                  value={password}
                  color={wrongLoginInfo ? "failure": "gray"}
                  onChange={(e) => {
                    setPassword(e.target.value)
                    setWrongLoginInfo(false)
                  }}
                  helperText={wrongLoginInfo ? <React.Fragment><span className="font-medium text-red-500 pl-1">Oops!</span>{' '}<span className="text-red-500">Email and/or password are incorrect!</span></React.Fragment>: ""}
                />
              </div>
              {!login &&
                <div>
                    <div className="mb-2 block">
                        <Label
                        htmlFor="password2"
                        value="Confirm Password"
                        />
                    </div>
                    <TextInput
                        id="password2"
                        type="password"
                        placeholder="••••••••"
                        required={true}
                        value={confirmPassword}
                        onChange={(e) => setConfirmPassword(e.target.value)}
                    />
                </div>
              }
              <div className="flex flex-row items-center justify-between">
                <div className="flex items-center gap-2">
                    <Checkbox id="remember" />
                    <Label htmlFor="remember">
                        <span className="text-gray-500 font-normal">Remember me</span>
                    </Label>
                </div>
                {login && 
                    <div className="pr-2">
                        <a className="text-sm font-medium text-c3 hover:underline dark:text-primary-500">Forgot password?</a>
                    </div>
                }
              </div>
              <Button type="submit" className="mt-2 bg-c1 hover:bg-c4" isProcessing={loading}>
                {login ? "Sign In": "Sign Up"}
              </Button>
              <div className="pt-2">
                {login ?
                <p className="text-sm font-light text-gray-500 dark:text-gray-400">
                    Don’t have an account yet? 
                    <a 
                    className="font-medium pl-1 text-c3 hover:underline dark:text-primary-500" 
                    onClick={() => setLogin(!login)}
                    >
                    Sign up
                    </a>
                </p>
                :
                <p className="text-sm font-light text-gray-500 dark:text-gray-400">
                    Already have an account? 
                    <a 
                    className="font-medium pl-1 text-c3 hover:underline dark:text-primary-500" 
                    onClick={() => setLogin(!login)}
                    >
                    Sign in
                    </a>
                </p>
                }
              </div>
            </form>
            </Card>
          </div>
    )
}