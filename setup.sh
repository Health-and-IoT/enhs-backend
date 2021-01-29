#!/bin/bash
cd
echo Setting Up PATH Variables.
echo export GOPATH="/shared/go" >> .profile
echo export GOBIN=$(go env GOPATH)/bin >> .profile
echo export PATH=$(go env GOPATH)/bin:$PATH >> .profile
echo $(whoami), you will need to restart the client for changes to take effect.
echo When you have logged back in, check to see if the GOPATH is set to "/shared/go" by using the command: "go env GOPATH"
