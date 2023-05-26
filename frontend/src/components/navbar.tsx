'use client'

import { Avatar, Button, Dropdown, Navbar, TextInput } from "flowbite-react";
import { useRouter } from 'next/navigation';
import Image from "next/image";
import Logo from '../../public/listennext_logo.png'
import { LogIn, LogOut, Menu, Settings } from "lucide-react";
import { useContext, useState } from "react";
import { UserContext } from "@/context/userContext";
import { Search } from "lucide-react";
import { SearchIcon } from "lucide-react";
import { AiOutlineSearch } from 'react-icons/ai'
import classNames from "classnames";

export default function NavBar() {
  const router = useRouter();
  const [searchInput, setSearchInput] = useState("")  
  const [isNavbarVisible, setIsNavbarVisible] = useState(false);
  const { user, setUser, loadingUser} = useContext(UserContext);

  const toggleNavbar = () => {
    setIsNavbarVisible(!isNavbarVisible);
  };

  const handleKeyDown = (event: any) => {
    if (event.key === "Enter" && searchInput.trim() !== "") {
      const url = `/search/${encodeURIComponent(searchInput)}`;
      router.push(url);
    }
  };

  const handleLogout = () => {
    router.push("/login");
    localStorage.removeItem('token');
    setUser(null)
  }
  
  return (
    <header className="sticky top-0 z-20 border-b border-gray-200 bg-white">
      <div className="max-w-7xl mx-auto px-6">
        <div className="sticky flex flex-wrap md:flex-nowrap items-center justify-between w-full h-max py-3 px-3">
            <div className="flex flex-row w-1/2 items-center">
                <Image
                alt="Flowbite logo"
                src={Logo}
                width="60"
                height="60"
                />
                <a className="self-center whitespace-nowrap pl-3 text-2xl font-semibold text-c1" href="/">
                ListenNext
                </a>
            </div>
            <div className="flex w-1/4 items-center justify-center md:justify-end md:order-1">
                <div className="flex">
                    {loadingUser ?
                    <div role="status" className="max-w-sm animate-pulse pt-1 md:pt-0">
                      <div className="h-8 bg-gray-200 rounded-full dark:bg-gray-700 w-20"></div>
                    </div>
                    :
                    user ?
                    <Dropdown
                    arrowIcon={false}
                    inline={true}
                    label={<Avatar alt="User settings" rounded={true}/>}
                    >
                    <Dropdown.Header>
                        <span className="block text-sm">
                        {user?.display_name}
                        </span>
                        <span className="block truncate text-sm font-medium">
                        {user?.email}
                        </span>
                    </Dropdown.Header>
                    <Dropdown.Item>
                        <Settings className="h-4 w-4"/> <span className="pl-2">Settings</span>
                    </Dropdown.Item>
                    <Dropdown.Item onClick={handleLogout}>
                        <LogOut className="h-4 w-4"/> <span className="pl-2">Sign out</span>
                    </Dropdown.Item>
                    </Dropdown>
                    :
                    <Button size="md" pill className="bg-c6 hover:bg-c1" onClick={() => {router.push("/login")}}>
                    <div className="flex flex-row items-center">
                        Login
                        <LogIn className="ml-2 h-5 w-5" />
                    </div>
                    </Button>
                    }
                    <Button color="light" className="px-0 py-0 w-8 border-0 md:hidden ml-2" data-collapse-toggle="navbar-sticky" aria-controls="navbar-sticky" onClick={toggleNavbar}>
                        <Menu/>
                    </Button>
                </div>
            </div>
            <div className={classNames(
                "w-full items-center justify-center pt-4 md:pt-0 md:flex",
                {
                    "hidden": !isNavbarVisible,
                    "flex": isNavbarVisible
                },
            )}>
                <TextInput
                    id="search1"
                    type="text"
                    icon={(props) => <AiOutlineSearch {...props} />}
                    placeholder="Search Artists..."
                    className="w-full md:w-5/6 lg:w-9/12"
                    value={searchInput}
                    onChange={(event) => setSearchInput(event.target.value)}
                    onKeyDown={handleKeyDown}
                />
            </div>
        </div>
      </div>
    </header>
    )
  }
  