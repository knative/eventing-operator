# Upgrade Guide

This document describes how to upgrade the eventing-operator to an expected
version.

## Backup precaution

As an administrator, you are recommended to save the content of the custom
resource for eventing-operator before upgrading your operator. Make sure that
you know the name and the namespace of your CR, and use the following command to
save the CR in a file called `eventing_operator_cr.yaml`:

For the version v0.12.0 or later:

```
kubectl get KnativeEventing <cr-name> -n <namespace> -o=yaml > eventing_operator_cr.yaml
```

For the version v0.11.0 or earlier:

```
kubectl get Eventing <cr-name> -n <namespace> -o=yaml > eventing_operator_cr.yaml
```

Replace `<cr-name>` with the name of your CR, and `<namespace>` with the
namespace.

One version of eventing-operator installs only one specific version of Knative
Eventing. With your operator successfully upgraded, your Knative Eventing is
upgraded as well.

## v0.11.0 -> v0.12.0

The Kind name of the custom resource has been changed from `Eventing` to
`KnativeEventing`. The version v0.12.0 is able to recognize both of the CRs.

Update the eventing-operator to the version v0.12.0 with the command:

```
kubectl apply -f https://github.com/knative/eventing-operator/releases/download/v0.12.0/eventing-operator.yaml
```

## v0.10.0 -> v0.11.0

If your existing operator is at the version v0.10.0, you cannot directly upgrade
it to v0.12.0 or later. We only support upgrade from v0.10.0 to v0.11.0.
Existing issues about directly upgrading from v0.10.0 to v0.12.0 has been report
at [Bug 54](https://github.com/knative/eventing-operator/issues/54).

Update the eventing-operator deployment to the version v0.11.0 with the command:

```
kubectl apply -f https://github.com/knative/eventing-operator/releases/download/v0.11.0/eventing-operator.yaml
```
