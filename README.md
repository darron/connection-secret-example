## Connection Secret Example

Small Golang app that executes `SET` and `GET` against Redis with random data when you hit it on `/redis`.

This is built to show how to move connection strings and passwords to a Kubernetes Secret - in order to:

1. Move Kubernetes applications from one service to another quickly.
2. You don't need to update **every** single application Deploy - just a single Secret.
3. Rollback is a single Secret file.
4. We can rollout connection changes application by application if desired. Can start with a small application - if it fails - we can quickly revert just that application.

To make a change:

1. Update the Secret with the new connection string AND/OR password.
2. Deploy that Secret.
3. Perform a `kubectl rollout restart deployment <name>`

Drawbacks?

1. A single secret change can affect multiple applications.
2. If a secret is deployed, the affected applications aren't automatically restarted.