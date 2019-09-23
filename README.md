# Firebase Remote Config with Golang
A simple example og using the Firebase Remote Config REST API.

# Setup
1. Create service account from your Firebase Console -> Service Account and download the JSON file.
2. Copy the JSON credential to root directory of the project and rename it as `service-account.json`
3. Change `PROJECT_ID` variable in **app.go** to your firebase project id.

# Run

`go run app.go` 

- Get Active Template (operation number 1)
  - result of the template will be store in file named `config.json`
  - Etag also stored in **etag.txt** it is for publishing new template.
  
- Publish Template (operation number 2)
  - If you want to change or add a new template/key-value, edit the `config.json` file
  - Publish new template using operation number 2.
  
- Rollback (operation number 3)
  - You can back to the previous version of the template by choosing operation number 3 and input version number you want.
  - `config.json` template will be change based on version you choose.
  
- Show version list (operation number 4) 
  - You can print list of version you have in firebase console, the printed version size depends on your input size.
 
# Etags

Each time the Remote Config template it retrieved an ETag is included. This ETag is a unique identifier of the current template on the server. When submitting updates to the template it contains the latest ETag from **etag.txt** to ensure that your updates are consistent.
