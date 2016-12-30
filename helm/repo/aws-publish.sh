#!/bin/bash

aws s3 sync charts/ s3://hspc-helm.preparedmind.net
