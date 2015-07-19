#GitStats

##Usage

./binary token

An authentication token is needed to be able to make 5000 requests/hour to github API.
You can generate one, on your account here : https://github.com/settings/tokens

##Libraries

GitStats uses two external packages:
import "github.com/google/go-github/github" ( a library which allow you to interact with git API with more ease )
import "golang.org/x/oauth2" ( an http authentication library )
