# github-pull-request-dashboard
This is a dashboard to see the open Pull Requests for a given repository, depending on how teams are configured they can display which team a pull request is currently waiting on and the dashboard can be filtered on a user so that they can periodically open it to see what is waiting for them.

## Install
in the terminal
1. run ``npm i`` in the root of the project
2. ``cd`` to the ``backend/`` directory and run ``go build .``
3. in the ``backend/`` directory you have to create a ``.env`` file with the following information:
```env
token=<Your GitHub Personal Access Token>
owner=<Owner of the repository>
repo=<Name of the repository>
```

4. ``cd ..`` back to project root, ``cd`` to ``frontend/``, run ``npm run build``
5. (optional) if you are hosting in a domain path (example.com/dashboard for example) then you have to create a ``.env`` file under frontend with the following field:
```env
VITE_URL_PATH=<if you are hosting on 'example.com/dashboard' then this should be '/dashboard'>
```
6. (optional), if you want to run the frontend on a different port than the default '5173' port then you can change the ``port`` value in ``vite.config.js`` 
7. run by going to the project root end executing ``npm run win`` or ``npm run lin`` depending on if you're on windows or Linux/Mac respectively

