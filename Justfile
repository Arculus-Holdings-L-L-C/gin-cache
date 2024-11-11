push-new-tag dry-run="false":
    #!/usr/bin/env bash
    # Get the latest tag from GitHub
    LATEST_TAG=$(git describe --tags `git rev-list --tags --max-count=1`)
    
    # Extract major, minor, patch from tag (assuming format like v1.2.3)
    IFS='.' read -r MAJOR MINOR PATCH <<< "${LATEST_TAG#v}"
    
    # Increment patch version
    NEW_PATCH=$((PATCH + 1))
    
    # Create new tag with incremented patch version
    NEW_TAG="v$MAJOR.$MINOR.$NEW_PATCH"
    
    if [[ "{{dry-run}}" == "true" ]]; then
        echo "Dry run mode - would create and push tag: $NEW_TAG"
    else
        # Create and push new tag
        git tag -a $NEW_TAG -m "Release $NEW_TAG"
        git push origin $NEW_TAG
    fi