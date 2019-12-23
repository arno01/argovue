# TODO

* Start a copy of pre-defined service per user with shared volumes and private volumes
* Automate access to volume with nginx pod
* Keep URL option for proxy services
* Add graph representation for workflows
* Add run duration to steps/pods/workflows and better time representation (moments)?
* Add workflow templates and their actions?
* Display user profile as we get it from OIDC

# DONE

* Helm deployment and image build workflow
* Reduce number of unauthorized requests
* Subscribe to single object (filter?)
* Add workflow actions, e.g. retry
* Honor service settings for redirect (maybe combine with annotations)
* proxy dex a-la argo ci
* Check why resubmit doesn't work
* Filter services and/or other objects by annotations?
* take name/namespace from crd, e.g. create watchable objects per namespace?
* Replace ad-hoc maps to properly defined CRD with generated code
* Use name/presentation name for objects.argovue.io kind
* Redirect to front page in case of SSE not-authorized error
* Display namespace for objects
