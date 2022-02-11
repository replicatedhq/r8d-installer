#!/bin/bash
# Author: K. Boyle & D. Stough

# TODO: Preflights
# Check for NetworkManager `systemctl reload NetworkManager`
# Check for AppArmor
# Check for Root

curl -sfL https://get.rke2.io | bash

systemctl enable rke2-server.service
systemctl start rke2-server.service

# Add the RKE2 Kube Tools to the path
chmod 644 "/etc/rancher/rke2/rke2.yaml" # TODO: This can be controlled in the configuration file for the server
echo "export KUBECONFIG=/etc/rancher/rke2/rke2.yaml" >> ~/.bashrc
echo 'export PATH=$PATH:/var/lib/rancher/rke2/bin/' >> ~/.bashrc

# Install the CLI binaries

# KOTS CLI
echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc

mkdir assets

LATEST_RELEASE=$(curl -L -s -H 'Accept: application/json' https://github.com/replicatedhq/kots/releases/latest)
LATEST_VERSION=$(echo $LATEST_RELEASE | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
ARTIFACT_URL="https://github.com/replicatedhq/kots/releases/download/$LATEST_VERSION/kots_linux_amd64.tar.gz"
curl -L "$ARTIFACT_URL" -o kots_linux_amd64.tar.gz
tar zxf kots_linux_amd64.tar.gz -C assets

mv assets/kots /usr/local/bin/kubectl-kots

# Troubleshoot CLI

LATEST_RELEASE=$(curl -L -s -H 'Accept: application/json' https://github.com/replicatedhq/troubleshoot/releases/latest)
LATEST_VERSION=$(echo $LATEST_RELEASE | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
ARTIFACT_URL="https://github.com/replicatedhq/troubleshoot/releases/download/$LATEST_VERSION/support-bundle_linux_amd64.tar.gz"
curl -L "$ARTIFACT_URL" -o support-bundle_linux_amd64.tar.gz
tar zxf support-bundle_linux_amd64.tar.gz -C assets

mv assets/support-bundle /usr/local/bin/kubectl-support_bundle

# Velero CLI 

LATEST_RELEASE=$(curl -L -s -H 'Accept: application/json' https://github.com/vmware-tanzu/velero/releases/latest)
LATEST_VERSION=$(echo $LATEST_RELEASE | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
ARTIFACT_URL="https://github.com/vmware-tanzu/velero/releases/download/${LATEST_VERSION}/velero-${LATEST_VERSION}-linux-amd64.tar.gz"
curl -L "$ARTIFACT_URL" -o velero-${LATEST_VERSION}-linux-amd64.tar.gz
tar zxf velero-${LATEST_VERSION}-linux-amd64.tar.gz -C assets

mv assets/velero-${LATEST_VERSION}-linux-amd64/velero /usr/local/bin/velero

# CSI Provider: Local-Path-Provisioner Manifests
# IMPORTANT: Doesn't support snapshots with restic
# sudo curl -sSL "https://raw.githubusercontent.com/rancher/local-path-provisioner/master/deploy/local-path-storage.yaml" -o /var/lib/rancher/rke2/server/manifests/local-path-storage.yaml

# OpenEBS Lite Operator
sudo curl -sSL "https://openebs.github.io/charts/openebs-operator-lite.yaml" -o /var/lib/rancher/rke2/server/manifests/openebs-operator-lite.yaml
sudo curl -sSL "https://openebs.github.io/charts/openebs-lite-sc.yaml" -o /var/lib/rancher/rke2/server/manifests/openebs-lite-sc.yaml
kubectl patch storageclass openebs-hostpath  -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'


# KOTS Manifests
mkdir -p assets
kubectl kots admin-console generate-manifests --with-minio=false --namespace=default --shared-password=password --rootdir=/home/dans/r8d-installer/assets
for i in ./assets/admin-console/*.yaml ; 
do 
    cat "${i}" >> kots.yaml
    echo "---" >> kots.yaml
done

mv kots.yaml /var/lib/rancher/rke2/server/manifests/kots.yaml

# Velero Manifests
# This won't work due to problems snapshoting local volume provisioners
# TODO: missing the kurl plugin here
# TODO: missing BSL will cause KOTS to think velero isn't installed
velero install --use-restic --no-secret --no-default-backup-location --namespace velero --plugins velero/velero-plugin-for-aws:v1.3.0,velero/velero-plugin-for-gcp:v1.3.0,velero/velero-plugin-for-microsoft-azure:v1.3.0,replicated/local-volume-provider:v0.3.0 --use-volume-snapshots=false --dry-run -o yaml > velero.yaml
mv velero.yaml /var/lib/rancher/rke2/server/manifests/velero.yaml

