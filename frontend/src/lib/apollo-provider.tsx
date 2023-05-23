'use client'

import {
    ApolloClient,
    ApolloLink,
    HttpLink,
    SuspenseCache,
  } from "@apollo/client";
  import {
    ApolloNextAppProvider,
    NextSSRInMemoryCache,
    SSRMultipartLink,
  } from "@apollo/experimental-nextjs-app-support/ssr";
  
  function makeClient() {
    const httpLink = new HttpLink({
        uri: "http://localhost:8000/query",
    });
  
    return new ApolloClient({
      connectToDevTools: true,
      cache: new NextSSRInMemoryCache(),
      link:
        typeof window === "undefined"
          ? ApolloLink.from([
              new SSRMultipartLink({
                stripDefer: true,
              }),
              httpLink,
            ])
          : httpLink,
    });
  }
  
  function makeSuspenseCache() {
    return new SuspenseCache();
  }
  
  export function ApolloWrapper({ children }: React.PropsWithChildren) {
    return (
      <ApolloNextAppProvider
        makeClient={makeClient}
        makeSuspenseCache={makeSuspenseCache}
      >
        {children}
      </ApolloNextAppProvider>
    );
  }