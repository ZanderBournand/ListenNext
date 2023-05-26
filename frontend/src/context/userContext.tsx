'use client'

import React, { FC, PropsWithChildren, createContext, useState } from 'react';

interface User {
    _id: number;
    display_name: string;
    email: string;
  }

interface UserContextType {
  user: User | null;
  setUser: (user: User | null) => void;
  loadingUser: Boolean;
  setLoadingUser: (loadingUser: Boolean) => void;
}

export const UserContext = createContext<UserContextType>({
  user: null,
  setUser: () => {},
  loadingUser: true,
  setLoadingUser: () => {},
});

const UserContextProvider: FC<PropsWithChildren> = ({ children }:  React.PropsWithChildren) => {
  const [user, setUser] = useState<User | null>(null);
  const [loadingUser, setLoadingUser] = useState<Boolean>(true);

  return (
    <UserContext.Provider value={{ user, setUser, loadingUser, setLoadingUser }}>
      {children}
    </UserContext.Provider>
  );
};

export default UserContextProvider;