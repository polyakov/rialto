#!/bin/bash

helm serve --repo-path charts/ &
helm repo add hspc-helm http://127.0.0.1:8879
