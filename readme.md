Backend application for a Todo App, written in Go with Echo and Bun ORM.

# Running the app

Run the app with:

`run go .`

# Swagger

The swagger UI is available at `localhost:1323/swagger/index.html`

To call authenticated routes with the Swagger UI you must prefix the token returned from `/signup` or `/login`  with `Bearer` as follows:

`Bearer 73a20efddf336f240075a45ffb7556f8d64d12856bce929ae447e0343d1ee234`