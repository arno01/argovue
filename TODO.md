# TODO

* Mount workflow volume and display associated service
* Add shared volumes to services
* Add graph representation for workflows
* Better kube error/info handling on service create/delete

* Add workflow templates and their actions?
* OIDC values remapping: e.g. OIDC groups and ID to more readable?

# FIX

* On connection break give it some time (10-15 seconds) before clean up
* Proxy service must authorize user by checking service label (or generate unique id)

# MAYBE

* Delete brokers on zero connections after timeout
* Better instance navigation (click to proxy)? Need ports
* Stream logs line by line with SSE, the same way as all objects
* Ingress objects for services
* Use helm operator to install services?

# DONE

* Keep URL option for proxy services (it should be always on)
* Display user profile as we get it from OIDC
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
* Delete service instances
