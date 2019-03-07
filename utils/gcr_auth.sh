#!/bin/bash
gcloud auth print-access-token | docker login -u oauth2accesstoken --password-stdin https://us.gcr.io

