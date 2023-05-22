'use client'

import { Footer } from "flowbite-react";
import Logo from '../../public/listennext_logo.png'
import Image from "next/image";

export default function Bottom() {
    return (
        <div className="max-w-7xl mx-auto px-6">
            <Footer container={true} className="flex w-full">
                <div className="w-full text-center">
                    <div className="w-full justify-between sm:flex sm:items-center sm:justify-between">
                    <div className="flex flex-row">
                        <Image
                            alt="Flowbite logo"
                            src={Logo}
                            width="50"
                            height="50"
                        />
                        <a href="/" className="self-center whitespace-nowrap px-3 text-2xl font-semibold text-c1">
                        ListenNext
                        </a>
                    </div>
                    <Footer.LinkGroup>
                        <Footer.Link href="#">
                        About
                        </Footer.Link>
                        <Footer.Link href="#">
                        Privacy Policy
                        </Footer.Link>
                        <Footer.Link href="#">
                        Contact
                        </Footer.Link>
                    </Footer.LinkGroup>
                    </div>
                    <Footer.Divider />
                    <Footer.Copyright
                    href="#"
                    by="ListenNext™"
                    year={2023}
                    />
                </div>
            </Footer>
        </div>
    )
}