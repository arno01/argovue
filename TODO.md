# TODO

* Check OIDC groups and how to use them to limit access, probably using another crd, e.g. rules

* Filter services and/or other objects by annotations?
* Keep URL option for proxy services
* take name/namespace from crd, e.g. create watchable objects per namespace?

* Replace ad-hoc maps to Properly defined CRD with generated code
* Use name/presentation name for objects.kubevue.io kind

* Add run duration to steps/pods/workflows and better time representation (moments)?

* Automate access to volume with nginx pod
* Add workflow templates and their actions?
* Redirect to front page in case of SSE not-authorized error

# DONE

* Helm deployment and image build workflow
* Reduce number of unauthorized requests
* Subscribe to single object (filter?)
* Add workflow actions, e.g. retry
* Honor service settings for redirect (maybe combine with annotations)
* proxy dex a-la argo ci
* Check why resubmit doesn't work
