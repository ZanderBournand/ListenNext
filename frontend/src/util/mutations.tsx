export const registerUser = `
  mutation SignUp($display_name: String!, $email: String!, $password: String!) {
    auth {
      register(input: {
        display_name: $display_name
        email: $email
        password: $password
      })
    }
  }
`

export const loginUser = `
  mutation SignIn($email: String!, $password: String!) {
    auth {
      login(email: $email, password: $password)
    }
  }
`

export const refreshLogin = `
  mutation RefreshLogin {
    auth{
      refreshLogin
    }
  }
`

export const spotifyLogin = `
  mutation SpotifyLogin($code: String!){
    auth{
      spotifyLogin(code: $code)
    }
  }
`