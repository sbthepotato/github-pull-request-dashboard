# github-pull-request-dashboard

This is a dashboard to see the open Pull Requests for a given repository, depending on how teams are configured they can display which team a pull request is currently waiting on and the dashboard can be filtered on a user so that they can periodically open it to see what is waiting for them.

## Configuration

Configuration is done through a `.env` file in the `backend/` directory. It is important to set this before building as the build steps rely on being done in order with that information filled out.

## Install

in the terminal

1. in the `backend/` directory you have to create a `.env` file with the following information:

```env
token=<Your GitHub Personal Access Token>
owner=<Owner of the repository>
repo=<Name of the repository>
port=<Port the server will run on>
base_path=<any subdomain being used (/dashboard for example)>
```

2. `cd` to the `frontend/` directory and run `npm run build`
3. `cd` back to the `backend/` directory and run `go build .`
4. the project can now be run from the executable in the backend directory
