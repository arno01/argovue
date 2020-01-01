# TODO

* Stream logs line by line with SSE, the same way as all objects
* Better pod/service presentation (display significant fields on main page)
* Add all deployment objects on service instance tab, with their statuses

* On connection break give it some time (10-15 seconds) before clean up
* Better kube error/info handling on service create/delete (with rollback?)
* OIDC values remapping: e.g. OIDC groups and ID to more readable?

# MAYBE

* Delete brokers on zero connections after timeout
* Better instance navigation (click to proxy)? Need ports
* Ingress objects for services
* Use helm operator to install services?
* Add workflow templates and their actions?
* Switch to redis (or any other distributed kv) to store sessions
* Fix dex proxy (service should be on even if oidc fails)

# DONE

* UI configurator: make env variables accessible?
* Send keep-alive events
* Keep tabs on navigation
* Generate global session key as part of deployment/start up
* Different node colors for statesz
* Clickable graph nodes (double click)
* Standalone UI
* Proxy service must authorize user by checking service label (or generate unique id)
* Display objects ownership (myself?/group)
* Remove Retry nodes from Graph
* Add graph representation for workflows
* Verify workflow file browser is mounted
* Display workflow filebrowser service and allow to remove it
* Mount workflow volume
* Add shared volumes to services
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
