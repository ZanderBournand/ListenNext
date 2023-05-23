import './globals.css'
import NavBar from '../components/navbar'
import FlowbiteContext from '@/context/flowbite'
import { ApolloWrapper } from '@/lib/apollo-provider'
import Bottom from '@/components/footer'
import { SidebarProvider } from '@/context/sidebarContext'

export const metadata = {
  title: 'Create Next App',
  description: 'Generated by create next app',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className="bg-white">
        <FlowbiteContext>
          <SidebarProvider>
            <ApolloWrapper>
              <NavBar/>
              {children}
              <Bottom/>
            </ApolloWrapper>
          </SidebarProvider>
        </FlowbiteContext>;
      </body>
    </html>
  )
}
