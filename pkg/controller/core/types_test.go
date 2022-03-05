package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestState(t *testing.T) {
	var s State
	assert.Equal(t, "state", s.Type())
	assert.Nil(t, s.Check())
	assert.Equal(t, StateInactive, s)
	assert.Equal(t, "inactive", s.String())

	statesStr := []string{"inactive", "active"}
	states := []State{StateInactive, StateActive}

	for index, stateStr := range statesStr {
		err := s.Set(stateStr)
		assert.Nil(t, err)
		assert.Equal(t, "state", s.Type())
		assert.Nil(t, s.Check())
		assert.Equal(t, states[index], s)
		assert.Equal(t, stateStr, s.String())
	}

	s = StateInactive - 1
	assert.NotNil(t, s.Check())
	s = StateActive + 1
	assert.NotNil(t, s.Check())
}

func TestMode(t *testing.T) {
	var m Mode
	assert.Equal(t, "mode", m.Type())
	assert.Nil(t, m.Check())
	assert.Equal(t, ModeUnknown, m)
	assert.Equal(t, "unknown", m.String())

	modesStr := []string{"unknown", "input", "output"}
	modes := []Mode{ModeUnknown, ModeInput, ModeOutput}

	for index, modeStr := range modesStr {
		err := m.Set(modeStr)
		assert.Nil(t, err)
		assert.Equal(t, "mode", m.Type())
		assert.Nil(t, m.Check())
		assert.Equal(t, modes[index], m)
		assert.Equal(t, modeStr, m.String())
	}

	m = ModeUnknown - 1
	assert.NotNil(t, m.Check())
	m = ModeOutput + 1
	assert.NotNil(t, m.Check())
}

func TestPull(t *testing.T) {
	var p Pull
	assert.Equal(t, "pull", p.Type())
	assert.Nil(t, p.Check())
	assert.Equal(t, PullUnknown, p)
	assert.Equal(t, "unknown", p.String())

	pullsStr := []string{"unknown", "disabled", "down", "up"}
	pulls := []Pull{PullUnknown, PullDisabled, PullDown, PullUp}

	for index, pullStr := range pullsStr {
		err := p.Set(pullStr)
		assert.Nil(t, err)
		assert.Equal(t, "pull", p.Type())
		assert.Nil(t, p.Check())
		assert.Equal(t, pulls[index], p)
		assert.Equal(t, pullStr, p.String())
	}

	p = PullUnknown - 1
	assert.NotNil(t, p.Check())
	p = PullUp + 1
	assert.NotNil(t, p.Check())
}

type TypeMocked struct {
	mock.Mock
}

func (t *TypeMocked) MarshalJSON() ([]byte, error) {
	t.Called()
	return json.Marshal("test")
}

func TestMarshal(t *testing.T) {
	var (
		m       Mode        = ModeOutput
		s       State       = StateActive
		p       Pull        = PullDown
		testObj *TypeMocked = new(TypeMocked)
	)

	mData, err := m.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, `"output"`, string(mData))

	sData, err := s.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, `"active"`, string(sData))

	pData, err := p.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, `"down"`, string(pData))

	testObj.On("MarshalJSON").Return([]byte(""), nil)
	testObj.AssertNotCalled(t, "MarshalJSON")

	tData, err := testObj.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, `"test"`, string(tData))

	x := struct {
		Mode  Mode
		State State
		Pull  Pull
		Test  *TypeMocked
	}{
		Mode:  m,
		State: s,
		Pull:  p,
		Test:  testObj,
	}

	data, err := json.Marshal(x)
	ok := assert.Nil(t, err)
	if !ok {
		t.Logf("err= %v", err.Error())
	}
	assert.Equal(t, `{"Mode":"output","State":"active","Pull":"down","Test":"test"}`, string(data))

	testObj.AssertNumberOfCalls(t, "MarshalJSON", 2)

	m, s, p = -1, -1, -1
	_, err = m.MarshalJSON()
	assert.NotNil(t, err)
	_, err = s.MarshalJSON()
	assert.NotNil(t, err)
	_, err = p.MarshalJSON()
	assert.NotNil(t, err)
}

func TestUnmarshal(t *testing.T) {
	var x struct {
		Mode  Mode
		State State
		Pull  Pull
	}

	err := json.Unmarshal([]byte(`{"Mode":"output","State":"active","Pull":"down"}`), &x)
	assert.Nil(t, err)
	assert.Equal(t, ModeOutput, x.Mode)
	assert.Equal(t, StateActive, x.State)
	assert.Equal(t, PullDown, x.Pull)
}
