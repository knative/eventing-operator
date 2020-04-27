FROM openshift/origin-release:golang-1.13 AS builder
WORKDIR ${GOPATH}/src/knative.dev/eventing-operator
COPY . .
ENV GOFLAGS="-mod=vendor"
RUN go build -o /tmp/manager ./cmd/manager
RUN cp -Lr ${GOPATH}/src/knative.dev/eventing-operator/cmd/manager/kodata /tmp

FROM openshift/origin-base
COPY --from=builder /tmp/manager /ko-app/eventing-operator
COPY --from=builder /tmp/kodata/ /var/run/ko
ENV KO_DATA_PATH="/var/run/ko"
LABEL \
    com.redhat.component="openshift-serverless-1-tech-preview-knative-eventing-rhel8-operator-container" \
    name="openshift-serverless-1-tech-preview/knative-eventing-rhel8-operator" \
    version="v0.14.1" \
    summary="Red Hat OpenShift Serverless 1 Knative Eventing Operator" \
    maintainer="serverless-support@redhat.com" \
    description="Red Hat OpenShift Serverless 1 Knative Eventing Operator" \
    io.k8s.display-name="Red Hat OpenShift Serverless 1 Knative Eventing Operator"

ENTRYPOINT ["/ko-app/eventing-operator"]
