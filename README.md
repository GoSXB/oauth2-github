## Simple Github OAuth2 authentication with Golang

### 1. Create a github oauth application

Go in your Profile -> Oauth Applications -> Register a new application

* Application name: whatever you want
* Homepage URL: whatever you want
* Application Description: whatever you want
* Authorization Callback URL: `http://localhost:5000`

Then click on 'Register application'

Then you'll get an Application ID and Application Secret

### 2. Run the app with the credentials

```
go run main.go <app id> <app secret>
```

* Follow the link displayed in your terminal
* Login in your github account
* Get redirected on the app
* You can now display the authenticated user profile
