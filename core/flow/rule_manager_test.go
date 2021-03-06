package flow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAndRemoveTrafficShapingGenerator(t *testing.T) {
	tsc := NewTrafficShapingController(nil, nil, nil)

	err := SetTrafficShapingGenerator(Reject, func(_ *FlowRule) *TrafficShapingController {
		return tsc
	})
	assert.Error(t, err, "default control behaviors are not allowed to be modified")
	err = RemoveTrafficShapingGenerator(Reject)
	assert.Error(t, err, "default control behaviors are not allowed to be removed")

	cb := ControlBehavior(9999)
	err = SetTrafficShapingGenerator(cb, func(_ *FlowRule) *TrafficShapingController {
		return tsc
	})
	assert.NoError(t, err)

	resource := "test-customized-tc"
	_, err = LoadRules([]*FlowRule{
		{
			ID:              10,
			Count:           20,
			MetricType:      QPS,
			Resource:        resource,
			ControlBehavior: cb,
		},
	})
	assert.NoError(t, err)
	assert.Contains(t, tcGenFuncMap, cb)
	assert.NotZero(t, len(tcMap[resource]))
	assert.Equal(t, tsc, tcMap[resource][0])

	err = RemoveTrafficShapingGenerator(cb)
	assert.NoError(t, err)
	assert.NotContains(t, tcGenFuncMap, cb)

	_, _ = LoadRules([]*FlowRule{})
}

func TestIsValidFlowRule(t *testing.T) {
	badRule1 := &FlowRule{ID: 1, Count: 1, MetricType: QPS, Resource: ""}
	badRule2 := &FlowRule{ID: 1, Count: -1.9, MetricType: QPS, Resource: "test"}
	badRule3 := &FlowRule{Count: 5, MetricType: QPS, Resource: "test", ClusterMode: true}
	badRule4 := &FlowRule{Count: 5, MetricType: QPS, Resource: "test", ControlBehavior: WarmUp}
	goodRule1 := &FlowRule{Count: 10, MetricType: QPS, Resource: "test", ControlBehavior: Throttling}

	assert.Error(t, IsValidFlowRule(badRule1))
	assert.Error(t, IsValidFlowRule(badRule2))
	assert.Error(t, IsValidFlowRule(badRule3))
	assert.Error(t, IsValidFlowRule(badRule4))
	assert.NoError(t, IsValidFlowRule(goodRule1))
}
