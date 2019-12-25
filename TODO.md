# TODO

* Better instance navigation (click to proxy)
* Delete service instances
* Delete watchers on connection close
* Stream logs line by line with SSE, the same way as all objects

* Keep URL option for proxy services
* Add graph representation for workflows
* Add workflow templates and their actions?
* Display user profile as we get it from OIDC
* OIDC values remapping: e.g. OIDC groups and ID to more readable?

# FIX

* On connection break give it some time (10-15 seconds) before clean up
* Proxy service must authorize user

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
* Add run duration to steps/pods/workflows and better time representation (moments)
* Start a copy of pre-defined service per user with shared volumes and private volumes
* Automate access to volume with nginx pod
* Lazy load pod logs
* List workflow PVC
* Add OIDC id as Service selector/copy labels from parent objects
* UI allows to command any workflow (must be only allowed ones)
* UI allows to view any service (must be only allowed ones)
