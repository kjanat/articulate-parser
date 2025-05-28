## Docker Non-Root User Implementation

I've addressed the concern about the Docker container running as root while the README claims non-root execution. The implementation includes:

1. Added a non-root user (appuser with UID 1000) in the builder stage
2. Copy the passwd file to the scratch image to make the USER directive work
3. Added USER directive to run the container as the non-root user
4. Created Dockerfile.dev for development with shell access (using Alpine instead of scratch)
5. Fixed docker-compose.yml to use proper YAML anchors

This change ensures that the Docker container actually runs as a non-privileged user as stated in the README, which is an important security best practice.

## Other PR Comments

Regarding the other PR comments:

1. For prerelease flag in release.yml:
   - The suggested change to add a comment explaining the rationale for using 'startsWith(github.ref, \"refs/tags/v0.\")' is good practice
   - v0.x.x versions follow semver convention for pre-1.0 software that may have breaking changes

2. For dependency-review.yml and codeql.yml workflow trigger changes:
   - Switching from 'pull_request' to 'workflow_call' changes how these workflows are triggered
   - Documentation should be added to explain when and how these workflows should be called
   - This may impact security scanning, so careful consideration is needed
