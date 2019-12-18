# TODO

* Use name/presentation name for objects.kubevue.io kind
* Check OIDC groups and how to use them to limit access, probably using another crd, e.g. rules
* Filter services and/or other objects by annotations?
* Keep URL option for proxy services
* Replace ad-hoc maps to Properly defined CRD with generated code
* Helm deployment
* Add run duration to steps/pods/workflows and better time representation (moments)?
* Automate access to volume with nginx pod
* Add workflow templates and their actions?
* Reduce number of unauthorized requests and take action if this happens?

# DONE

* Subscribe to single object (filter?)
* Add workflow actions, e.g. retry
* Honor service settings for redirect (maybe combine with annotations)
* proxy dex a-la argo ci
* Check why resubmit doesn't work
