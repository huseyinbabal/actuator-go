package actuator_test

import (
	"encoding/json"
	"fmt"
	"github.com/huseyinbabal/actuator-go"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestActuator_Health(t *testing.T) {
	a := actuator.New("http://localhost:8081/actuator")
	health, err := a.Health()
	assert.NoError(t, err)

	assert.Equal(t, "UP", health.Status)

}

func TestActuator_Info(t *testing.T) {
	a := actuator.New("http://localhost:8081/actuator")
	info, err := a.Info()
	assert.NoError(t, err)

	assert.Equal(t, "", info.Git.Branch)
}

func TestActuator_HeapDump(t *testing.T) {
	a := actuator.New("http://localhost:8081/actuator")
	heapDump, err := a.HeapDump()
	assert.NoError(t, err)
	dumpFile, err := os.CreateTemp(os.TempDir(), "a.txt")
	assert.NoError(t, err)
	marshal, err := json.Marshal(heapDump)
	err = os.WriteFile(dumpFile.Name(), marshal, os.FileMode(0777))
	assert.NoError(t, err)
	fmt.Print(dumpFile)
	assert.Contains(t, heapDump.Location, "hprof.")
}
