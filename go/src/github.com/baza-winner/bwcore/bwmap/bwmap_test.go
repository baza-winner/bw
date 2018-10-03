package bwmap

import (
  "fmt"
  "github.com/baza-winner/bwcore/ansi"
  "github.com/baza-winner/bwcore/bwtesting"
  "github.com/mohae/deepcopy"
  "testing"
)

type testCropMapStruct struct {
  m      interface{}
  crop   interface{}
  result interface{}
}

func TestCropMap(t *testing.T) {
  tests := map[string]testCropMapStruct{
    "string": {
      m: map[string]interface{}{
        "some": "thing",
        "good": "is not bad",
      },
      crop: `some`,
      result: map[string]interface{}{
        "some": "thing",
      },
    },
    "[]string": {
      m: map[string]interface{}{
        "A": 1,
        "B": 2,
        "C": 3,
        "D": 4,
      },
      crop: []string{"B", "C"},
      result: map[string]interface{}{
        "B": 2,
        "C": 3,
      },
    },
    "map[string]interface{}": {
      m: map[string]int{
        "A": 1,
        "B": 2,
        "C": 3,
        "D": 4,
      },
      crop: map[string]interface{}{
        "A": struct{}{},
        "D": struct{}{},
      },
      result: map[string]int{
        "A": 1,
        "D": 4,
      },
    },
  }
  testsToRun := tests
  for testName, test := range testsToRun {
    t.Logf(ansi.Ansi(`Header`, "Running test case <ansiPrimaryLiteral>%s"), testName)
    result := deepcopy.Copy(test.m)
    CropMap(result, test.crop)
    testTitle := fmt.Sprintf("CropMap(%+v, %+v)\n", test.m, test.crop)
    bwtesting.DeepEqual(t, result, test.result, testTitle)
  }
}
