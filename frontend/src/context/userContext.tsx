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
}

export const UserContext = createContext<UserContextType>({
  user: null,
  setUser: () => {},
});

const UserContextProvider: FC<PropsWithChildren> = ({ children }:  React.PropsWithChildren) => {
  const [user, setUser] = useState<User | null>(null);

  return (
    <UserContext.Provider value={{ user, setUser }}>
      {children}
    </UserContext.Provider>
  );
};

export default UserContextProvider;