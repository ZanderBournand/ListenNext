'use client'

import { Button, Navbar } from "flowbite-react";
import { useRouter } from 'next/navigation';
import Image from "next/image";
import Logo from '../../public/listennext_logo.png'
import { LogIn } from "lucide-react";
import { useState } from "react";

export default function NavBar() {
  const router = useRouter();
  const [searchInput, setSearchInput] = useState("")  

  const handleKeyDown = (event: any) => {
    if (event.key === "Enter" && searchInput.trim() !== "") {
      const url = `/search/${encodeURIComponent(searchInput)}`;
      router.push(url);
    }
  };
  
  return (
    <header className="sticky top-0 z-20 bg-white">
      <div className="max-w-7xl mx-auto px-6">
        <Navbar fluid>
          <Navbar.Brand href="/">
            <Image
              alt="Flowbite logo"
              src={Logo}
              width="60"
              height="60"
            />
            <span className="self-center whitespace-nowrap pl-3 text-2xl font-semibold text-c1">
              ListenNext
            </span>
          </Navbar.Brand>
          <div className="flex md:order-1">
            <Button size="md" pill className="bg-c6 hover:bg-c1">
              <div className="flex flex-row items-center">
                Login
                <LogIn className="ml-2 h-5 w-5" />
              </div>
            </Button>
            <Navbar.Toggle className="hover:bg-c2 focus:bg-c2 focus:border-0"/>
          </div>
          <Navbar.Collapse className="pt-2">
            <div className="relative md:block">
              <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                <svg className="w-5 h-5 text-gray-500" aria-hidden="true" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd"></path></svg>
                <span className="sr-only">Search icon</span>
              </div>
              <input 
                type="text" 
                id="search-navbar" 
                className="block p-2 pl-10 text-sm text-gray-900 border border-gray-300 rounded-lg bg-gray-50 w-full sm:w-full md:w-72 lg:w-96 focus:border-c6 focus:ring-c2" 
                placeholder="Search Artists..."
                value={searchInput}
                onChange={(event) => setSearchInput(event.target.value)}
                onKeyDown={handleKeyDown}
              />
            </div>
          </Navbar.Collapse>
        </Navbar>
      </div>
    </header>
    )
  }
  