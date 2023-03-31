# ConsulExample

api.AgentServiceCheck

- DeregisterCriticalServiceAfter
- TLSSSkipVerify
- TTL
- CheckID

api.AgentServiceRegistration

- ID
- Name
- Tags
- Address
- Port
- Check

c.Agent().UpdateTTL(check.checkID, reason, api.HealthPassing)

# watcher

- query (type, service=clustername, passinonly=true)

watch.Parse(query)

plan.HybridHander = func(index watch.BlockingParamVal, result interface{})
plan.RunWithConfig("", &api.Config{}) => blocking
