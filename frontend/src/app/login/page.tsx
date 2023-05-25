import Image from "next/image";
import Logo from '../../../public/listennext_logo.png'
import { Button, Card, Checkbox, Label, TextInput } from "flowbite-react";
import Login from "@/components/login";

export default async function Home() {      
  return (
    <main>
      <div className="bg-white py-10">
        <div className="max-w-7xl mx-auto px-6 flex flex-col items-center py-12 sm:py-24">
          <div className="flex flex-row items-center">
            <Image
              alt="Flowbite logo"
              src={Logo}
              width="60"
              height="60"
            />
            <a className="self-center whitespace-nowrap pl-3 text-3xl font-semibold text-c1" href="/">
              ListenNext
            </a>
          </div>
          <Login/>
        </div>
      </div>
    </main>
  );  
}