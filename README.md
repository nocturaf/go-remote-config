# go-remote-config
Simple Firebase Remote Config REST API with Golang.

# Setup

1. Please generate your own private credentials service account from Firebase Console, place it on root directory
2. Change PROJECT_ID to yours in app.go
3. Run GetRemoteConfig at first to get latest etag before start publishing (add/update key). 
