#!/bin/bash
git filter-branch --force --index-filter 'git rm --cached -r --ignore-unmatch web-ui' --prune-empty --tag-name-filter cat -- --all
git push origin master --force