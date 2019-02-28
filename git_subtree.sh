#!/bin/bash
git remote add backend git@github.com:ZedYeung/nearby-backend.git
git subtree add --prefix=backend backend master
# git subtree push --prefix=backend backend nearby/master

git remote add frontend git@github.com:ZedYeung/nearby-frontend.git
git subtree add --prefix=frontend frontend master
# git subtree push --prefix=frontend frontend nearby/master