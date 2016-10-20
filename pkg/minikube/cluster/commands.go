/*
Copyright (C) 2016 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cluster

import (
	"fmt"

	"github.com/jimmidyson/minishift/pkg/minikube/constants"
)

// Kill any running instances.
var stopCommand = "sudo killall openshift || true"

var startCommandFmtStr = `
# Run with nohup so it stays up. Redirect logs to useful places.
cd /var/lib/minishift;
sudo mkdir -p /mnt/sda1/var/lib/minishift/openshift.local.volumes /mnt/sda1/var/lib/minishift/openshift.local.config /mnt/sda1/var/lib/minishift/openshift.local.etcd || true
sudo sh -c 'PATH=/usr/local/sbin:$PATH nohup openshift cli cluster up \
		--host-config-dir=$(pwd)/openshift.local.config \
		--host-data-dir=$(pwd)/openshift.local.etcd \
		--host-volumes-dir=$(pwd)/openshift.local.volumes \
		--routing-suffix=%s.nip.io \
		--use-existing-config \
		> %s 2> %s < /dev/null &'
until $(curl --output /dev/null --silent --fail -k https://localhost:%d/healthz/ready); do
    printf '.'
    sleep 1
done;
sudo /usr/local/bin/openshift admin policy add-cluster-role-to-user cluster-admin admin --config=$(pwd)/openshift.local.config/master/admin.kubeconfig
`

var (
	logsCommand      = fmt.Sprintf("tail -n +1 %s %s", constants.RemoteOpenShiftErrPath, constants.RemoteOpenShiftOutPath)
	getCACertCommand = fmt.Sprintf("cat %s", constants.RemoteOpenShiftCAPath)
)

func GetStartCommand(ip string) string {
	return fmt.Sprintf(startCommandFmtStr, ip, constants.RemoteOpenShiftErrPath, constants.RemoteOpenShiftOutPath, constants.APIServerPort)
}
