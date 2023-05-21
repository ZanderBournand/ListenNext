'use client'

import { Button, Navbar } from "flowbite-react";
import Image from "next/image";
import Logo from '../../public/listennext_logo.png'

export default function NavBar() {
    return (
    <header className="sticky top-0 z-20 bg-white">
      <Navbar fluid>
          <Navbar.Brand href="/">
            <Image
              alt="Flowbite logo"
              src={Logo}
              width="70"
              height="70"
            />
            <span className="self-center whitespace-nowrap px-3 text-2xl font-semibold text-c1">
              ListenNext
            </span>
          </Navbar.Brand>
          <div className="flex md:order-1">
            <Button size="lg" pill className="bg-c3 hover:bg-c1">
              Login
            </Button>
            <Navbar.Toggle />
          </div>
          <Navbar.Collapse className="mr-24 pt-2">
            <Navbar.Link href="/" active>Home</Navbar.Link>
            <Navbar.Link href="/">Recommendations</Navbar.Link>
          </Navbar.Collapse>
        </Navbar>
    </header>
    )
  }
  