Here's a typical authentication flow and how each part fits into it:

Register: This endpoint (/register) allows users to create an account by providing their name, username, email, and password. Upon receiving the registration request, the handler hashes the password and stores the user's information in the database.

Login: The login functionality allows users who already have an account to authenticate themselves. Typically, users provide their credentials (email/username and password) to log in. The server then verifies the credentials against the stored information in the database. If the credentials are correct, the server generates a token and returns it to the client. This token can then be used to authenticate subsequent requests.

Token Generation: The token generation endpoint (/token) is usually associated with the login process. After successfully validating the user's credentials, the server generates a JSON Web Token (JWT) and returns it to the client. The client can include this token in the headers of subsequent requests to authenticate itself.

If you decide to implement logout functionality in your application, it typically involves:
Clearing the token from the client-side storage (e.g., localStorage, sessionStorage, cookies).
Optionally, sending a request to the server to invalidate the token or perform any necessary cleanup actions.
Providing feedback to the user confirming that they have been successfully logged out.